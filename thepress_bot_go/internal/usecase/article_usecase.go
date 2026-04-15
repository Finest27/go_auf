package usecase

import (
	"context"
	"database/sql"
	"sync"
	"thepress_bot_go/internal/config"
	"thepress_bot_go/internal/domain/models"
	"thepress_bot_go/internal/domain/services"
	"thepress_bot_go/internal/infra/ai"
	"thepress_bot_go/internal/infra/publisher"
	"thepress_bot_go/internal/infra/repository"
	"thepress_bot_go/internal/infra/scraper"
	"thepress_bot_go/internal/infra/utils"
	"time"
)

type ArticleUseCase struct {
	repo   *repository.SQLiteArticleRepository
	linker *publisher.InternalLinker
}

func NewArticleUseCase(repo *repository.SQLiteArticleRepository, linker *publisher.InternalLinker) *ArticleUseCase {
	return &ArticleUseCase{
		repo:   repo,
		linker: linker,
	}
}

func (u *ArticleUseCase) ExecuteScrapeCycle(ctx context.Context, cfg config.Config) {
	nvidia := ai.NewNvidiaProvider(cfg.AI.NvidiaAPIKey, cfg.Prompts.SystemPromptNvidia)

	failedArticles, _ := u.repo.GetFailed(ctx, 10)
	for _, art := range failedArticles {
		utils.BroadcastLog("[RETRY] Փորձում ենք վերամշակել հոդվածը (#%d): %s", art.RetryCount+1, art.Title)
		u.processWithAI(ctx, &art, nvidia)
	}

	for _, topic := range cfg.Topics {
		select {
		case <-ctx.Done(): return
		default:
			utils.BroadcastLog("[SYSTEM] Ստուգում ենք նորություններ: %s", topic.Name)
			browser, err := scraper.NewStealthBrowser()
			if err != nil { continue }

			links, err := scraper.FetchRSSLinks(ctx, browser, topic.RSSURL)
			if err != nil { browser.Close(); continue }

			universal := scraper.NewUniversalScraper(browser)
			var wg sync.WaitGroup
			sem := make(chan struct{}, 3)

			for _, link := range links {
				if exists, _ := u.repo.Exists(ctx, link); exists { continue }

				wg.Add(1)
				sem <- struct{}{}
				go func(targetLink string, catID int) {
					defer wg.Done()
					defer func() { <-sem }()

					scrapeCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
					defer cancel()

					title, content, image, err := universal.Scrape(scrapeCtx, targetLink)
					if err != nil || len(content) < cfg.Advanced.MinArticleLen { return }

					art := &models.Article{
						Title:      title,
						Content:    content,
						SourceURL:  targetLink,
						ImageURL:   sql.NullString{String: image, Valid: image != ""},
						Status:     "pending",
						CategoryID: sql.NullInt64{Int64: int64(catID), Valid: true},
					}
					u.repo.Save(ctx, art)
					u.processWithAI(ctx, art, nvidia)
				}(link, topic.WPCategoryID)
			}
			wg.Wait()
			browser.Close()
		}
	}
}

func (u *ArticleUseCase) processWithAI(ctx context.Context, art *models.Article, nvidia *ai.NvidiaProvider) {
	aiCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	res, err := nvidia.ProcessArticle(aiCtx, art.Title, art.Content)
	if err == nil && res != nil {
		art.Status = "rewritten"
		art.RewrittenContent = sql.NullString{String: res.RewrittenContent, Valid: true}
		art.MetaDescription = sql.NullString{String: res.MetaDescription, Valid: true}
		art.Slug = sql.NullString{String: res.Slug, Valid: true}
		art.ImageAlt = sql.NullString{String: res.ImageAlt, Valid: true}
		art.RetryCount = 0
		art.NextRetryAt = sql.NullTime{Valid: false}
	} else {
		art.Status = "failed"
		art.RetryCount++
		backoff := []time.Duration{10 * time.Minute, 30 * time.Minute, 2 * time.Hour, 6 * time.Hour, 12 * time.Hour}
		delay := backoff[4]
		if art.RetryCount <= len(backoff) {
			delay = backoff[art.RetryCount-1]
		}
		art.NextRetryAt = sql.NullTime{Time: time.Now().Add(delay), Valid: true}
		utils.BroadcastLog("[AI ERROR] %v. Հոդված #%d կփորձվի %v հետո:", err, art.RetryCount, delay)
	}
	u.repo.Update(ctx, art)
}

func (u *ArticleUseCase) PublishSingle(ctx context.Context, cfg config.Config, art *models.Article) {
	nvidia := ai.NewNvidiaProvider(cfg.AI.NvidiaAPIKey, cfg.Prompts.SystemPromptNvidia)
	var imageProvider services.AIProvider
	if cfg.AI.ModelsLabKey != "" {
		imageProvider = ai.NewModelsLabProvider(cfg.AI.ModelsLabKey)
	} else {
		imageProvider = nvidia
	}

	wp := publisher.NewWPClient(cfg.WordPress.URL, cfg.WordPress.Username, cfg.WordPress.AppPassword, imageProvider)
	
	finalHTML := u.linker.InjectLinks(ctx, art.ID, art.RewrittenContent.String)
	art.RewrittenContent.String = finalHTML

	pubLink, pubErr := wp.Publish(art, int(art.CategoryID.Int64))
	if pubErr != nil {
		art.Status = "failed"
		art.RetryCount++
		art.NextRetryAt = sql.NullTime{Time: time.Now().Add(30 * time.Minute), Valid: true}
		u.repo.Update(ctx, art)
		utils.BroadcastLog("[PUBLISHER ERROR] %v", pubErr)
		return
	}

	art.Status = "published"
	art.PublishDate = sql.NullTime{Time: time.Now(), Valid: true}
	u.repo.Update(ctx, art)
	utils.BroadcastLog("[SUCCESS] Հրապարակված է: %s", pubLink)
}

func (u *ArticleUseCase) ExecutePublishCycle(ctx context.Context, cfg config.Config, startIndex int) int {
	if len(cfg.Topics) == 0 { return 0 }

	for i := 0; i < len(cfg.Topics); i++ {
		indexToCheck := (startIndex + i) % len(cfg.Topics)
		topic := cfg.Topics[indexToCheck]

		art, err := u.repo.GetOneRewrittenByCategory(ctx, topic.WPCategoryID)
		if err == nil && art != nil {
			u.PublishSingle(ctx, cfg, art)
			return (indexToCheck + 1) % len(cfg.Topics)
		}
	}
	return startIndex
}
