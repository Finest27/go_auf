package services

import (
	"context"
	"thepress_bot_go/internal/domain/models"
)

type Publisher interface {
	Publish(article *models.Article, catID int) (string, error)
}

type Linker interface {
	InjectLinks(ctx context.Context, excludeID int64, htmlContent string) string
}

// ScraperService defines the contract for fetching articles
type ScraperService interface {
	FetchRSSLinks(ctx context.Context, url string) ([]string, error)
	ScrapeArticle(ctx context.Context, url string) (title string, content string, image string, err error)
	Close()
}
