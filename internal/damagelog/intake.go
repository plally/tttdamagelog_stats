package damagelog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/plally/damagelog_parser/internal/dal"
	"github.com/plally/damagelog_parser/internal/models"
)

func createRound(ctx context.Context, q *dal.Queries, e Entry) (*dal.Round, error) {
	round, err := q.CreateRound(ctx, dal.CreateRoundParams{
		Map: e.Map,
		StartTime: pgtype.Timestamp{
			Valid: true,
			Time:  time.Unix(int64(e.Date), 0),
		},
		EndTime: pgtype.Timestamp{
			Valid: true,
			Time:  time.Unix(int64(e.Date), 0),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create round: %w", err)
	}

	return &round, nil
}

func parseShootTable(round *dal.Round, shootTable map[string][][]any, roles []DamagelogRole) []dal.CreateEventsParams {
	var events []dal.CreateEventsParams
	for k, event := range shootTable {
		roundTime, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		if len(event) == 0 {
			continue
		}
		shootData := event[0]

		attackerID := shootData[0].(float64)
		attacker := roles[int(attackerID)-1].Steamid64

		weapon := shootData[1].(string)

		shootEvent := models.ShootEventData{
			Attacker: attacker,
			Weapon:   weapon,
		}
		eventData, err := json.Marshal(shootEvent)
		if err != nil {
			continue
		}

		events = append(events, dal.CreateEventsParams{
			RoundTime: float64(roundTime),
			RoundID:   round.ID,
			EventType: models.EventTypeShoot,
			EventData: eventData,
		})
	}

	return events
}

func parseDamageTable(round *dal.Round, damageTable []DamagelogEvent, roles []DamagelogRole) ([]dal.CreateEventsParams, error) {
	var events []dal.CreateEventsParams
	for _, event := range damageTable {
		eventType, eventData := EventMapping(roles, &event)
		eventDataBytes, err := json.Marshal(eventData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal event data: %w", err)
		}

		events = append(events, dal.CreateEventsParams{
			RoundTime: float64(event.Time),
			RoundID:   round.ID,
			EventType: eventType,
			EventData: eventDataBytes,
		})
	}
	return events, nil
}

func ProcessData(ctx context.Context, reader io.Reader, q *dal.Queries) error {
	var events []dal.CreateEventsParams

	var e Entry
	json.NewDecoder(reader).Decode(&e)

	round, err := createRound(ctx, q, e)
	if err != nil {
		return fmt.Errorf("failed to create round: %w", err)
	}

	var dmgLog Damagelog
	json.Unmarshal([]byte(e.Damagelog), &dmgLog)

	newEvents := parseShootTable(round, dmgLog.ShootTable, dmgLog.Roles)
	events = append(events, newEvents...)

	newEvents, err = parseDamageTable(round, dmgLog.DamageTable, dmgLog.Roles)
	if err != nil {
		return fmt.Errorf("failed to parse damage table: %w", err)
	}
	events = append(events, newEvents...)

	slog.With("eventsLen", len(events)).Debug("creating events")

	endTime := 0
	for i, e := range events {
		if int(e.RoundTime) > endTime {
			endTime = int(e.RoundTime)
		}
		events[i].EventTime = pgtype.Timestamp{
			Valid: true,
			Time:  round.StartTime.Time.Add(time.Duration(e.RoundTime) * time.Second),
		}
	}

	for _, role := range dmgLog.Roles {
		steamid, err := strconv.ParseInt(role.Steamid64, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse steamid: %w", err)
		}

		_, err = q.PutUser(ctx, dal.PutUserParams{
			Steamid64: steamid,
			Nickname:  role.Nick,
		})
		if err != nil {
			return fmt.Errorf("failed to put user: %w", err)
		}
	}

	_, err = q.CreateEvents(ctx, events)
	if err != nil {
		return fmt.Errorf("failed to create events: %w", err)
	}

	_, err = q.UpdateRoundEndTime(ctx, dal.UpdateRoundEndTimeParams{
		ID: round.ID,
		EndTime: pgtype.Timestamp{
			Valid: true,
			Time:  round.StartTime.Time.Add(time.Duration(endTime) * time.Second),
		},
	})

	return err
}

type Entry struct {
	Damagelog string
	Round     int
	Map       string
	Date      float64
}

type Damagelog struct {
	Roles       []DamagelogRole
	ShootTable  map[string][][]any
	DamageTable []DamagelogEvent `json:"DamageTable"`
}

type DamagelogEvent struct {
	Time  float64 `json:"time"`
	Infos []any   `json:"infos"`
	ID    float64 `json:"id"`
}

type DamagelogRole struct {
	Nick      string
	Role      float64
	Steamid64 string
}
