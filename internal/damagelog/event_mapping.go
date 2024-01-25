package damagelog

import (
	"fmt"

	"github.com/plally/damagelog_parser/internal/models"
)

var mappings = map[int]string{
	5:  "DMG",
	11: "KILL",
}

func EventMapping(roles []DamagelogRole, event *DamagelogEvent) (string, any) {
	eventType := mappings[int(event.ID)]

	switch eventType {
	case "DMG":
		if len(event.Infos) < 4 {
			return models.EventTypeInfo, event.Infos
		}

		victimID, ok := event.Infos[0].(float64)
		if !ok {
			return models.EventTypeInfo, event.Infos
		}
		attackID, ok := event.Infos[1].(float64)
		if !ok {
			return models.EventTypeInfo, event.Infos
		}

		victim := roles[int(victimID)-1]
		attacker := roles[int(attackID)-1]

		return models.EventTypeDamage, models.DamageEventData{
			Victim:   victim.Steamid64,
			Attacker: attacker.Steamid64,
			Damage:   getInfos[float64](event.Infos, 2),
			Weapon:   getInfos[string](event.Infos, 3),
		}
	case "KILL":
		if len(event.Infos) < 3 {
			fmt.Println("KILL INVALID", event.Infos)
			return models.EventTypeInfo, event.Infos
		}
		fmt.Println("KILL VALID", event.Infos)
		attackerID, ok := event.Infos[0].(float64)
		if !ok {
			return models.EventTypeInfo, event.Infos
		}
		attacker := roles[int(attackerID)-1]

		victimID, ok := event.Infos[1].(float64)
		if !ok {
			return models.EventTypeInfo, event.Infos
		}
		victim := roles[int(victimID)-1]

		weapon := getInfos[string](event.Infos, 2)

		return models.EventTypeKill, models.KillEventData{
			Victim:   victim.Steamid64,
			Attacker: attacker.Steamid64,
			Weapon:   weapon,
		}
	}
	return models.EventTypeInfo, event.Infos
}

func getInfos[T any](infos []any, index int) T {
	v, ok := infos[index].(T)
	if !ok {
		var zero T
		return zero
	}
	return v
}
