# ThePress Bot - Ամբողջական Քարտեզ (Complete Map)
Այս ֆայլը պարունակում է բոտի բոլոր ֆայլերի, ֆունկցիաների և տվյալների բազայի (DB) կառուցվածքի մասին տեղեկատվություն։

## Տվյալների բազա (Database Schema)
Բազան SQLite է, որը գտնվում է `bot_ultimate.db` ֆայլում։ Ունի հետևյալ աղյուսակները՝
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

```

## Կոդի կառուցվածք (Code Structure & Functions)
### thepress_bot_go/cmd/bot/main.go
```go
func NewBotRunner(uc *usecase.ArticleUseCase) *BotRunner {
func (r *BotRunner) SetServer(s *api.Server) {
func (r *BotRunner) Start() {
func (r *BotRunner) Stop() {
func main() {
func (r *BotRunner) runScrapeLoop(ctx context.Context) {
func (r *BotRunner) runAILoop(ctx context.Context) {
func (r *BotRunner) runPublishLoop(ctx context.Context) {
type BotRunner struct {
```

### thepress_bot_go/internal/config/config.go
```go
func InitDB(db *sqlx.DB) {
func Load() error {
func loadFromSQLiteLocked() error {
func Save(cfg Config) error {
func Get() Config {
func saveConfigToSQLiteLocked(cfg Config) error {
type Config struct {
type Topic struct {
```

### thepress_bot_go/internal/domain/models/article.go
```go
type Article struct {
```

### thepress_bot_go/internal/domain/models/feed.go
```go
type Feed struct {
```

### thepress_bot_go/internal/domain/repository/interfaces.go
```go
type ArticleRepository interface {
type FeedRepository interface {
```

### thepress_bot_go/internal/domain/services/ai_interface.go
```go
type AIResult struct {
type AIProvider interface {
```

### thepress_bot_go/internal/domain/services/interfaces.go
```go
type Publisher interface {
type Linker interface {
type ScraperService interface {
```

### thepress_bot_go/internal/infra/ai/factory.go
```go
func NewProvider(ctx context.Context, tool, apiKey, sysPrompt string) (services.AIProvider, error) {
```

### thepress_bot_go/internal/infra/ai/modelslab.go
```go
func NewModelsLabProvider(apiKey string) *ModelsLabProvider {
func (m *ModelsLabProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error) {
func (m *ModelsLabProvider) ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error) {
type ModelsLabProvider struct {
```

### thepress_bot_go/internal/infra/ai/nvidia.go
```go
func NewNvidiaProvider(apiKey, sysPrompt string) *NvidiaProvider {
func (n *NvidiaProvider) ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error) {
func (n *NvidiaProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error) {
func (n *NvidiaProvider) runSimpleCompletion(ctx context.Context, url, model, prompt string) (string, error) {
func (n *NvidiaProvider) tryModel(ctx context.Context, url, model, prompt string) (*services.AIResult, error) {
func (n *NvidiaProvider) executeRequest(ctx context.Context, url string, bodyData interface{}) (string, error) {
func (n *NvidiaProvider) Close() {}
type NvidiaProvider struct {
```

### thepress_bot_go/internal/infra/api/server.go
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
type Server struct {
```

### thepress_bot_go/internal/infra/database/sqlite.go
```go
func NewSQLiteDB(dbPath string) *sqlx.DB {
```

### thepress_bot_go/internal/infra/factory.go
```go
func NewDefaultProviderFactory() *DefaultProviderFactory {
func (f *DefaultProviderFactory) CreateAI(cfg config.Config) services.AIProvider {
func (f *DefaultProviderFactory) CreatePublisher(cfg config.Config) services.Publisher {
func (f *DefaultProviderFactory) CreateScraper() (services.ScraperService, error) {
type DefaultProviderFactory struct{}
```

### thepress_bot_go/internal/infra/publisher/linker.go
```go
func NewInternalLinker(repo repository.ArticleRepository) *InternalLinker {
func (l *InternalLinker) InjectLinks(ctx context.Context, articleID int64, htmlContent string) string {
type InternalLinker struct {
```

### thepress_bot_go/internal/infra/publisher/wordpress.go
```go
func NewWPClient(url, user, pass string, ai services.AIProvider) *WPClient {
func (wp *WPClient) doRequest(req *http.Request) (*http.Response, error) {
func (wp *WPClient) UploadImageFromURL(imageURL, altText string) (int, error) {
func (wp *WPClient) UploadImageToMediaLibrary(imagePath, altText string) (int, string, error) {
func (wp *WPClient) ProcessInlineImages(html string) (string, error) {
func (wp *WPClient) Publish(article *models.Article, catID int) (string, error) {
type WPClient struct {
```

### thepress_bot_go/internal/infra/repository/sqlite_article_repo.go
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
type SQLiteArticleRepository struct {
```

### thepress_bot_go/internal/infra/scraper/client.go
```go
func GetProxyURL() *url.URL {
func NewStealthClient() *StealthClient {
func (s *StealthClient) Get(url string) (string, error) {
type StealthClient struct {
```

### thepress_bot_go/internal/infra/scraper/rss.go
```go
func FetchRSSLinks(ctx context.Context, browser *rod.Browser, url string) ([]string, error) {
```

### thepress_bot_go/internal/infra/scraper/scraper_test.go
```go
func TestURLValidation(t *testing.T) {
```

### thepress_bot_go/internal/infra/scraper/stealth.go
```go
func NewStealthBrowser() (*rod.Browser, error) {
func PreparePage(b *rod.Browser) *rod.Page {
```

### thepress_bot_go/internal/infra/scraper/universal.go
```go
func NewUniversalScraper(b *rod.Browser) *UniversalScraper {
func (s *UniversalScraper) Scrape(ctx context.Context, rawURL string) (string, string, string, error) {
type UniversalScraper struct {
```

### thepress_bot_go/internal/infra/scraper/wrapper.go
```go
func NewWrapper() (services.ScraperService, error) {
func (w *Wrapper) FetchRSSLinks(ctx context.Context, url string) ([]string, error) {
func (w *Wrapper) ScrapeArticle(ctx context.Context, url string) (title string, content string, image string, err error) {
func (w *Wrapper) Close() {
type Wrapper struct {
```

### thepress_bot_go/internal/infra/utils/browser.go
```go
func GetRandomUserAgent() string {
func OpenBrowser(url string) error {
```

### thepress_bot_go/internal/infra/utils/downloader.go
```go
func isSafeURL(rawURL string) error {
func getSafeClient() *http.Client {
func ToBase64(data []byte) string {
func DownloadFileToBytes(url string) ([]byte, error) {
func DownloadFileToPath(url, path string) error {
```

### thepress_bot_go/internal/infra/utils/logger.go
```go
func (h *LogHub) Register() chan string {
func (h *LogHub) Unregister(ch chan string) {
func BroadcastEvent(eventType string, data interface{}) {
func BroadcastLog(format string, v ...interface{}) {
type LogHub struct {
type wsEvent struct {
```

### thepress_bot_go/internal/infra/utils/sanitizer.go
```go
func SanitizeHTML(html string) string {
func CleanJSON(raw string) string {
```

### thepress_bot_go/internal/usecase/article_usecase.go
```go
func NewArticleUseCase(repo repository.ArticleRepository, linker services.Linker, factory ProviderFactory) *ArticleUseCase {
func (u *ArticleUseCase) ExecuteScrapeCycle(ctx context.Context, cfg config.Config, startIndex int) int {
func (u *ArticleUseCase) ExecuteAICycle(ctx context.Context, cfg config.Config) {
func (u *ArticleUseCase) processWithAI(ctx context.Context, art *models.Article, aiProv services.AIProvider) {
func (u *ArticleUseCase) PublishSingle(ctx context.Context, cfg config.Config, art *models.Article) {
func (u *ArticleUseCase) ExecutePublishCycle(ctx context.Context, cfg config.Config, startIndex int) int {
type ArticleUseCase struct {
type ProviderFactory interface {
```

### thepress_bot_go/internal/usecase/article_usecase_test.go
```go
func TestExecutePublishCycle_EmptyTopics(t *testing.T) {
func TestContextTimeout(t *testing.T) {
```
