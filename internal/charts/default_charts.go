package charts

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/plally/damagelog_parser/internal/dal"
)

type chartGenerator func(ctx context.Context, steamID64 string, q *dal.Queries) (Chart, error)

var defaultCharts = map[string]chartGenerator{
	"PlayerKilledByMost": PlayerKilledByMost,
	"MostPlayedMaps":     MostPlayedMaps,
	"MostUsedWeapon":     MostUsedWeapon,
	"PlayerTopKills":     PlayerTopKills,
}

func GetRandomChart(ctx context.Context, steamID64 string, q *dal.Queries) (Chart, error) {
	length := len(defaultCharts)
	i := rand.Intn(length)
	for _, generator := range defaultCharts {
		if i == 0 {
			return generator(ctx, steamID64, q)
		}
		i--
	}

	return Chart{}, fmt.Errorf("failed to get random chart")
}

func GetDefaultCharts(ctx context.Context, steamID64 string, q *dal.Queries) (map[string]Chart, error) {

	charts := make(map[string]Chart)
	for name, generator := range defaultCharts {
		chart, err := generator(ctx, steamID64, q)
		if err != nil {
			return nil, err
		}

		charts[name] = chart
	}

	return charts, nil
}

func MostPlayedMaps(ctx context.Context, steamID64 string, q *dal.Queries) (Chart, error) {
	input := PieChartInput{
		Label: "Most Played Maps in the last 30 days",
	}

	dbData, err := q.StatsGetMostPlayedMaps(ctx, pgtype.Timestamp{
		Valid: true,
		Time:  time.Now().Add(-30 * 24 * time.Hour),
	})
	if err != nil {
		return Chart{}, err
	}
	if len(dbData) > 6 {
		dbData = dbData[:6]
	}
	for _, row := range dbData {
		input.Items = append(input.Items, PieChartItem{
			Label: row.Map,
			Data:  int(row.PlayCount),
		})
	}

	return PieChart(input), nil
}

func PlayerKilledByMost(ctx context.Context, steamID64 string, q *dal.Queries) (Chart, error) {
	input := BarChartInput{
		Label: "Players who have killed you the most",
	}

	dbData, err := q.StatsGetPlayerKilledByMost(ctx, steamID64)
	if err != nil {
		return Chart{}, err
	}

	for _, row := range dbData {
		input.Items = append(input.Items, BarChartItem{
			Data:  int(row.TotalKills),
			Label: GetAttackerName(ctx, row.Attacker, q),
		})
	}

	return BarChart(input), nil
}

func MostUsedWeapon(ctx context.Context, steamID64 string, q *dal.Queries) (Chart, error) {
	input := PieChartInput{
		Label: "Most used weapons in the last 30 days",
	}

	dbData, err := q.StatsGetMostUsedWeaponsPerRound(ctx, pgtype.Timestamp{
		Valid: true,
		Time:  time.Now().Add(-30 * 24 * time.Hour),
	})
	if err != nil {
		return Chart{}, err
	}

	for _, row := range dbData {
		if row.Weapon == "" {
			continue
		}
		if row.Weapon == nil {
			continue
		}

		input.Items = append(input.Items, PieChartItem{
			Data:  int(row.UsageCount),
			Label: fmt.Sprintf("%v", row.Weapon),
		})
		if len(input.Items) > 6 {
			break
		}
	}

	return PieChart(input), nil
}

func PlayerTopKills(ctx context.Context, steamID64 string, q *dal.Queries) (Chart, error) {
	input := BarChartInput{
		Label: "Players with the most kills in the last 7 days",
	}

	dbData, err := q.StatsGetPlayersTopKills(ctx, pgtype.Timestamp{
		Valid: true,
		Time:  time.Now().Add(-7 * 24 * time.Hour),
	})
	if err != nil {
		return Chart{}, err
	}

	for _, row := range dbData {
		input.Items = append(input.Items, BarChartItem{
			Data:  int(row.TotalKills),
			Label: GetAttackerName(ctx, row.Attacker, q),
		})
	}

	return BarChart(input), nil
}

func GetAttackerName(ctx context.Context, v any, q *dal.Queries) string {
	switch v := v.(type) {
	case string:
		steamID := v
		steamID64, err := strconv.ParseInt(steamID, 10, 64)
		if err != nil {
			return steamID
		}

		user, err := q.GetUser(ctx, int64(steamID64))
		if err != nil {
			return steamID
		}

		return user.Nickname
	default:
		return fmt.Sprintf("%v", v)
	}
}
