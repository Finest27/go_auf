# ThePressUSA Auto-Journalist Bot - Complete Architecture Map

## Codebase Structure & Components

### File: `./cmd/bot/main.go`
```go
package main

type BotRunner struct {
func NewBotRunner(uc *usecase.ArticleUseCase) *BotRunner
func (r *BotRunner) SetServer(s *api.Server)
func (r *BotRunner) Start()
func (r *BotRunner) Stop()
func main()
func (r *BotRunner) runScrapeLoop(ctx context.Context)
func (r *BotRunner) runAILoop(ctx context.Context)
func (r *BotRunner) runPublishLoop(ctx context.Context)
```

### File: `./internal/config/config.go`
```go

type Config struct {
type Topic struct {
func InitDB(db *sqlx.DB)
func Load() error
func loadFromSQLiteLocked() error
func Save(cfg Config) error
func Get() Config
func saveConfigToSQLiteLocked(cfg Config) error
```

### File: `./internal/domain/models/article.go`
```go
package models

type Article struct {
```

### File: `./internal/domain/models/feed.go`
```go
package models

type Feed struct {
```

### File: `./internal/domain/repository/interfaces.go`
```go
package repository

type ArticleRepository interface {
type FeedRepository interface {
```

### File: `./internal/domain/services/ai_interface.go`
```go
package services

type AIProvider interface {
type AIResult struct {
```

### File: `./internal/domain/services/interfaces.go`
```go
package services

type Publisher interface {
type Linker interface {
type ScraperService interface {
```

### File: `./internal/infra/ai/factory.go`
```go

func NewProvider(ctx context.Context, tool, apiKey, sysPrompt string) (services.AIProvider, error)
```

### File: `./internal/infra/ai/modelslab.go`
```go
package ai

type ModelsLabProvider struct {
func NewModelsLabProvider(apiKey string) *ModelsLabProvider
func (m *ModelsLabProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error)
func (m *ModelsLabProvider) ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error)
```

### File: `./internal/infra/ai/nvidia.go`
```go
package ai

type NvidiaProvider struct {
func NewNvidiaProvider(apiKey, sysPrompt string) *NvidiaProvider
func (n *NvidiaProvider) ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error)
func (n *NvidiaProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error)
func (n *NvidiaProvider) runSimpleCompletion(ctx context.Context, url, model, prompt string) (string, error)
func (n *NvidiaProvider) tryModel(ctx context.Context, url, model, prompt string) (*services.AIResult, error)
func (n *NvidiaProvider) executeRequest(ctx context.Context, url string, bodyData interface{}) (string, error)
func (n *NvidiaProvider) Close()
```

### File: `./internal/infra/api/server.go`
```go
package api

type Server struct {
func NewServer(onStart, onStop func(), repo *repository.SQLiteArticleRepository, uc *usecase.ArticleUseCase) *Server
func (s *Server) SetBotRunning(running bool)
func (s *Server) handleGetQueueItem(c *fiber.Ctx) error
func (s *Server) handleUpdateQueueItem(c *fiber.Ctx) error
func (s *Server) handlePublishItem(c *fiber.Ctx) error
func (s *Server) handleGetAnalytics(c *fiber.Ctx) error
func (s *Server) handleGetQueue(c *fiber.Ctx) error
func (s *Server) handleDeleteItem(c *fiber.Ctx) error
func (s *Server) handleClearQueue(c *fiber.Ctx) error
func (s *Server) handleStatus(c *fiber.Ctx) error
func (s *Server) handleToggle(c *fiber.Ctx) error
func (s *Server) handleGetSettings(c *fiber.Ctx) error
func (s *Server) handlePostSettings(c *fiber.Ctx) error
func (s *Server) Listen(addr string) error
```

### File: `./internal/infra/database/sqlite.go`
```go

func NewSQLiteDB(dbPath string) *sqlx.DB
```

### File: `./internal/infra/factory.go`
```go
package infra

type DefaultProviderFactory struct{}
func NewDefaultProviderFactory() *DefaultProviderFactory
func (f *DefaultProviderFactory) CreateAI(cfg config.Config) services.AIProvider
func (f *DefaultProviderFactory) CreatePublisher(cfg config.Config) services.Publisher
func (f *DefaultProviderFactory) CreateScraper() (services.ScraperService, error)
```

### File: `./internal/infra/publisher/linker.go`
```go
package publisher

type InternalLinker struct {
func NewInternalLinker(repo repository.ArticleRepository) *InternalLinker
func (l *InternalLinker) InjectLinks(ctx context.Context, articleID int64, htmlContent string) string
```

### File: `./internal/infra/publisher/wordpress.go`
```go
package publisher

type WPClient struct {
func NewWPClient(url, user, pass string, ai services.AIProvider) *WPClient
func (wp *WPClient) doRequest(req *http.Request) (*http.Response, error)
func (wp *WPClient) UploadImageFromURL(imageURL, altText string) (int, error)
func (wp *WPClient) UploadImageToMediaLibrary(imagePath, altText string) (int, string, error)
func (wp *WPClient) ProcessInlineImages(html string) (string, error)
func (wp *WPClient) Publish(article *models.Article, catID int) (string, error)
```

### File: `./internal/infra/repository/sqlite_article_repo.go`
```go
package repository

type SQLiteArticleRepository struct {
func NewSQLiteArticleRepository(db *sqlx.DB) *SQLiteArticleRepository
func (r *SQLiteArticleRepository) Save(ctx context.Context, a *models.Article) error
func (r *SQLiteArticleRepository) Update(ctx context.Context, a *models.Article) error
func (r *SQLiteArticleRepository) GetUnprocessed(ctx context.Context, limit int) ([]models.Article, error)
func (r *SQLiteArticleRepository) GetPending(ctx context.Context, limit int) ([]models.Article, error)
func (r *SQLiteArticleRepository) GetFailed(ctx context.Context, limit int) ([]models.Article, error)
func (r *SQLiteArticleRepository) Exists(ctx context.Context, url string) (bool, error)
func (r *SQLiteArticleRepository) GetByID(ctx context.Context, id int64) (*models.Article, error)
func (r *SQLiteArticleRepository) GetOneRewrittenByCategory(ctx context.Context, catID int) (*models.Article, error)
func (r *SQLiteArticleRepository) GetRelated(ctx context.Context, id int64, limit int) ([]models.Article, error)
func (r *SQLiteArticleRepository) GetStats(ctx context.Context) (published, pending, failed int, err error)
func (r *SQLiteArticleRepository) Delete(ctx context.Context, id int64) error
func (r *SQLiteArticleRepository) ClearQueue(ctx context.Context) error
```

### File: `./internal/infra/scraper/client.go`
```go

type StealthClient struct {
func GetProxyURL() *url.URL
func NewStealthClient() *StealthClient
func (s *StealthClient) Get(url string) (string, error)
```

### File: `./internal/infra/scraper/rss.go`
```go
package scraper

func FetchRSSLinks(ctx context.Context, browser *rod.Browser, url string) ([]string, error)
```

### File: `./internal/infra/scraper/scraper_test.go`
```go
package scraper

func TestURLValidation(t *testing.T)
```

### File: `./internal/infra/scraper/stealth.go`
```go

func NewStealthBrowser() (*rod.Browser, error)
func PreparePage(b *rod.Browser) *rod.Page
```

### File: `./internal/infra/scraper/universal.go`
```go
package scraper

type UniversalScraper struct {
func NewUniversalScraper(b *rod.Browser) *UniversalScraper
func (s *UniversalScraper) Scrape(ctx context.Context, rawURL string) (string, string, string, error)
```

### File: `./internal/infra/scraper/wrapper.go`
```go
package scraper

type Wrapper struct {
func NewWrapper() (services.ScraperService, error)
func (w *Wrapper) FetchRSSLinks(ctx context.Context, url string) ([]string, error)
func (w *Wrapper) ScrapeArticle(ctx context.Context, url string) (title string, content string, image string, err error)
func (w *Wrapper) Close()
```

### File: `./internal/infra/utils/browser.go`
```go

func GetRandomUserAgent() string
func OpenBrowser(url string) error
```

### File: `./internal/infra/utils/downloader.go`
```go
package utils

func isSafeURL(rawURL string) error
func getSafeClient() *http.Client
func ToBase64(data []byte) string
func DownloadFileToBytes(url string) ([]byte, error)
func DownloadFileToPath(url, path string) error
```

### File: `./internal/infra/utils/logger.go`
```go
package utils

type LogHub struct {
type wsEvent struct {
func (h *LogHub) Register() chan string
func (h *LogHub) Unregister(ch chan string)
func BroadcastEvent(eventType string, data interface{})
func BroadcastLog(format string, v ...interface{})
```

### File: `./internal/infra/utils/sanitizer.go`
```go

func SanitizeHTML(html string) string
func CleanJSON(raw string) string
```

### File: `./internal/usecase/article_usecase.go`
```go
package usecase

type ProviderFactory interface {
type ArticleUseCase struct {
func NewArticleUseCase(repo repository.ArticleRepository, linker services.Linker, factory ProviderFactory) *ArticleUseCase
func (u *ArticleUseCase) ExecuteScrapeCycle(ctx context.Context, cfg config.Config, startIndex int) int
func (u *ArticleUseCase) ExecuteAICycle(ctx context.Context, cfg config.Config)
func (u *ArticleUseCase) processWithAI(ctx context.Context, art *models.Article, aiProv services.AIProvider)
func (u *ArticleUseCase) PublishSingle(ctx context.Context, cfg config.Config, art *models.Article)
func (u *ArticleUseCase) ExecutePublishCycle(ctx context.Context, cfg config.Config, startIndex int) int
```

### File: `./internal/usecase/article_usecase_test.go`
```go
package usecase

func TestExecutePublishCycle_EmptyTopics(t *testing.T)
func TestContextTimeout(t *testing.T)
```
