package config

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	WordPress struct {
		URL         string `json:"url" db:"wp_url"`
		Username    string `json:"username" db:"wp_username"`
		AppPassword string `json:"app_password" db:"wp_app_password"`
	} `json:"wordpress"`
	AI struct {
		Tool         string `json:"tool" db:"ai_tool"`
		NvidiaAPIKey string `json:"nvidia_api_key" db:"nvidia_api_key"`
		ModelsLabKey string `json:"modelslab_api_key" db:"modelslab_api_key"`
	} `json:"ai"`
	Advanced struct {
		MinArticleLen int `json:"min_article_len" db:"min_article_len"`
	} `json:"advanced"`
	Bot struct {
		RunIntervalHours       int  `json:"run_interval_hours" db:"run_interval_hours"`
		RunIntervalMinutes     int  `json:"run_interval_minutes" db:"run_interval_minutes"`
		PublishIntervalMinutes int  `json:"publish_interval_minutes" db:"publish_interval_minutes"`
		AutoPublish            bool `json:"auto_publish" db:"auto_publish"`
	} `json:"bot"`
	Topics  []Topic `json:"topics"`
	Prompts struct {
		SystemPromptNvidia    string `json:"system_prompt_nvidia" db:"system_prompt_nvidia"`
		SystemPromptModelsLab string `json:"system_prompt_modelslab" db:"system_prompt_modelslab"`
	} `json:"prompts"`
}

type Topic struct {
	Name         string `json:"name" db:"name"`
	WPCategoryID int    `json:"wp_category_id" db:"wp_category_id"`
	RSSURL       string `json:"rss_url" db:"rss_url"`
}

var (
	GlobalConfig Config
	mu           sync.RWMutex
	dbConn       *sqlx.DB
)

func InitDB(db *sqlx.DB) {
	mu.Lock()
	dbConn = db
	mu.Unlock()
}

func Load() error {
	mu.Lock()
	defer mu.Unlock()

	if dbConn == nil {
		return nil
	}

	var dbCount int
	err := dbConn.Get(&dbCount, "SELECT COUNT(*) FROM app_settings")
	if err == nil && dbCount > 0 {
		return loadFromSQLiteLocked()
	}

	// Fallback to settings.json only if DB is empty
	if _, err := os.Stat("settings.json"); err == nil {
		data, _ := os.ReadFile("settings.json")
		if err := json.Unmarshal(data, &GlobalConfig); err == nil {
			_ = migrateConfigToSQLiteLocked(GlobalConfig)
			os.Rename("settings.json", "settings.json.bak")
			return nil
		}
	}

	// Default settings if nothing exists
	GlobalConfig.Prompts.SystemPromptNvidia = "You are a professional journalist. Rewrite the following article in Armenian, making it engaging and SEO-friendly."
	return nil
}

func loadFromSQLiteLocked() error {
	type flatSettings struct {
		WpUrl                  string `db:"wp_url"`
		WpUsername             string `db:"wp_username"`
		WpAppPassword          string `db:"wp_app_password"`
		AiTool                 string `db:"ai_tool"`
		NvidiaApiKey           string `db:"nvidia_api_key"`
		ModelslabApiKey        string `db:"modelslab_api_key"`
		MinArticleLen          int    `db:"min_article_len"`
		RunIntervalHours       int    `db:"run_interval_hours"`
		RunIntervalMinutes     int    `db:"run_interval_minutes"`
		PublishIntervalMinutes int    `db:"publish_interval_minutes"`
		AutoPublish            bool   `db:"auto_publish"`
		SystemPromptNvidia    string `db:"system_prompt_nvidia"`
		SystemPromptModelsLab string `db:"system_prompt_modelslab"`
	}

	var flat flatSettings
	query := "SELECT wp_url, wp_username, wp_app_password, ai_tool, nvidia_api_key, modelslab_api_key, min_article_len, run_interval_hours, run_interval_minutes, publish_interval_minutes, auto_publish, system_prompt_nvidia, system_prompt_modelslab FROM app_settings LIMIT 1"

	err := dbConn.Get(&flat, query)
	if err != nil {
		return err
	}

	GlobalConfig.WordPress.URL = flat.WpUrl
	GlobalConfig.WordPress.Username = flat.WpUsername
	GlobalConfig.WordPress.AppPassword = flat.WpAppPassword
	GlobalConfig.AI.Tool = flat.AiTool
	GlobalConfig.AI.NvidiaAPIKey = flat.NvidiaApiKey
	GlobalConfig.AI.ModelsLabKey = flat.ModelslabApiKey
	GlobalConfig.Advanced.MinArticleLen = flat.MinArticleLen
	GlobalConfig.Bot.RunIntervalHours = flat.RunIntervalHours
	GlobalConfig.Bot.RunIntervalMinutes = flat.RunIntervalMinutes
	GlobalConfig.Bot.PublishIntervalMinutes = flat.PublishIntervalMinutes
	GlobalConfig.Bot.AutoPublish = flat.AutoPublish
	GlobalConfig.Prompts.SystemPromptNvidia = flat.SystemPromptNvidia
	GlobalConfig.Prompts.SystemPromptModelsLab = flat.SystemPromptModelsLab

	var topics []Topic
	_ = dbConn.Select(&topics, "SELECT name, wp_category_id, rss_url FROM rss_topics")
	GlobalConfig.Topics = topics

	return nil
}

func Save(cfg Config) error {
	mu.Lock()
	defer mu.Unlock()

	// Preserving sensitive data if not provided in update
	if cfg.WordPress.AppPassword == "" { cfg.WordPress.AppPassword = GlobalConfig.WordPress.AppPassword }
	if cfg.AI.NvidiaAPIKey == "" { cfg.AI.NvidiaAPIKey = GlobalConfig.AI.NvidiaAPIKey }
	if cfg.AI.ModelsLabKey == "" { cfg.AI.ModelsLabKey = GlobalConfig.AI.ModelsLabKey }

	if dbConn != nil {
		err := migrateConfigToSQLiteLocked(cfg)
		if err != nil { return err }
	}

	GlobalConfig = cfg
	return nil
}

func Get() Config {
	mu.RLock()
	defer mu.RUnlock()
	return GlobalConfig
}

func migrateConfigToSQLiteLocked(cfg Config) error {
	query := "INSERT INTO app_settings (id, wp_url, wp_username, wp_app_password, ai_tool, nvidia_api_key, modelslab_api_key, min_article_len, run_interval_hours, run_interval_minutes, publish_interval_minutes, auto_publish, system_prompt_nvidia, system_prompt_modelslab) VALUES (1, :wp_url, :wp_username, :wp_app_password, :ai_tool, :nvidia_api_key, :modelslab_api_key, :min_article_len, :run_interval_hours, :run_interval_minutes, :publish_interval_minutes, :auto_publish, :system_prompt_nvidia, :system_prompt_modelslab) ON CONFLICT(id) DO UPDATE SET wp_url=excluded.wp_url, wp_username=excluded.wp_username, wp_app_password=excluded.wp_app_password, ai_tool=excluded.ai_tool, nvidia_api_key=excluded.nvidia_api_key, modelslab_api_key=excluded.modelslab_api_key, min_article_len=excluded.min_article_len, run_interval_hours=excluded.run_interval_hours, run_interval_minutes=excluded.run_interval_minutes, publish_interval_minutes=excluded.publish_interval_minutes, auto_publish=excluded.auto_publish, system_prompt_nvidia=excluded.system_prompt_nvidia, system_prompt_modelslab=excluded.system_prompt_modelslab;"

	flat := map[string]interface{}{
		"wp_url":                   cfg.WordPress.URL,
		"wp_username":              cfg.WordPress.Username,
		"wp_app_password":          cfg.WordPress.AppPassword,
		"ai_tool":                  cfg.AI.Tool,
		"nvidia_api_key":           cfg.AI.NvidiaAPIKey,
		"modelslab_api_key":        cfg.AI.ModelsLabKey,
		"min_article_len":          cfg.Advanced.MinArticleLen,
		"run_interval_hours":       cfg.Bot.RunIntervalHours,
		"run_interval_minutes":     cfg.Bot.RunIntervalMinutes,
		"publish_interval_minutes": cfg.Bot.PublishIntervalMinutes,
		"auto_publish":             cfg.Bot.AutoPublish,
		"system_prompt_nvidia":    cfg.Prompts.SystemPromptNvidia,
		"system_prompt_modelslab": cfg.Prompts.SystemPromptModelsLab,
	}

	tx, err := dbConn.Beginx()
	if err != nil { return err }

	_, err = tx.NamedExec(query, flat)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, _ = tx.Exec("DELETE FROM rss_topics")
	for _, topic := range cfg.Topics {
		_, _ = tx.NamedExec("INSERT INTO rss_topics (name, wp_category_id, rss_url) VALUES (:name, :wp_category_id, :rss_url)", topic)
	}

	return tx.Commit()
}
