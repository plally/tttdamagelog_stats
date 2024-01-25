package models

type Event struct {
	ID        int
	EventType string
	EventData any
	RoundID   int
}

type DamageEventData struct {
	Victim   string
	Attacker string
	Damage   float64
	Weapon   string
}

type KillEventData struct {
	Victim   string
	Attacker string
	Weapon   string
}

type ShootEventData struct {
	Attacker string
	Weapon   string
}

const (
	EventTypeShoot  = "shoot"
	EventTypeDamage = "damage"
	EventTypeKill   = "kill"
	EventTypeInfo   = "infos"
)
