package storage

import (
	"context"
	"tg-bot-supchick/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type SourcePostgresStorage struct {
	db *sqlx.DB
}

type dbSource struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]models.Source, error) {
	conn, err := s.db.Connx(ctx)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	var dbSources []dbSource
	if err := conn.SelectContext(ctx, &dbSources, `SELECT * FROM sources`); err != nil {
		return nil, err
	}

	var sources []models.Source
	for _, item := range dbSources {
		sources = append(sources, models.Source{
			ID:        item.ID,
			Name:      item.Name,
			FeedURL:   item.FeedURL,
			CreatedAt: item.CreatedAt,
		})
	}

	return sources, nil
}

func (s *SourcePostgresStorage) SourceByID(ctx context.Context, id int64) (*models.Source, error) {
	conn, err := s.db.Connx(ctx)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	var dbSource dbSource
	if err := conn.GetContext(ctx, &dbSource, `SELECT * FROM sources WHERE id = $1`, id); err != nil {
		return nil, err
	}

	return &models.Source{
		ID:        dbSource.ID,
		Name:      dbSource.Name,
		FeedURL:   dbSource.FeedURL,
		CreatedAt: dbSource.CreatedAt,
	}, nil
}

func (s *SourcePostgresStorage) Add(ctx context.Context, m models.Source) (int64, error) {
	conn, err := s.db.Connx(ctx)

	if err != nil {
		return 0, err
	}

	defer conn.Close()

	var id int64
	row := conn.QueryRowxContext(
		ctx,
		`INSERT INTO sources (name, feed_url, created_at) VALUES ($1, $2, $3)`,
		m.Name,
		m.FeedURL,
		m.CreatedAt,
	)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SourcePostgresStorage) Delete(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)

	if err != nil {
		return err
	}

	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `DELETE FROM sources WHERE id = $1`, id); err != nil {
		return err
	}

	return nil
}
