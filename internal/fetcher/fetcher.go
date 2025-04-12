package fetcher

import (
	"context"
	"log"
	"strings"
	"sync"
	"tg-bot-supchick/helper"
	"tg-bot-supchick/internal/models"
	"tg-bot-supchick/internal/source"
	"time"
)

type ArticleStorage interface {
	Store(ctx context.Context, a models.Article) error
}

type SourceProvider interface {
	Sources(ctx context.Context) ([]models.Source, error)
}

type Source interface {
	ID() int64
	Name() string
	Fetch(ctx context.Context) ([]models.Item, error)
}

type Fetcher struct {
	articles ArticleStorage
	sources  SourceProvider

	fetchInterval time.Duration
	fetchKeywords []string
}

func NewFetcher(
	articles ArticleStorage,
	sources SourceProvider,
	fetchInterval time.Duration,
	fetchKeywords []string,
) *Fetcher {
	return &Fetcher{
		articles:      articles,
		sources:       sources,
		fetchInterval: fetchInterval,
		fetchKeywords: fetchKeywords,
	}
}

func (f *Fetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(f.fetchInterval)
	defer ticker.Stop()

	if err := f.Fetch(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return err
			}
		}
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	sources, err := f.sources.Sources(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, src := range sources {
		wg.Add(1)

		RSSSource := source.NewRSSSourceFromModel(&src)

		go func(source Source) {
			defer wg.Done()

			items, err := RSSSource.Fetch(ctx)
			if err != nil {
				log.Printf(err.Error())
				return
			}

			if err := f.processItems(ctx, source, items); err == nil {
				log.Printf(err.Error())
				return
			}

		}(RSSSource)
	}

	wg.Wait()

	return nil
}

func (f *Fetcher) processItems(ctx context.Context, source Source, items []models.Item) error {
	for _, item := range items {
		item.Date = item.Date.UTC()

		if f.itemShouldBeSkipped(item) {
			continue
		}

		if err := f.articles.Store(ctx, models.Article{
			SourceID:    source.ID(),
			Title:       item.Title,
			Link:        item.Link,
			Summary:     item.Summary,
			PublishedAt: item.Date,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fetcher) itemShouldBeSkipped(item models.Item) bool {
	categories := item.Categories

	for _, keyWord := range f.fetchKeywords {
		if helper.Contains(categories, keyWord) ||
			helper.Contains(categories, strings.ToLower(item.Title)) {
			return true
		}
	}

	return false
}
