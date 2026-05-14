# Բոտի Ամբողջական Քարտեզ (Complete Codebase Map)

## Պանակների Կառուցվածք (Directory Structure)
```
.
./cmd
./cmd/bot
./internal
./internal/config
./internal/domain
./internal/domain/models
./internal/domain/repository
./internal/domain/services
./internal/infra
./internal/infra/ai
./internal/infra/api
./internal/infra/api/static
./internal/infra/api/static/uploads
./internal/infra/database
./internal/infra/publisher
./internal/infra/repository
./internal/infra/scraper
./internal/infra/utils
./internal/usecase
```

## Գո Ֆայլեր և Ֆունկցիաներ (Go Files and Functions)
### ./cmd/bot/main.go
```go
func NewBotRunner(uc *usecase.ArticleUseCase) *BotRunner {
func (r *BotRunner) SetServer(s *api.Server) {
func (r *BotRunner) Start() {
func (r *BotRunner) Stop() {
func main() {
func (r *BotRunner) runScrapeLoop(ctx context.Context) {
func (r *BotRunner) runAILoop(ctx context.Context) {
func (r *BotRunner) runPublishLoop(ctx context.Context) {
```

### ./internal/config/config.go
```go
func InitDB(db *sqlx.DB) {
func Load() error {
func loadFromSQLiteLocked() error {
func Save(cfg Config) error {
func Get() Config {
func saveConfigToSQLiteLocked(cfg Config) error {
```

### ./internal/domain/models/article.go
```go
```

### ./internal/domain/models/feed.go
```go
```

### ./internal/domain/repository/interfaces.go
```go
```

### ./internal/domain/services/ai_interface.go
```go
```

### ./internal/domain/services/interfaces.go
```go
```

### ./internal/infra/ai/factory.go
```go
func NewProvider(ctx context.Context, tool, apiKey, sysPrompt string) (services.AIProvider, error) {
```

### ./internal/infra/ai/modelslab.go
```go
func NewModelsLabProvider(apiKey string) *ModelsLabProvider {
func (m *ModelsLabProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error) {
func (m *ModelsLabProvider) ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error) {
```

### ./internal/infra/ai/nvidia.go
```go
func NewNvidiaProvider(apiKey, sysPrompt string) *NvidiaProvider {
func (n *NvidiaProvider) ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error) {
func (n *NvidiaProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error) {
func (n *NvidiaProvider) runSimpleCompletion(ctx context.Context, url, model, prompt string) (string, error) {
func (n *NvidiaProvider) tryModel(ctx context.Context, url, model, prompt string) (*services.AIResult, error) {
func (n *NvidiaProvider) executeRequest(ctx context.Context, url string, bodyData interface{}) (string, error) {
func (n *NvidiaProvider) Close() {}
```

### ./internal/infra/api/server.go
```go
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

### ./internal/infra/database/sqlite.go
```go
func NewSQLiteDB(dbPath string) *sqlx.DB {
```

### ./internal/infra/factory.go
```go
func NewDefaultProviderFactory() *DefaultProviderFactory {
func (f *DefaultProviderFactory) CreateAI(cfg config.Config) services.AIProvider {
func (f *DefaultProviderFactory) CreatePublisher(cfg config.Config) services.Publisher {
func (f *DefaultProviderFactory) CreateScraper() (services.ScraperService, error) {
```

### ./internal/infra/publisher/linker.go
```go
func NewInternalLinker(repo repository.ArticleRepository) *InternalLinker {
func (l *InternalLinker) InjectLinks(ctx context.Context, articleID int64, htmlContent string) string {
```

### ./internal/infra/publisher/wordpress.go
```go
func NewWPClient(url, user, pass string, ai services.AIProvider) *WPClient {
func (wp *WPClient) doRequest(req *http.Request) (*http.Response, error) {
func (wp *WPClient) UploadImageFromURL(imageURL, altText string) (int, error) {
func (wp *WPClient) UploadImageToMediaLibrary(imagePath, altText string) (int, string, error) {
func (wp *WPClient) ProcessInlineImages(html string) (string, error) {
func (wp *WPClient) Publish(article *models.Article, catID int) (string, error) {
```

### ./internal/infra/repository/sqlite_article_repo.go
```go
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

### ./internal/infra/scraper/client.go
```go
func GetProxyURL() *url.URL {
func NewStealthClient() *StealthClient {
func (s *StealthClient) Get(url string) (string, error) {
```

### ./internal/infra/scraper/rss.go
```go
func FetchRSSLinks(ctx context.Context, browser *rod.Browser, url string) ([]string, error) {
```

### ./internal/infra/scraper/scraper_test.go
```go
func TestURLValidation(t *testing.T) {
```

### ./internal/infra/scraper/stealth.go
```go
func NewStealthBrowser() (*rod.Browser, error) {
func PreparePage(b *rod.Browser) *rod.Page {
```

### ./internal/infra/scraper/universal.go
```go
func NewUniversalScraper(b *rod.Browser) *UniversalScraper {
func (s *UniversalScraper) Scrape(ctx context.Context, rawURL string) (string, string, string, error) {
```

### ./internal/infra/scraper/wrapper.go
```go
func NewWrapper() (services.ScraperService, error) {
func (w *Wrapper) FetchRSSLinks(ctx context.Context, url string) ([]string, error) {
func (w *Wrapper) ScrapeArticle(ctx context.Context, url string) (title string, content string, image string, err error) {
func (w *Wrapper) Close() {
```

### ./internal/infra/utils/browser.go
```go
func GetRandomUserAgent() string {
func OpenBrowser(url string) error {
```

### ./internal/infra/utils/downloader.go
```go
func isSafeURL(rawURL string) error {
func getSafeClient() *http.Client {
func ToBase64(data []byte) string {
func DownloadFileToBytes(url string) ([]byte, error) {
func DownloadFileToPath(url, path string) error {
```

### ./internal/infra/utils/logger.go
```go
func (h *LogHub) Register() chan string {
func (h *LogHub) Unregister(ch chan string) {
func BroadcastEvent(eventType string, data interface{}) {
func BroadcastLog(format string, v ...interface{}) {
```

### ./internal/infra/utils/sanitizer.go
```go
func SanitizeHTML(html string) string {
func CleanJSON(raw string) string {
```

### ./internal/usecase/article_usecase.go
```go
func NewArticleUseCase(repo repository.ArticleRepository, linker services.Linker, factory ProviderFactory) *ArticleUseCase {
func (u *ArticleUseCase) ExecuteScrapeCycle(ctx context.Context, cfg config.Config, startIndex int) int {
func (u *ArticleUseCase) ExecuteAICycle(ctx context.Context, cfg config.Config) {
func (u *ArticleUseCase) processWithAI(ctx context.Context, art *models.Article, aiProv services.AIProvider) {
func (u *ArticleUseCase) PublishSingle(ctx context.Context, cfg config.Config, art *models.Article) {
func (u *ArticleUseCase) ExecutePublishCycle(ctx context.Context, cfg config.Config, startIndex int) int {
```

### ./internal/usecase/article_usecase_test.go
```go
func TestExecutePublishCycle_EmptyTopics(t *testing.T) {
func TestContextTimeout(t *testing.T) {
```

## Տվյալների Բազայի Կառուցվածք (Database Information)
Բազան SQLite է։ Ստորև ներկայացված են SQL հրամանները՝ աղյուսակներ ստեղծելու համար.
```sql
./internal/infra/database/sqlite.go-
./internal/infra/database/sqlite.go-	schema := `
./internal/infra/database/sqlite.go:	CREATE TABLE IF NOT EXISTS app_settings (
./internal/infra/database/sqlite.go-		id INTEGER PRIMARY KEY CHECK (id = 1),
./internal/infra/database/sqlite.go-		wp_url TEXT,
./internal/infra/database/sqlite.go-		wp_username TEXT,
./internal/infra/database/sqlite.go-		wp_app_password TEXT,
./internal/infra/database/sqlite.go-		ai_tool TEXT,
./internal/infra/database/sqlite.go-		nvidia_api_key TEXT,
./internal/infra/database/sqlite.go-		modelslab_api_key TEXT,
./internal/infra/database/sqlite.go-		min_article_len INTEGER,
./internal/infra/database/sqlite.go-		run_interval_hours INTEGER,
./internal/infra/database/sqlite.go-		run_interval_minutes INTEGER,
./internal/infra/database/sqlite.go-		publish_interval_minutes INTEGER,
./internal/infra/database/sqlite.go-		auto_publish BOOLEAN,
./internal/infra/database/sqlite.go-		system_prompt_nvidia TEXT,
./internal/infra/database/sqlite.go-		system_prompt_modelslab TEXT
./internal/infra/database/sqlite.go-	);
./internal/infra/database/sqlite.go-
./internal/infra/database/sqlite.go:	CREATE TABLE IF NOT EXISTS rss_topics (
./internal/infra/database/sqlite.go-		id INTEGER PRIMARY KEY AUTOINCREMENT,
./internal/infra/database/sqlite.go-		name TEXT,
./internal/infra/database/sqlite.go-		wp_category_id INTEGER,
./internal/infra/database/sqlite.go-		rss_url TEXT UNIQUE
./internal/infra/database/sqlite.go-	);
./internal/infra/database/sqlite.go-
./internal/infra/database/sqlite.go:	CREATE TABLE IF NOT EXISTS articles (
./internal/infra/database/sqlite.go-		id INTEGER PRIMARY KEY AUTOINCREMENT,
./internal/infra/database/sqlite.go-		source_url TEXT UNIQUE,
./internal/infra/database/sqlite.go-		title TEXT,
./internal/infra/database/sqlite.go-		content TEXT,
./internal/infra/database/sqlite.go-		rewritten_content TEXT,
./internal/infra/database/sqlite.go-		image_url TEXT,
./internal/infra/database/sqlite.go-		meta_description TEXT,
./internal/infra/database/sqlite.go-		focus_keywords TEXT,
./internal/infra/database/sqlite.go-		slug TEXT,
./internal/infra/database/sqlite.go-		category TEXT,
./internal/infra/database/sqlite.go-		category_id INTEGER,
./internal/infra/database/sqlite.go-		tags TEXT,
./internal/infra/database/sqlite.go-		image_alt TEXT,
./internal/infra/database/sqlite.go-		status TEXT,
./internal/infra/database/sqlite.go-		retry_count INTEGER DEFAULT 0,
./internal/infra/database/sqlite.go-		next_retry_at DATETIME,
./internal/infra/database/sqlite.go-		publish_date DATETIME,
./internal/infra/database/sqlite.go-		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
./internal/infra/database/sqlite.go-	);
./internal/infra/database/sqlite.go-
```
