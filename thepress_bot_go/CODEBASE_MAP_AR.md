# ThePress Bot Ultimate - Ամբողջական Ճարտարապետական Քարտեզ

Այս փաստաթուղթը պարունակում է բոտի ամբողջական կոդի քարտեզը, ներառյալ ֆայլերը, ֆունկցիաները, կառուցվածքները (structs) և տվյալների բազայի սխեման:

## Տվյալների Բազա (Database)
Բազան SQLite է (`bot_ultimate.db`): Հիմնական աղյուսակներն են.
- `app_settings` - բոտի գլոբալ կարգավորումներ (WP, AI keys, prompts)
- `rss_topics` - RSS հոսքերի ցանկ
- `articles` - հավաքագրված և մշակված հոդվածների պահոց (կարգավիճակներով՝ pending, published, failed)

## Ֆայլերի և Ֆունկցիաների Ցանկ

### `internal/config/config.go` (Package: `unknown`)
**Կառուցվածքներ (Structs):** Config, Topic
**Ֆունկցիաներ:**
- `InitDB`
- `Load`
- `loadFromSQLiteLocked`
- `Save`
- `Get`
- `saveConfigToSQLiteLocked`

### `internal/domain/services/ai_interface.go` (Package: `services`)
**Կառուցվածքներ (Structs):** AIResult
**Ինտերֆեյսներ (Interfaces):** AIProvider

### `internal/domain/services/interfaces.go` (Package: `services`)
**Ինտերֆեյսներ (Interfaces):** Publisher, Linker, ScraperService

### `internal/domain/repository/interfaces.go` (Package: `repository`)
**Ինտերֆեյսներ (Interfaces):** ArticleRepository, FeedRepository

### `internal/domain/models/article.go` (Package: `models`)
**Կառուցվածքներ (Structs):** Article

### `internal/domain/models/feed.go` (Package: `models`)
**Կառուցվածքներ (Structs):** Feed

### `internal/infra/factory.go` (Package: `infra`)
**Կառուցվածքներ (Structs):** DefaultProviderFactory
**Ֆունկցիաներ:**
- `NewDefaultProviderFactory`
- `(f *DefaultProviderFactory) CreateAI`
- `(f *DefaultProviderFactory) CreatePublisher`
- `(f *DefaultProviderFactory) CreateScraper`

### `internal/infra/database/sqlite.go` (Package: `unknown`)
**Ֆունկցիաներ:**
- `NewSQLiteDB`

### `internal/infra/repository/sqlite_article_repo.go` (Package: `repository`)
**Կառուցվածքներ (Structs):** SQLiteArticleRepository
**Ֆունկցիաներ:**
- `NewSQLiteArticleRepository`
- `(r *SQLiteArticleRepository) Save`
- `(r *SQLiteArticleRepository) Update`
- `(r *SQLiteArticleRepository) GetUnprocessed`
- `(r *SQLiteArticleRepository) GetPending`
- `(r *SQLiteArticleRepository) GetFailed`
- `(r *SQLiteArticleRepository) Exists`
- `(r *SQLiteArticleRepository) GetByID`
- `(r *SQLiteArticleRepository) GetOneRewrittenByCategory`
- `(r *SQLiteArticleRepository) GetRelated`
- `(r *SQLiteArticleRepository) GetStats`
- `(r *SQLiteArticleRepository) Delete`
- `(r *SQLiteArticleRepository) ClearQueue`

### `internal/infra/api/server.go` (Package: `api`)
**Կառուցվածքներ (Structs):** Server
**Ֆունկցիաներ:**
- `NewServer`
- `(s *Server) SetBotRunning`
- `(s *Server) handleGetQueueItem`
- `(s *Server) handleUpdateQueueItem`
- `(s *Server) handlePublishItem`
- `(s *Server) handleGetAnalytics`
- `(s *Server) handleGetQueue`
- `(s *Server) handleDeleteItem`
- `(s *Server) handleClearQueue`
- `(s *Server) handleStatus`
- `(s *Server) handleToggle`
- `(s *Server) handleGetSettings`
- `(s *Server) handlePostSettings`
- `(s *Server) Listen`

### `internal/infra/utils/browser.go` (Package: `unknown`)
**Ֆունկցիաներ:**
- `GetRandomUserAgent`
- `OpenBrowser`

### `internal/infra/utils/downloader.go` (Package: `utils`)
**Ֆունկցիաներ:**
- `isSafeURL`
- `getSafeClient`
- `ToBase64`
- `DownloadFileToBytes`
- `DownloadFileToPath`

### `internal/infra/utils/logger.go` (Package: `utils`)
**Կառուցվածքներ (Structs):** LogHub, wsEvent
**Ֆունկցիաներ:**
- `(h *LogHub) Register`
- `(h *LogHub) Unregister`
- `BroadcastEvent`
- `BroadcastLog`

### `internal/infra/utils/sanitizer.go` (Package: `unknown`)
**Ֆունկցիաներ:**
- `SanitizeHTML`
- `CleanJSON`

### `internal/infra/publisher/linker.go` (Package: `publisher`)
**Կառուցվածքներ (Structs):** InternalLinker
**Ֆունկցիաներ:**
- `NewInternalLinker`
- `(l *InternalLinker) InjectLinks`

### `internal/infra/publisher/wordpress.go` (Package: `publisher`)
**Կառուցվածքներ (Structs):** WPClient
**Ֆունկցիաներ:**
- `NewWPClient`
- `(wp *WPClient) doRequest`
- `(wp *WPClient) UploadImageFromURL`
- `(wp *WPClient) UploadImageToMediaLibrary`
- `(wp *WPClient) ProcessInlineImages`
- `(wp *WPClient) Publish`

### `internal/infra/scraper/client.go` (Package: `unknown`)
**Կառուցվածքներ (Structs):** StealthClient
**Ֆունկցիաներ:**
- `GetProxyURL`
- `NewStealthClient`
- `(s *StealthClient) Get`

### `internal/infra/scraper/rss.go` (Package: `scraper`)
**Ֆունկցիաներ:**
- `FetchRSSLinks`

### `internal/infra/scraper/scraper_test.go` (Package: `scraper`)
**Ֆունկցիաներ:**
- `TestURLValidation`

### `internal/infra/scraper/stealth.go` (Package: `unknown`)
**Ֆունկցիաներ:**
- `NewStealthBrowser`
- `PreparePage`

### `internal/infra/scraper/universal.go` (Package: `scraper`)
**Կառուցվածքներ (Structs):** UniversalScraper
**Ֆունկցիաներ:**
- `NewUniversalScraper`
- `(s *UniversalScraper) Scrape`

### `internal/infra/scraper/wrapper.go` (Package: `scraper`)
**Կառուցվածքներ (Structs):** Wrapper
**Ֆունկցիաներ:**
- `NewWrapper`
- `(w *Wrapper) FetchRSSLinks`
- `(w *Wrapper) ScrapeArticle`
- `(w *Wrapper) Close`

### `internal/infra/ai/factory.go` (Package: `unknown`)
**Ֆունկցիաներ:**
- `NewProvider`

### `internal/infra/ai/modelslab.go` (Package: `ai`)
**Կառուցվածքներ (Structs):** ModelsLabProvider
**Ֆունկցիաներ:**
- `NewModelsLabProvider`
- `(m *ModelsLabProvider) ProcessArticle`
- `(m *ModelsLabProvider) ProcessImage`

### `internal/infra/ai/nvidia.go` (Package: `ai`)
**Կառուցվածքներ (Structs):** NvidiaProvider
**Ֆունկցիաներ:**
- `NewNvidiaProvider`
- `(n *NvidiaProvider) ProcessImage`
- `(n *NvidiaProvider) ProcessArticle`
- `(n *NvidiaProvider) runSimpleCompletion`
- `(n *NvidiaProvider) tryModel`
- `(n *NvidiaProvider) executeRequest`
- `(n *NvidiaProvider) Close`

### `internal/usecase/article_usecase.go` (Package: `usecase`)
**Կառուցվածքներ (Structs):** ArticleUseCase
**Ինտերֆեյսներ (Interfaces):** ProviderFactory
**Ֆունկցիաներ:**
- `NewArticleUseCase`
- `(u *ArticleUseCase) ExecuteScrapeCycle`
- `(u *ArticleUseCase) ExecuteAICycle`
- `(u *ArticleUseCase) processWithAI`
- `(u *ArticleUseCase) PublishSingle`
- `(u *ArticleUseCase) ExecutePublishCycle`

### `internal/usecase/article_usecase_test.go` (Package: `usecase`)
**Ֆունկցիաներ:**
- `TestExecutePublishCycle_EmptyTopics`
- `TestContextTimeout`

### `cmd/bot/main.go` (Package: `main`)
**Կառուցվածքներ (Structs):** BotRunner
**Ֆունկցիաներ:**
- `NewBotRunner`
- `(r *BotRunner) SetServer`
- `(r *BotRunner) Start`
- `(r *BotRunner) Stop`
- `main`
- `(r *BotRunner) runScrapeLoop`
- `(r *BotRunner) runAILoop`
- `(r *BotRunner) runPublishLoop`
