package repository

import (
	"context"
	"thepress_bot_go/internal/domain/models"
)

type ArticleRepository interface {
	Save(ctx context.Context, article *models.Article) error
	Update(ctx context.Context, article *models.Article) error
	GetByID(ctx context.Context, id int64) (*models.Article, error)
	GetPending(ctx context.Context, limit int) ([]models.Article, error)
	Exists(ctx context.Context, url string) (bool, error)
	GetRelated(ctx context.Context, articleID int64, limit int) ([]models.Article, error)
}

type FeedRepository interface {
	Add(ctx context.Context, feed *models.Feed) error
	GetAllActive(ctx context.Context) ([]models.Feed, error)
	UpdateLastCheck(ctx context.Context, id int64) error
}
