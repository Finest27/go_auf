# Complete Architecture Map
## Database Schema
CREATE TABLE feeds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT UNIQUE,
		name TEXT,
		last_check DATETIME,
		active BOOLEAN DEFAULT 1
	);
CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE articles (
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
		tags TEXT,
		status TEXT,
		publish_date DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	, category_id INTEGER, image_alt TEXT, retry_count INTEGER DEFAULT 0, next_retry_at DATETIME);
CREATE INDEX idx_articles_status ON articles(status);
CREATE INDEX idx_articles_source_url ON articles(source_url);
CREATE INDEX idx_articles_created_at ON articles(created_at DESC);
CREATE TABLE app_settings (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		wp_url TEXT,
		wp_username TEXT,
		wp_app_password TEXT,
		ai_tool TEXT,
		nvidia_api_key TEXT,
		gemini_api_key TEXT,
		modelslab_api_key TEXT,
		min_article_len INTEGER,
		run_interval_hours INTEGER,
		run_interval_minutes INTEGER,
		publish_interval_minutes INTEGER,
		auto_publish BOOLEAN,
		system_prompt TEXT
	, system_prompt_gemini TEXT, system_prompt_nvidia TEXT, system_prompt_modelslab TEXT);
CREATE TABLE rss_topics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		wp_category_id INTEGER,
		rss_url TEXT UNIQUE
	);

## Directory Structure
.
├── bot.exe
├── bot.exe~
├── bot_architecture_map.md
├── bot_state.db
├── bot_ultimate.db
├── cmd
│   └── bot
│       └── main.go
├── complete_map.md
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── domain
│   │   ├── models
│   │   │   ├── article.go
│   │   │   └── feed.go
│   │   ├── repository
│   │   │   └── interfaces.go
│   │   └── services
│   │       ├── ai_interface.go
│   │       └── interfaces.go
│   ├── infra
│   │   ├── ai
│   │   │   ├── factory.go
│   │   │   ├── modelslab.go
│   │   │   └── nvidia.go
│   │   ├── api
│   │   │   ├── server.go
│   │   │   └── static
│   │   ├── database
│   │   │   └── sqlite.go
│   │   ├── factory.go
│   │   ├── publisher
│   │   │   ├── linker.go
│   │   │   └── wordpress.go
│   │   ├── repository
│   │   │   └── sqlite_article_repo.go
│   │   ├── scraper
│   │   │   ├── client.go
│   │   │   ├── rss.go
│   │   │   ├── scraper_test.go
│   │   │   ├── stealth.go
│   │   │   ├── universal.go
│   │   │   └── wrapper.go
│   │   └── utils
│   │       ├── browser.go
│   │       ├── downloader.go
│   │       ├── logger.go
│   │       └── sanitizer.go
│   └── usecase
│       ├── article_usecase.go
│       └── article_usecase_test.go
├── map_generator.sh
├── run_bot.bat
└── test_build.exe

19 directories, 39 files

## Go Files, Structs, and Functions
### File: ./internal/config/config.go
```go
type Config struct {
type Topic struct {
func InitDB(db *sqlx.DB) {
func Load() error {
func loadFromSQLiteLocked() error {
func Save(cfg Config) error {
func Get() Config {
func saveConfigToSQLiteLocked(cfg Config) error {
```

### File: ./internal/domain/services/ai_interface.go
```go
type AIResult struct {
type AIProvider interface {
```

### File: ./internal/domain/services/interfaces.go
```go
type Publisher interface {
type Linker interface {
type ScraperService interface {
```

### File: ./internal/domain/repository/interfaces.go
```go
type ArticleRepository interface {
type FeedRepository interface {
```

### File: ./internal/domain/models/feed.go
```go
type Feed struct {
```

### File: ./internal/domain/models/article.go
```go
type Article struct {
```

### File: ./internal/infra/database/sqlite.go
```go
func NewSQLiteDB(dbPath string) *sqlx.DB {
```

### File: ./internal/infra/repository/sqlite_article_repo.go
```go
type SQLiteArticleRepository struct {
func NewSQLiteArticleRepository(db *sqlx.DB) *SQLiteArticleRepository {
func (r *SQLiteArticleRepository) Save(ctx context.Context, a *models.Article) error {
func (r *SQLiteArticleRepository) Update(ctx context.Context, a *models.Article) error {
func (r *SQLiteArticleRepository) GetUnprocessed(ctx context.Context, limit int) ([]models.Article, error) {
func (r *SQLiteArticleRepository) GetPending(ctx context.Context, limit int) ([]models.Article, error) {
func (r *SQLiteArticleRepository) GetFailed(ctx context.Context, limit int) ([]models.Article, error) {
func (r *SQLiteArticleRepository) Exists(ctx context.Context, url string) (bool, error) {
func (r *SQLiteArticleRepository) GetByID(ctx context.Context, id int64) (*models.Article, error) {
func (r *SQLiteArticleRepository) GetOneRewrittenByCategory(ctx context.Context, catID int) (*models.Article, error) {
func (r *SQLiteArticleRepository) GetRelated(ctx context.Context, id int64, limit int) ([]models.Article, error) {
func (r *SQLiteArticleRepository) GetStats(ctx context.Context) (published, pending, failed int, err error) {
func (r *SQLiteArticleRepository) Delete(ctx context.Context, id int64) error {
func (r *SQLiteArticleRepository) ClearQueue(ctx context.Context) error {
```

### File: ./internal/infra/api/server.go
```go
type Server struct {
func NewServer(onStart, onStop func(), repo *repository.SQLiteArticleRepository, uc *usecase.ArticleUseCase) *Server {
func (s *Server) SetBotRunning(running bool) {
func (s *Server) handleGetQueueItem(c *fiber.Ctx) error {
func (s *Server) handleUpdateQueueItem(c *fiber.Ctx) error {
func (s *Server) handlePublishItem(c *fiber.Ctx) error {
func (s *Server) handleGetAnalytics(c *fiber.Ctx) error {
func (s *Server) handleGetQueue(c *fiber.Ctx) error {
func (s *Server) handleDeleteItem(c *fiber.Ctx) error {
func (s *Server) handleClearQueue(c *fiber.Ctx) error {
func (s *Server) handleStatus(c *fiber.Ctx) error {
func (s *Server) handleToggle(c *fiber.Ctx) error {
func (s *Server) handleGetSettings(c *fiber.Ctx) error {
func (s *Server) handlePostSettings(c *fiber.Ctx) error {
func (s *Server) Listen(addr string) error {
```

### File: ./internal/infra/utils/browser.go
```go
func GetRandomUserAgent() string {
func OpenBrowser(url string) error {
```

### File: ./internal/infra/utils/logger.go
```go
type LogHub struct {
func (h *LogHub) Register() chan string {
func (h *LogHub) Unregister(ch chan string) {
type wsEvent struct {
func BroadcastEvent(eventType string, data interface{}) {
func BroadcastLog(format string, v ...interface{}) {
```

### File: ./internal/infra/utils/sanitizer.go
```go
func SanitizeHTML(html string) string {
func CleanJSON(raw string) string {
```

### File: ./internal/infra/utils/downloader.go
```go
func isSafeURL(rawURL string) error {
func getSafeClient() *http.Client {
func ToBase64(data []byte) string {
func DownloadFileToBytes(url string) ([]byte, error) {
func DownloadFileToPath(url, path string) error {
```

### File: ./internal/infra/publisher/wordpress.go
```go
type WPClient struct {
func NewWPClient(url, user, pass string, ai services.AIProvider) *WPClient {
func (wp *WPClient) doRequest(req *http.Request) (*http.Response, error) {
func (wp *WPClient) UploadImageFromURL(imageURL, altText string) (int, error) {
func (wp *WPClient) UploadImageToMediaLibrary(imagePath, altText string) (int, string, error) {
func (wp *WPClient) ProcessInlineImages(html string) (string, error) {
func (wp *WPClient) Publish(article *models.Article, catID int) (string, error) {
```

### File: ./internal/infra/publisher/linker.go
```go
type InternalLinker struct {
func NewInternalLinker(repo repository.ArticleRepository) *InternalLinker {
func (l *InternalLinker) InjectLinks(ctx context.Context, articleID int64, htmlContent string) string {
```

### File: ./internal/infra/scraper/wrapper.go
```go
type Wrapper struct {
func NewWrapper() (services.ScraperService, error) {
func (w *Wrapper) FetchRSSLinks(ctx context.Context, url string) ([]string, error) {
func (w *Wrapper) ScrapeArticle(ctx context.Context, url string) (title string, content string, image string, err error) {
func (w *Wrapper) Close() {
```

### File: ./internal/infra/scraper/universal.go
```go
type UniversalScraper struct {
func NewUniversalScraper(b *rod.Browser) *UniversalScraper {
func (s *UniversalScraper) Scrape(ctx context.Context, rawURL string) (string, string, string, error) {
```

### File: ./internal/infra/scraper/client.go
```go
type StealthClient struct {
func GetProxyURL() *url.URL {
func NewStealthClient() *StealthClient {
func (s *StealthClient) Get(url string) (string, error) {
```

### File: ./internal/infra/scraper/rss.go
```go
func FetchRSSLinks(ctx context.Context, browser *rod.Browser, url string) ([]string, error) {
```

### File: ./internal/infra/scraper/scraper_test.go
```go
func TestURLValidation(t *testing.T) {
```

### File: ./internal/infra/scraper/stealth.go
```go
func NewStealthBrowser() (*rod.Browser, error) {
func PreparePage(b *rod.Browser) *rod.Page {
```

### File: ./internal/infra/ai/nvidia.go
```go
type NvidiaProvider struct {
func NewNvidiaProvider(apiKey, sysPrompt string) *NvidiaProvider {
func (n *NvidiaProvider) ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error) {
func (n *NvidiaProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error) {
func (n *NvidiaProvider) runSimpleCompletion(ctx context.Context, url, model, prompt string) (string, error) {
func (n *NvidiaProvider) tryModel(ctx context.Context, url, model, prompt string) (*services.AIResult, error) {
func (n *NvidiaProvider) executeRequest(ctx context.Context, url string, bodyData interface{}) (string, error) {
func (n *NvidiaProvider) Close() {}
```

### File: ./internal/infra/ai/modelslab.go
```go
type ModelsLabProvider struct {
func NewModelsLabProvider(apiKey string) *ModelsLabProvider {
func (m *ModelsLabProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error) {
func (m *ModelsLabProvider) ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error) {
```

### File: ./internal/infra/ai/factory.go
```go
func NewProvider(ctx context.Context, tool, apiKey, sysPrompt string) (services.AIProvider, error) {
```

### File: ./internal/infra/factory.go
```go
type DefaultProviderFactory struct{}
func NewDefaultProviderFactory() *DefaultProviderFactory {
func (f *DefaultProviderFactory) CreateAI(cfg config.Config) services.AIProvider {
func (f *DefaultProviderFactory) CreatePublisher(cfg config.Config) services.Publisher {
func (f *DefaultProviderFactory) CreateScraper() (services.ScraperService, error) {
```

### File: ./internal/usecase/article_usecase.go
```go
type ProviderFactory interface {
type ArticleUseCase struct {
func NewArticleUseCase(repo repository.ArticleRepository, linker services.Linker, factory ProviderFactory) *ArticleUseCase {
func (u *ArticleUseCase) ExecuteScrapeCycle(ctx context.Context, cfg config.Config, startIndex int) int {
func (u *ArticleUseCase) ExecuteAICycle(ctx context.Context, cfg config.Config) {
func (u *ArticleUseCase) processWithAI(ctx context.Context, art *models.Article, aiProv services.AIProvider) {
func (u *ArticleUseCase) PublishSingle(ctx context.Context, cfg config.Config, art *models.Article) {
func (u *ArticleUseCase) ExecutePublishCycle(ctx context.Context, cfg config.Config, startIndex int) int {
```

### File: ./internal/usecase/article_usecase_test.go
```go
func TestExecutePublishCycle_EmptyTopics(t *testing.T) {
func TestContextTimeout(t *testing.T) {
```

### File: ./cmd/bot/main.go
```go
type BotRunner struct {
func NewBotRunner(uc *usecase.ArticleUseCase) *BotRunner {
func (r *BotRunner) SetServer(s *api.Server) {
func (r *BotRunner) Start() {
func (r *BotRunner) Stop() {
func main() {
func (r *BotRunner) runScrapeLoop(ctx context.Context) {
func (r *BotRunner) runAILoop(ctx context.Context) {
func (r *BotRunner) runPublishLoop(ctx context.Context) {
```
