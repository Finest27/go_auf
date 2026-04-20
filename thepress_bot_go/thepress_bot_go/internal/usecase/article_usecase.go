package usecase

import (
	"context"
	"database/sql"
		"thepress_bot_go/internal/config"
	"thepress_bot_go/internal/domain/models"
	"thepress_bot_go/internal/domain/repository"
	"thepress_bot_go/internal/domain/services"
	"thepress_bot_go/internal/infra/utils"
	"time"
)

type ProviderFactory interface {
	CreateAI(cfg config.Config) services.AIProvider
	CreatePublisher(cfg config.Config) services.Publisher
	CreateScraper() (services.ScraperService, error)
}

type ArticleUseCase struct {
	repo    repository.ArticleRepository
	linker  services.Linker
	factory ProviderFactory
}

func NewArticleUseCase(repo repository.ArticleRepository, linker services.Linker, factory ProviderFactory) *ArticleUseCase {
	return &ArticleUseCase{
		repo:    repo,
		linker:  linker,
		factory: factory,
	}
}

func (u *ArticleUseCase) ExecuteScrapeCycle(ctx context.Context, cfg config.Config, startIndex int) int {
        _, pendingRewritten, failedCount, err := u.repo.GetStats(ctx)
        if err == nil {
                if failedCount > 0 {
                        failedArticles, err := u.repo.GetFailed(ctx, 1)
                        if err == nil && len(failedArticles) > 0 {
                                art := failedArticles[0]
                                utils.BroadcastLog("[RETRY] Կրկնակի փորձ վերաշարադրելու հոդվածը (#%d): %s", art.RetryCount+1, art.Title)
                                u.processWithAI(ctx, &art, u.factory.CreateAI(cfg))
                                return startIndex
                        }
                }

                if pendingRewritten > 0 {
                        utils.BroadcastLog("[SCRAPE] Կա սպասող հոդված (%d հատ): Սպասում ենք հրապարակմանը...", pendingRewritten)
                        return startIndex
                }
        }

        aiProv := u.factory.CreateAI(cfg)

        if len(cfg.Topics) == 0 {
                return 0
        }

        for i := 0; i < len(cfg.Topics); i++ {
                indexToCheck := (startIndex + i) % len(cfg.Topics)
                topic := cfg.Topics[indexToCheck]

                select {
                case <-ctx.Done():
                        return startIndex
                default:
                        utils.BroadcastLog("[SYSTEM] Ստուգում ենք RSS թեման: %s", topic.Name) 
                        scraperService, err := u.factory.CreateScraper()
                        if err != nil {
                                utils.BroadcastLog("[SYSTEM ERROR] Failed to create scraper: %v", err)
                                continue
                        }

                        processedOne := false
                        func() {
                                defer scraperService.Close()

                                links, err := scraperService.FetchRSSLinks(ctx, topic.RSSURL) 
                                if err != nil {
                                        utils.BroadcastLog("[SYSTEM ERROR] Failed to fetch links for %s: %v", topic.Name, err)
                                        return
                                }

                                for _, link := range links {
                                        exists, err := u.repo.Exists(ctx, link)
                                        if err != nil {
                                                utils.BroadcastLog("[SYSTEM ERROR] DB error checking exists: %v", err)
                                                continue
                                        }
                                        if exists {
                                                continue
                                        }

                                        scrapeCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
                                        title, content, image, err := scraperService.ScrapeArticle(scrapeCtx, link)
                                        cancel()

                                        if err != nil {
                                                continue
                                        }

                                        if len(content) < cfg.Advanced.MinArticleLen {        
                                                continue
                                        }

                                        art := &models.Article{
                                                Title:      title,
                                                Content:    content,
                                                SourceURL:  link,
                                                ImageURL:   sql.NullString{String: image, Valid: image != ""},
                                                Status:     "pending",
                                                CategoryID: sql.NullInt64{Int64: int64(topic.WPCategoryID), Valid: true},
                                        }
                                        if err := u.repo.Save(ctx, art); err != nil {
                                                utils.BroadcastLog("[SYSTEM ERROR] Failed to save article: %v", err)
                                                return
                                        }

                                        u.processWithAI(ctx, art, aiProv)
                                        processedOne = true
                                        return
                                }
                        }()

                        if processedOne {
                                return (indexToCheck + 1) % len(cfg.Topics)
                        }
                }
        }
        return startIndex
}

func (u *ArticleUseCase) processWithAI(ctx context.Context, art *models.Article, aiProv services.AIProvider) {
	aiCtx, cancel := context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	res, err := aiProv.ProcessArticle(aiCtx, art.Title, art.Content)
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
		utils.BroadcastLog("[AI ERROR] %v. Հոդված #%d կփորձվի %v հետո", err, art.RetryCount, delay)
	}
	if err := u.repo.Update(ctx, art); err != nil {
		utils.BroadcastLog("[SYSTEM ERROR] Failed to update article AI status: %v", err)
	}
}

func (u *ArticleUseCase) PublishSingle(ctx context.Context, cfg config.Config, art *models.Article) {
	pub := u.factory.CreatePublisher(cfg)

	finalHTML := u.linker.InjectLinks(ctx, art.ID, art.RewrittenContent.String)
	art.RewrittenContent.String = finalHTML

	pubLink, pubErr := pub.Publish(art, int(art.CategoryID.Int64))
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
	if err := u.repo.Update(ctx, art); err != nil {
		utils.BroadcastLog("[SYSTEM ERROR] Failed to update article publish status: %v", err)
	}
	utils.BroadcastLog("[SUCCESS] Հրապարակված է: %s", pubLink)
}

func (u *ArticleUseCase) ExecutePublishCycle(ctx context.Context, cfg config.Config, startIndex int) int {
	if len(cfg.Topics) == 0 {
		return 0
	}

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
