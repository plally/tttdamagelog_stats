newmigration:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir db/migrations -seq $$name

migrateup:
	migrate -database ${POSTGRESQL_URL} -path db/migrations up
