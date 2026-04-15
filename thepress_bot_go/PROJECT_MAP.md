# ThePressUSA Բոտի Քարտեզ (Project Map)

Այս փաստաթուղթը պարունակում է նախագծի ամբողջական կառուցվածքը, ֆայլերի նկարագրությունն ու նրանցում առկա ֆունկցիաները, ինչպես նաև SQLite տվյալների բազայի կառուցվածքը։

## 1. Պանակների Կառուցվածքը (Folder Structure)

```text
thepress_bot_go
├── README.md                 - Նախագծի ընդհանուր նկարագրություն
├── bot.exe                   - Windows-ի համար կոմպիլացված գործարկվող ֆայլ
├── bot_live.log              - Log ֆայլ, որտեղ պահպանվում են իրադարձությունները
├── bot_ultimate.db           - Հիմնական տվյալների բազան (SQLite)
├── bot_ultimate.db-shm       - SQLite-ի օժանդակ ֆայլ
├── bot_ultimate.db-wal       - SQLite-ի օժանդակ ֆայլ
├── cmd
│   └── bot
│       └── main.go           - Ծրագրի մուտքի կետը (Entry point)
├── go.mod                    - Go մոդուլների նկարագրություն
├── go.sum                    - Go մոդուլների հեշեր
├── internal
│   ├── config
│   │   └── config.go         - Կարգավորումների և բազայի սինխրոնիզացիա
│   ├── domain
│   │   ├── models
│   │   │   ├── article.go    - Հոդվածի մոդել
│   │   │   └── feed.go       - RSS աղբյուրի մոդել
│   │   ├── repository
│   │   │   └── interfaces.go - Տվյալների բազայի ինտերֆեյսներ
│   │   └── services
│   │       └── ai_interface.go - AI պրովայդերների ինտերֆեյսներ
│   ├── infra
│   │   ├── ai
│   │   │   ├── factory.go    - AI մոդելների ընտրություն
│   │   │   ├── modelslab.go  - ModelsLab API-ի ինտեգրացիա
│   │   │   └── nvidia.go     - Nvidia API-ի ինտեգրացիա
│   │   ├── api
│   │   │   ├── server.go     - Fiber Վեբ սերվեր (UI-ի համար)
│   │   │   └── static        - HTML/CSS/JS (Live UI ֆայլեր)
│   │   ├── database
│   │   │   └── sqlite.go     - SQLite բազայի ստեղծում և կառավարում
│   │   ├── publisher
│   │   │   ├── linker.go     - Ինքնաշխատ ներքին հղումներ (Internal Linker)
│   │   │   └── wordpress.go  - WordPress REST API ինտեգրացիա
│   │   ├── repository
│   │   │   └── sqlite_article_repo.go - DB հարցումներ (CRUD) հոդվածների համար
│   │   ├── scraper
│   │   │   ├── client.go     - HTTP client ռեսուրսների համար
│   │   │   ├── rss.go        - RSS ֆիդերի կարդացում
│   │   │   ├── scraper_test.go - Scraper-ի թեստեր
│   │   │   ├── stealth.go    - Քրոուլինգ (stealth crawling) գործիք
│   │   │   └── universal.go  - GoReadability տեքստի մաքրում
│   │   └── utils
│   │       ├── browser.go    - Բրաուզերի օժանդակ ֆունկցիաներ
│   │       ├── downloader.go - Ֆայլերի անվտանգ ներբեռնում
│   │       ├── logger.go     - WebSocket լոգավորում
│   │       └── sanitizer.go  - HTML-ի անվտանգ մաքրում (XSS պաշտպանություն)
│   └── usecase
│       ├── article_usecase.go      - Բիզնես տրամաբանություն (Scrape & Publish)
│       └── article_usecase_test.go - Usecase-ի թեստեր
└── run_bot.bat               - Գործարկման Windows bat ֆայլ
```

---

## 2. Ֆայլեր և Ֆունկցիաներ

### `cmd/bot/main.go`
Ծրագրի հիմնական գործարկման կետն է:
- `NewBotRunner(uc *usecase.ArticleUseCase) *BotRunner` - Ստեղծում է բոտի աշխատանքի կառավարիչ:
- `Start()` - Սկսում է բոտի աշխատանքը:
- `Stop()` - Կանգնեցնում է բոտը:
- `main()` - Գլխավոր ֆունկցիա:
- `runScrapeLoop(ctx context.Context)` - Նորությունների հավաքագրման (scrape) ցիկլը:
- `runPublishLoop(ctx context.Context)` - Հրապարակման (publish) ցիկլը:

### `internal/config/config.go`
Կարգավորումների վարում:
- `InitDB(db *sqlx.DB)` - Կապում է DB-ն կոնֆիգուրացիայի հետ:
- `Load() error` - Բեռնում է կարգավորումները:
- `loadFromSQLiteLocked() error` - Բեռնում է DB-ից անվտանգ կերպով:
- `Save(cfg Config) error` - Պահպանում է կարգավորումները:
- `Get() Config` - Վերադարձնում է ընթացիկ կարգավորումները:
- `migrateConfigToSQLiteLocked(cfg Config) error` - Միգրացիա դեպի DB:

### `internal/infra/database/sqlite.go`
Տվյալների բազայի կառավարում:
- `NewSQLiteDB(dbPath string) *sqlx.DB` - Ստեղծում է նոր կապ տվյալների բազայի հետ և կիրառում PRAGMA օպտիմիզացիաներ:

### `internal/infra/repository/sqlite_article_repo.go`
Հոդվածների պահպանման և կառավարման ֆունկցիաներ:
- `NewSQLiteArticleRepository(db *sqlx.DB) *SQLiteArticleRepository`
- `Save(ctx context.Context, a *models.Article) error`
- `Update(ctx context.Context, a *models.Article) error`
- `GetPending(ctx context.Context, limit int) ([]models.Article, error)`
- `GetFailed(ctx context.Context, limit int) ([]models.Article, error)`
- `Exists(ctx context.Context, url string) (bool, error)`
- `GetByID(ctx context.Context, id int64) (*models.Article, error)`
- `GetOneRewrittenByCategory(ctx context.Context, catID int) (*models.Article, error)`
- `GetRelated(ctx context.Context, id int64, limit int) ([]models.Article, error)`
- `GetStats(ctx context.Context) (published, pending, failed int, err error)`
- `Delete(ctx context.Context, id int64) error`
- `ClearQueue(ctx context.Context) error`

### `internal/infra/api/server.go`
Բոտի կառավարման վահանակի (Web UI) API-ն:
- `NewServer(onStart, onStop func(), repo *repository.SQLiteArticleRepository, uc *usecase.ArticleUseCase) *Server`
- HTTP Հարցումների սպասարկումներ (handlers)` `handleGetQueueItem`, `handleUpdateQueueItem`, `handlePublishItem`, `handleGetAnalytics`, `handleGetQueue`, `handleDeleteItem`, `handleClearQueue`, `handleStatus`, `handleToggle`, `handleGetSettings`, `handlePostSettings`:
- `Listen(addr string) error` - Սկսում է լսել նշված հասցեն:

### `internal/infra/publisher/wordpress.go`
WordPress կայքում հրապարակելու տրամաբանություն:
- `NewWPClient(...) *WPClient`
- `doRequest(req *http.Request) (*http.Response, error)`
- `UploadImageFromURL(imageURL, altText string) (int, error)`
- `UploadImageToMediaLibrary(imagePath, altText string) (int, string, error)`
- `ProcessInlineImages(html string) (string, error)`
- `Publish(article *models.Article, catID int) (string, error)`

### `internal/infra/publisher/linker.go`
Ներքին հղումների ավտոմատ կցում հոդվածներին:
- `NewInternalLinker(repo repository.ArticleRepository) *InternalLinker`
- `InjectLinks(ctx context.Context, articleID int64, htmlContent string) string`

### `internal/infra/scraper/universal.go` & `stealth.go` & `client.go` & `rss.go`
Նորությունների և RSS-ների հավաքագրում:
- `NewUniversalScraper(b *rod.Browser) *UniversalScraper`
- `Scrape(ctx context.Context, rawURL string) (string, string, string, error)`
- `NewStealthClient() *StealthClient`
- `Get(url string) (string, error)`
- `FetchRSSLinks(ctx context.Context, browser *rod.Browser, url string) ([]string, error)`
- `NewStealthBrowser() (*rod.Browser, error)`
- `PreparePage(b *rod.Browser) *rod.Page`

### `internal/infra/ai/nvidia.go` & `modelslab.go` & `factory.go`
Արհեստական բանականության (AI) ինտեգրացիաներ:
- `NewProvider(ctx context.Context, tool, apiKey, sysPrompt string) (services.AIProvider, error)`
- `NewNvidiaProvider(...) *NvidiaProvider`
- `ProcessImage(...)` և `ProcessArticle(...)` - Նկարների և տեքստերի մշակում
- `NewModelsLabProvider(apiKey string) *ModelsLabProvider`

### `internal/infra/utils/`
Օժանդակ ֆունկցիաներ անվտանգության, լոգավորման և բրաուզերի համար:
- `browser.go` -> `GetRandomUserAgent()`, `OpenBrowser(url string)`
- `logger.go` -> `BroadcastEvent(...)`, `BroadcastLog(...)`
- `sanitizer.go` -> `SanitizeHTML(html string)`, `CleanJSON(raw string)`
- `downloader.go` -> `isSafeURL(rawURL string)`, `DownloadFileToBytes(...)`, `DownloadFileToPath(...)`

### `internal/usecase/article_usecase.go`
Հիմնական բիզնես-տրամաբանությունը, որը միավորում է Scraper-ը, AI-ը և Publisher-ը:
- `NewArticleUseCase(...) *ArticleUseCase`
- `ExecuteScrapeCycle(ctx context.Context, cfg config.Config)`
- `processWithAI(ctx context.Context, art *models.Article, nvidia *ai.NvidiaProvider)`
- `PublishSingle(ctx context.Context, cfg config.Config, art *models.Article)`
- `ExecutePublishCycle(ctx context.Context, cfg config.Config, startIndex int) int`

---

## 3. Տվյալների Բազայի Սխեմա (Database Schema)

Բոտն օգտագործում է SQLite (`bot_ultimate.db`) հետևյալ կառուցվածքով.

### 3.1. Աղյուսակ `app_settings`
Պահում է բոտի կարգավորումները և API բանալիները:
```sql
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
```

### 3.2. Աղյուսակ `rss_topics`
Պահում է RSS աղբյուրները և համապատասխան WordPress Category ID-ները:
```sql
CREATE TABLE IF NOT EXISTS rss_topics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    wp_category_id INTEGER,
    rss_url TEXT UNIQUE
);
```

### 3.3. Աղյուսակ `articles`
Հիմնական աղյուսակը հոդվածների համար:
```sql
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
    status TEXT, -- pending, rewritten, published, failed
    retry_count INTEGER DEFAULT 0,
    next_retry_at DATETIME,
    publish_date DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_articles_status ON articles(status);
CREATE INDEX IF NOT EXISTS idx_articles_source_url ON articles(source_url);
CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at DESC);
```
