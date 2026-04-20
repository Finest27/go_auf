package repository

import (
	"context"
	"sync"
	"github.com/jmoiron/sqlx"
	"thepress_bot_go/internal/domain/models"
	"thepress_bot_go/internal/infra/utils"
)

type SQLiteArticleRepository struct {
	db *sqlx.DB
	mu sync.RWMutex // Use RWMutex for better concurrency
}

func NewSQLiteArticleRepository(db *sqlx.DB) *SQLiteArticleRepository {
	return &SQLiteArticleRepository{db: db}
}

func (r *SQLiteArticleRepository) Save(ctx context.Context, a *models.Article) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Using RETURNING id (supported in modern SQLite) to avoid extra SELECT
	query := `INSERT INTO articles (source_url, title, content, status, image_url, category_id, retry_count)
                  VALUES (?, ?, ?, ?, ?, ?, 0) ON CONFLICT(source_url) DO NOTHING RETURNING id`
	
	err := r.db.GetContext(ctx, &a.ID, query, a.SourceURL, a.Title, a.Content, a.Status, a.ImageURL, a.CategoryID)
	if err != nil {
		// If nothing was inserted (DO NOTHING), we might need to get the existing ID if it is expected
		if err.Error() == "sql: no rows in result set" {
			return nil 
		}
		return err
	}

	go utils.BroadcastEvent("queue_update", nil)
	return nil
}

func (r *SQLiteArticleRepository) Update(ctx context.Context, a *models.Article) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	query := `UPDATE articles SET status = ?, rewritten_content = ?, meta_description = ?, focus_keywords = ?,
                slug = ?, category = ?, tags = ?, image_alt = ?, category_id = ?, retry_count = ?, next_retry_at = ?, publish_date = ?
                WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, a.Status, a.RewrittenContent, a.MetaDescription, a.FocusKeywords,
		a.Slug, a.Category, a.Tags, a.ImageAlt, a.CategoryID, a.RetryCount, a.NextRetryAt, a.PublishDate, a.ID)

	if err == nil {
		go utils.BroadcastEvent("queue_update", nil)
	}

	return err
}

func (r *SQLiteArticleRepository) GetUnprocessed(ctx context.Context, limit int) ([]models.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var articles []models.Article
	err := r.db.SelectContext(ctx, &articles, "SELECT * FROM articles WHERE status = 'pending' ORDER BY id ASC LIMIT ?", limit)
	return articles, err
}

func (r *SQLiteArticleRepository) GetPending(ctx context.Context, limit int) ([]models.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var articles []models.Article
	err := r.db.SelectContext(ctx, &articles, "SELECT * FROM articles WHERE status IN ('pending', 'rewritten', 'failed') ORDER BY id DESC LIMIT ?", limit)
	return articles, err
}

func (r *SQLiteArticleRepository) GetFailed(ctx context.Context, limit int) ([]models.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var articles []models.Article
	query := `SELECT * FROM articles WHERE status = 'failed' AND retry_count < 5
                AND (next_retry_at IS NULL OR next_retry_at <= CURRENT_TIMESTAMP)
                ORDER BY retry_count ASC, id DESC LIMIT ?`
	err := r.db.SelectContext(ctx, &articles, query, limit)
	return articles, err
}

func (r *SQLiteArticleRepository) Exists(ctx context.Context, url string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM articles WHERE source_url = ?", url)
	return count > 0, err
}

func (r *SQLiteArticleRepository) GetByID(ctx context.Context, id int64) (*models.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var a models.Article
	err := r.db.GetContext(ctx, &a, "SELECT * FROM articles WHERE id = ?", id)
	return &a, err
}

func (r *SQLiteArticleRepository) GetOneRewrittenByCategory(ctx context.Context, catID int) (*models.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var a models.Article
	err := r.db.GetContext(ctx, &a, "SELECT * FROM articles WHERE status = 'rewritten' AND category_id = ? ORDER BY id ASC LIMIT 1", catID)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *SQLiteArticleRepository) GetRelated(ctx context.Context, id int64, limit int) ([]models.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var articles []models.Article
	err := r.db.SelectContext(ctx, &articles, "SELECT * FROM articles WHERE id != ? AND status = 'published' ORDER BY created_at DESC LIMIT ?", id, limit)
	return articles, err
}

func (r *SQLiteArticleRepository) GetStats(ctx context.Context) (published, pending, failed int, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	err = r.db.GetContext(ctx, &published, "SELECT COUNT(*) FROM articles WHERE status = 'published'")
	if err != nil { return }
	err = r.db.GetContext(ctx, &pending, "SELECT COUNT(*) FROM articles WHERE status IN ('pending', 'rewritten')")
	if err != nil { return }
	err = r.db.GetContext(ctx, &failed, "SELECT COUNT(*) FROM articles WHERE status = 'failed'")
	return
}

func (r *SQLiteArticleRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, "DELETE FROM articles WHERE id = ?", id)

	if err == nil {
		go utils.BroadcastEvent("queue_update", nil)
	}

	return err
}

func (r *SQLiteArticleRepository) ClearQueue(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, "DELETE FROM articles WHERE status IN ('pending', 'rewritten', 'failed')")

	if err == nil {
		go utils.BroadcastEvent("queue_update", nil)
	}

	return err
}