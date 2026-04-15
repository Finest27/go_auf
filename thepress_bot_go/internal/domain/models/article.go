package models

import (
	"database/sql"
	"time"
)

type Article struct {
	ID               int64          `db:"id"`
	SourceURL        string         `db:"source_url"`
	Title            string         `db:"title"`
	Content          string         `db:"content"`
	RewrittenContent sql.NullString `db:"rewritten_content"`
	ImageURL         sql.NullString `db:"image_url"`
	MetaDescription  sql.NullString `db:"meta_description"`
	FocusKeywords    sql.NullString `db:"focus_keywords"`
	Slug             sql.NullString `db:"slug"`
	Category         sql.NullString `db:"category"`
	CategoryID       sql.NullInt64  `db:"category_id"`
	Tags             sql.NullString `db:"tags"`
	ImageAlt         sql.NullString `db:"image_alt"`
	Status           string         `db:"status"` // pending, rewritten, published, failed
	RetryCount       int            `db:"retry_count"`
	NextRetryAt      sql.NullTime   `db:"next_retry_at"`
	PublishDate      sql.NullTime   `db:"publish_date"`
	CreatedAt        time.Time      `db:"created_at"`
}
