MIGRATE_DB_DSN=host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable

migrate-up:
	goose postgres "$(MIGRATE_DB_DSN)" up
