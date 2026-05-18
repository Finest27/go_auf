# Ամբողջական Քարտեզ (Complete Architecture Map)

Այս փաստաթուղթը պարունակում է ThePress Bot-ի բոլոր ֆայլերի, ֆունկցիաների, ստրուկտուրաների և տվյալների բազայի մանրամասները։

## Տվյալների Բազա (Database)
- **Տեսակը:** SQLite
- **Ֆայլը:** `bot_ultimate.db`
- **Հիմնական աղյուսակները:** `app_settings`, `rss_topics`, `feeds`, `articles`

## Կոդի Կառուցվածք (Code Structure)

### Թղթապանակ (Directory): `./internal/config`

#### Ֆայլ (File): `config.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `Config` (struct)
- `Topic` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `InitDB()`
- `Load()`
- `loadFromSQLiteLocked()`
- `Save()`
- `Get()`
- `saveConfigToSQLiteLocked()`

### Թղթապանակ (Directory): `./internal/domain/services`

#### Ֆայլ (File): `ai_interface.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `AIResult` (struct)
- `AIProvider` (interface)

#### Ֆայլ (File): `interfaces.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `Publisher` (interface)
- `Linker` (interface)
- `ScraperService` (interface)

### Թղթապանակ (Directory): `./internal/domain/repository`

#### Ֆայլ (File): `interfaces.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `ArticleRepository` (interface)
- `FeedRepository` (interface)

### Թղթապանակ (Directory): `./internal/domain/models`

#### Ֆայլ (File): `feed.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `Feed` (struct)

#### Ֆայլ (File): `article.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `Article` (struct)

### Թղթապանակ (Directory): `./internal/infra`

#### Ֆայլ (File): `factory.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `DefaultProviderFactory` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `NewDefaultProviderFactory()`
- `CreateAI()`
- `CreatePublisher()`
- `CreateScraper()`

### Թղթապանակ (Directory): `./internal/infra/database`

#### Ֆայլ (File): `sqlite.go`
**Ֆունկցիաներ (Functions & Methods):**
- `NewSQLiteDB()`

### Թղթապանակ (Directory): `./internal/infra/repository`

#### Ֆայլ (File): `sqlite_article_repo.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `SQLiteArticleRepository` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `NewSQLiteArticleRepository()`
- `Save()`
- `Update()`
- `GetUnprocessed()`
- `GetPending()`
- `GetFailed()`
- `Exists()`
- `GetByID()`
- `GetOneRewrittenByCategory()`
- `GetRelated()`
- `GetStats()`
- `Delete()`
- `ClearQueue()`

### Թղթապանակ (Directory): `./internal/infra/api`

#### Ֆայլ (File): `server.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `Server` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `NewServer()`
- `SetBotRunning()`
- `handleGetQueueItem()`
- `handleUpdateQueueItem()`
- `handlePublishItem()`
- `handleGetAnalytics()`
- `handleGetQueue()`
- `handleDeleteItem()`
- `handleClearQueue()`
- `handleStatus()`
- `handleToggle()`
- `handleGetSettings()`
- `handlePostSettings()`
- `Listen()`

### Թղթապանակ (Directory): `./internal/infra/utils`

#### Ֆայլ (File): `browser.go`
**Ֆունկցիաներ (Functions & Methods):**
- `GetRandomUserAgent()`
- `OpenBrowser()`

#### Ֆայլ (File): `logger.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `LogHub` (struct)
- `wsEvent` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `Register()`
- `Unregister()`
- `BroadcastEvent()`
- `BroadcastLog()`

#### Ֆայլ (File): `sanitizer.go`
**Ֆունկցիաներ (Functions & Methods):**
- `SanitizeHTML()`
- `CleanJSON()`

#### Ֆայլ (File): `downloader.go`
**Ֆունկցիաներ (Functions & Methods):**
- `isSafeURL()`
- `getSafeClient()`
- `ToBase64()`
- `DownloadFileToBytes()`
- `DownloadFileToPath()`

### Թղթապանակ (Directory): `./internal/infra/publisher`

#### Ֆայլ (File): `wordpress.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `WPClient` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `NewWPClient()`
- `doRequest()`
- `UploadImageFromURL()`
- `UploadImageToMediaLibrary()`
- `ProcessInlineImages()`
- `Publish()`

#### Ֆայլ (File): `linker.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `InternalLinker` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `NewInternalLinker()`
- `InjectLinks()`

### Թղթապանակ (Directory): `./internal/infra/scraper`

#### Ֆայլ (File): `wrapper.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `Wrapper` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `NewWrapper()`
- `FetchRSSLinks()`
- `ScrapeArticle()`
- `Close()`

#### Ֆայլ (File): `universal.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `UniversalScraper` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `NewUniversalScraper()`
- `Scrape()`

#### Ֆայլ (File): `client.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `StealthClient` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `GetProxyURL()`
- `NewStealthClient()`
- `Get()`

#### Ֆայլ (File): `rss.go`
**Ֆունկցիաներ (Functions & Methods):**
- `FetchRSSLinks()`

#### Ֆայլ (File): `scraper_test.go`
**Ֆունկցիաներ (Functions & Methods):**
- `TestURLValidation()`

#### Ֆայլ (File): `stealth.go`
**Ֆունկցիաներ (Functions & Methods):**
- `NewStealthBrowser()`
- `PreparePage()`

### Թղթապանակ (Directory): `./internal/infra/ai`

#### Ֆայլ (File): `nvidia.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `NvidiaProvider` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `NewNvidiaProvider()`
- `ProcessImage()`
- `ProcessArticle()`
- `runSimpleCompletion()`
- `tryModel()`
- `executeRequest()`
- `Close()`

#### Ֆայլ (File): `modelslab.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `ModelsLabProvider` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `NewModelsLabProvider()`
- `ProcessArticle()`
- `ProcessImage()`

#### Ֆայլ (File): `factory.go`
**Ֆունկցիաներ (Functions & Methods):**
- `NewProvider()`

### Թղթապանակ (Directory): `./internal/usecase`

#### Ֆայլ (File): `article_usecase.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `ProviderFactory` (interface)
- `ArticleUseCase` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `NewArticleUseCase()`
- `ExecuteScrapeCycle()`
- `ExecuteAICycle()`
- `processWithAI()`
- `PublishSingle()`
- `ExecutePublishCycle()`

#### Ֆայլ (File): `article_usecase_test.go`
**Ֆունկցիաներ (Functions & Methods):**
- `TestExecutePublishCycle_EmptyTopics()`
- `TestContextTimeout()`

### Թղթապանակ (Directory): `./cmd/bot`

#### Ֆայլ (File): `main.go`
**Ստրուկտուրաներ և Ինտերֆեյսներ (Types):**
- `BotRunner` (struct)
**Ֆունկցիաներ (Functions & Methods):**
- `NewBotRunner()`
- `SetServer()`
- `Start()`
- `Stop()`
- `main()`
- `runScrapeLoop()`
- `runAILoop()`
- `runPublishLoop()`
