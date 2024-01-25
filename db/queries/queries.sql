-- name: PutUser :one
INSERT INTO users (steamid64, nickname) VALUES ($1, $2) ON CONFLICT (steamid64) DO UPDATE SET nickname = $2 RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE steamid64 = $1;

-- name: CreateRound :one
INSERT INTO rounds (map, start_time, end_time) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateRoundEndTime :one
UPDATE rounds SET end_time = $1 WHERE id = $2 RETURNING *;

-- name: CreateEvent :one
INSERT INTO events (round_id, event_type, event_data, round_time, event_time) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: CreateEvents :copyfrom
INSERT INTO events (round_id, event_type, event_data, round_time, event_time) VALUES ($1, $2, $3, $4, $5);

-- name: GetRound :one
SELECT * FROM rounds WHERE id = $1;

-- name: GetRoundEvents :many
SELECT * FROM events WHERE round_id = $1;

-- name: GetEvents :many
SELECT * FROM events;

-- name: StatsGetMostUsedWeaponsPerRound :many
SELECT event_data->>'Weapon' AS weapon, COUNT(DISTINCT round_id) as usage_count FROM events WHERE event_time > $1 GROUP BY event_data->>'Weapon' ORDER BY usage_count DESC;

-- name: StatsGetMostPlayedMaps :many
SELECT map, COUNT(map) as play_count FROM rounds WHERE start_time > $1 GROUP BY map ORDER BY play_count DESC;

-- name: StatsGetTotalWeaponDamage :many
SELECT SUM((event_data->>'Damage')::int) as total_damage, event_data->>'Weapon' as weapon FROM events WHERE event_time > $1 AND event_type='damage' GROUP BY event_data->>'Weapon' ORDER BY total_damage DESC;

-- name: StatsGetPlayerKilledByMost :many
SELECT event_data->>'Attacker' as attacker, COUNT(event_data->>'Attacker') as total_kills FROM events WHERE event_type='kill' AND event_data->>'Victim'=sqlc.arg(victim)::text GROUP BY event_data->>'Attacker' ORDER BY total_kills DESC LIMIT 5;
