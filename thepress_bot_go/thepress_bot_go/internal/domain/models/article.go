package models

import (
	"database/sql"
	"time"
)

type Article struct {
	ID               int64          `db:"id" json:"id"`
	SourceURL        string         `db:"source_url" json:"source_url"`
	Title            string         `db:"title" json:"title"`
	Content          string         `db:"content" json:"content"`
	RewrittenContent sql.NullString `db:"rewritten_content" json:"rewritten_content"`
	ImageURL         sql.NullString `db:"image_url" json:"image_url"`
	MetaDescription  sql.NullString `db:"meta_description" json:"meta_description"`
	FocusKeywords    sql.NullString `db:"focus_keywords" json:"focus_keywords"`
	Slug             sql.NullString `db:"slug" json:"slug"`
	Category         sql.NullString `db:"category" json:"category"`
	CategoryID       sql.NullInt64  `db:"category_id" json:"category_id"`
	Tags             sql.NullString `db:"tags" json:"tags"`
	ImageAlt         sql.NullString `db:"image_alt" json:"image_alt"`
	Status           string         `db:"status" json:"status"`
	RetryCount       int            `db:"retry_count" json:"retry_count"`
	NextRetryAt      sql.NullTime   `db:"next_retry_at" json:"next_retry_at"`
	PublishDate      sql.NullTime   `db:"publish_date" json:"publish_date"`
	CreatedAt        time.Time      `db:"created_at" json:"created_at"`
}
