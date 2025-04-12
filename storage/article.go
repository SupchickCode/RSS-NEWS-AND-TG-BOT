package storage

import (
	"context"
	"tg-bot-supchick/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

type dbArticle struct {
	ID          int64     `db:"id"`
	SourceID    int64     `db:"source_id"`
	Title       string    `db:"title"`
	Link        string    `db:"link"`
	Summary     string    `db:"summary"`
	PublishedAt time.Time `db:"published_at"`
	PostedAt    time.Time `db:"posted_at"`
	CreatedAt   time.Time `db:"created_at"`
}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit int64) ([]models.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var dbArticles []dbArticle
	if err := conn.SelectContext(ctx,
		&dbArticles,
		`SELECT * FROM articles 
		WHERE posted_at IS NULL
		AND published_at >= $1::timestamp
		ORDER BY DESC published_at LIMIT $2`,
		since.UTC().Format(time.RFC3339),
		limit,
	); err != nil {
		return nil, err
	}

	var articles []models.Article
	for _, item := range dbArticles {
		articles = append(articles, models.Article{
			ID:          item.ID,
			SourceID:    item.SourceID,
			Title:       item.Title,
			Summary:     item.Summary,
			PublishedAt: item.PublishedAt,
			PostedAt:    item.PostedAt,
			CreatedAt:   item.CreatedAt,
		})
	}

	return articles, nil
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, m models.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO articles (source_id, title, summary, link, published_at) VALUES ($1, $2, $3, $4, $5)`,
		m.SourceID,
		m.Title,
		m.Summary,
		m.Link,
		m.PostedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`UPDATE sources SET posted_at = $1::timestamp WHERE id = $2`,
		time.Now().UTC().Format(time.RFC3339), id,
	); err != nil {
		return err
	}

	return nil
}
