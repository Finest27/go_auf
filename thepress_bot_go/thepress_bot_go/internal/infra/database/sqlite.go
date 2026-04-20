package database

import (
	"github.com/jmoiron/sqlx"
	"log"
	_ "modernc.org/sqlite"
)

func NewSQLiteDB(dbPath string) *sqlx.DB {
	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	optimizations := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA foreign_keys = ON;",
	}

	for _, opt := range optimizations {
		if _, err := db.Exec(opt); err != nil {
			log.Printf("[DB WARNING] Failed to set pragma %s: %v", opt, err)
		}
	}

	schema := `
	CREATE TABLE IF NOT EXISTS app_settings (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		wp_url TEXT,
		wp_username TEXT,
		wp_app_password TEXT,
		ai_tool TEXT,
		nvidia_api_key TEXT,
		modelslab_api_key TEXT,
		min_article_len INTEGER,
		run_interval_hours INTEGER,
		run_interval_minutes INTEGER,
		publish_interval_minutes INTEGER,
		auto_publish BOOLEAN,
		system_prompt_nvidia TEXT,
		system_prompt_modelslab TEXT
	);

	CREATE TABLE IF NOT EXISTS rss_topics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		wp_category_id INTEGER,
		rss_url TEXT UNIQUE
	);

	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		source_url TEXT UNIQUE,
		title TEXT,
		content TEXT,
		rewritten_content TEXT,
		image_url TEXT,
		meta_description TEXT,
		focus_keywords TEXT,
		slug TEXT,
		category TEXT,
		category_id INTEGER,
		tags TEXT,
		image_alt TEXT,
		status TEXT,
		retry_count INTEGER DEFAULT 0,
		next_retry_at DATETIME,
		publish_date DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_articles_status ON articles(status);
	CREATE INDEX IF NOT EXISTS idx_articles_source_url ON articles(source_url);
	CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at DESC);
	`
	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	// Schema Updates (Safe Migrations - check error but don't fail if column already exists)
	updates := []string{
		"ALTER TABLE app_settings ADD COLUMN system_prompt_nvidia TEXT;",
		"ALTER TABLE app_settings ADD COLUMN system_prompt_modelslab TEXT;",
		"ALTER TABLE articles ADD COLUMN category_id INTEGER;",
		"ALTER TABLE articles ADD COLUMN retry_count INTEGER DEFAULT 0;",
		"ALTER TABLE articles ADD COLUMN next_retry_at DATETIME;",
	}

	for _, update := range updates {
		_, _ = db.Exec(update) // SQLite doesn't support 'IF NOT EXISTS' for columns, so we ignore errors here
	}

	return db
}
