# Full Codebase Map
## Files
./bot.exe~
./bot_architecture_map.md
./cmd/bot/main.go
./complete_map.md
./go.mod
./go.sum
./internal/config/config.go
./internal/domain/models/article.go
./internal/domain/models/feed.go
./internal/domain/repository/interfaces.go
./internal/domain/services/ai_interface.go
./internal/domain/services/interfaces.go
./internal/infra/ai/factory.go
./internal/infra/ai/modelslab.go
./internal/infra/ai/nvidia.go
./internal/infra/api/server.go
./internal/infra/api/static/index.html
./internal/infra/api/static/script.js
./internal/infra/api/static/style.css
./internal/infra/api/static/temp_image.jpg
./internal/infra/api/static/uploads/.gitkeep
./internal/infra/api/static/uploads/placeholder.txt
./internal/infra/database/sqlite.go
./internal/infra/factory.go
./internal/infra/publisher/linker.go
./internal/infra/publisher/wordpress.go
./internal/infra/repository/sqlite_article_repo.go
./internal/infra/scraper/client.go
./internal/infra/scraper/rss.go
./internal/infra/scraper/scraper_test.go
./internal/infra/scraper/stealth.go
./internal/infra/scraper/universal.go
./internal/infra/scraper/wrapper.go
./internal/infra/utils/browser.go
./internal/infra/utils/downloader.go
./internal/infra/utils/logger.go
./internal/infra/utils/sanitizer.go
./internal/usecase/article_usecase.go
./internal/usecase/article_usecase_test.go
./map_generator.sh
./run_bot.bat

## DB Schema
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

## Functions & Interfaces
./internal/config/config.go:type Config struct {
./internal/config/config.go:type Topic struct {
./internal/config/config.go:func InitDB(db *sqlx.DB) {
./internal/config/config.go:func Load() error {
./internal/config/config.go:func loadFromSQLiteLocked() error {
./internal/config/config.go:	type flatSettings struct {
./internal/config/config.go:func Save(cfg Config) error {
./internal/config/config.go:func Get() Config {
./internal/config/config.go:func saveConfigToSQLiteLocked(cfg Config) error {
./internal/domain/services/ai_interface.go:type AIResult struct {
./internal/domain/services/ai_interface.go:type AIProvider interface {
./internal/domain/services/interfaces.go:type Publisher interface {
./internal/domain/services/interfaces.go:type Linker interface {
./internal/domain/services/interfaces.go:type ScraperService interface {
./internal/domain/repository/interfaces.go:type ArticleRepository interface {
./internal/domain/repository/interfaces.go:type FeedRepository interface {
./internal/domain/models/feed.go:type Feed struct {
./internal/domain/models/article.go:type Article struct {
./internal/infra/database/sqlite.go:func NewSQLiteDB(dbPath string) *sqlx.DB {
./internal/infra/repository/sqlite_article_repo.go:type SQLiteArticleRepository struct {
./internal/infra/repository/sqlite_article_repo.go:func NewSQLiteArticleRepository(db *sqlx.DB) *SQLiteArticleRepository {
./internal/infra/repository/sqlite_article_repo.go:func (r *SQLiteArticleRepository) Save(ctx context.Context, a *models.Article) error {
./internal/infra/repository/sqlite_article_repo.go:func (r *SQLiteArticleRepository) Update(ctx context.Context, a *models.Article) error {
./internal/infra/repository/sqlite_article_repo.go:func (r *SQLiteArticleRepository) GetUnprocessed(ctx context.Context, limit int) ([]models.Article, error) {
./internal/infra/repository/sqlite_article_repo.go:func (r *SQLiteArticleRepository) GetPending(ctx context.Context, limit int) ([]models.Article, error) {
./internal/infra/repository/sqlite_article_repo.go:func (r *SQLiteArticleRepository) GetFailed(ctx context.Context, limit int) ([]models.Article, error) {
./internal/infra/repository/sqlite_article_repo.go:func (r *SQLiteArticleRepository) Exists(ctx context.Context, url string) (bool, error) {
./internal/infra/repository/sqlite_article_repo.go:func (r *SQLiteArticleRepository) GetByID(ctx context.Context, id int64) (*models.Article, error) {
./internal/infra/repository/sqlite_article_repo.go:func (r *SQLiteArticleRepository) GetOneRewrittenByCategory(ctx context.Context, catID int) (*models.Article, error) {
./internal/infra/repository/sqlite_article_repo.go:func (r *SQLiteArticleRepository) GetRelated(ctx context.Context, id int64, limit int) ([]models.Article, error) {
./internal/infra/repository/sqlite_article_repo.go:func (r *SQLiteArticleRepository) GetStats(ctx context.Context) (published, pending, failed int, err error) {
./internal/infra/repository/sqlite_article_repo.go:func (r *SQLiteArticleRepository) Delete(ctx context.Context, id int64) error {
./internal/infra/repository/sqlite_article_repo.go:func (r *SQLiteArticleRepository) ClearQueue(ctx context.Context) error {
./internal/infra/api/server.go:type Server struct {
./internal/infra/api/server.go:func NewServer(onStart, onStop func(), repo *repository.SQLiteArticleRepository, uc *usecase.ArticleUseCase) *Server {
./internal/infra/api/server.go:func (s *Server) SetBotRunning(running bool) {
./internal/infra/api/server.go:func (s *Server) handleGetQueueItem(c *fiber.Ctx) error {
./internal/infra/api/server.go:func (s *Server) handleUpdateQueueItem(c *fiber.Ctx) error {
./internal/infra/api/server.go:func (s *Server) handlePublishItem(c *fiber.Ctx) error {
./internal/infra/api/server.go:func (s *Server) handleGetAnalytics(c *fiber.Ctx) error {
./internal/infra/api/server.go:func (s *Server) handleGetQueue(c *fiber.Ctx) error {
./internal/infra/api/server.go:func (s *Server) handleDeleteItem(c *fiber.Ctx) error {
./internal/infra/api/server.go:func (s *Server) handleClearQueue(c *fiber.Ctx) error {
./internal/infra/api/server.go:func (s *Server) handleStatus(c *fiber.Ctx) error {
./internal/infra/api/server.go:func (s *Server) handleToggle(c *fiber.Ctx) error {
./internal/infra/api/server.go:func (s *Server) handleGetSettings(c *fiber.Ctx) error {
./internal/infra/api/server.go:func (s *Server) handlePostSettings(c *fiber.Ctx) error {
./internal/infra/api/server.go:func (s *Server) Listen(addr string) error {
./internal/infra/utils/browser.go:func GetRandomUserAgent() string {
./internal/infra/utils/browser.go:func OpenBrowser(url string) error {
./internal/infra/utils/logger.go:type LogHub struct {
./internal/infra/utils/logger.go:func (h *LogHub) Register() chan string {
./internal/infra/utils/logger.go:func (h *LogHub) Unregister(ch chan string) {
./internal/infra/utils/logger.go:type wsEvent struct {
./internal/infra/utils/logger.go:func BroadcastEvent(eventType string, data interface{}) {
./internal/infra/utils/logger.go:func BroadcastLog(format string, v ...interface{}) {
./internal/infra/utils/sanitizer.go:func SanitizeHTML(html string) string {
./internal/infra/utils/sanitizer.go:func CleanJSON(raw string) string {
./internal/infra/utils/downloader.go:func isSafeURL(rawURL string) error {
./internal/infra/utils/downloader.go:func getSafeClient() *http.Client {
./internal/infra/utils/downloader.go:func ToBase64(data []byte) string {
./internal/infra/utils/downloader.go:func DownloadFileToBytes(url string) ([]byte, error) {
./internal/infra/utils/downloader.go:func DownloadFileToPath(url, path string) error {
./internal/infra/publisher/wordpress.go:type WPClient struct {
./internal/infra/publisher/wordpress.go:func NewWPClient(url, user, pass string, ai services.AIProvider) *WPClient {
./internal/infra/publisher/wordpress.go:func (wp *WPClient) doRequest(req *http.Request) (*http.Response, error) {
./internal/infra/publisher/wordpress.go:func (wp *WPClient) UploadImageFromURL(imageURL, altText string) (int, error) {
./internal/infra/publisher/wordpress.go:func (wp *WPClient) UploadImageToMediaLibrary(imagePath, altText string) (int, string, error) {
./internal/infra/publisher/wordpress.go:func (wp *WPClient) ProcessInlineImages(html string) (string, error) {
./internal/infra/publisher/wordpress.go:	type uploadResult struct {
./internal/infra/publisher/wordpress.go:func (wp *WPClient) Publish(article *models.Article, catID int) (string, error) {
./internal/infra/publisher/linker.go:type InternalLinker struct {
./internal/infra/publisher/linker.go:func NewInternalLinker(repo repository.ArticleRepository) *InternalLinker {
./internal/infra/publisher/linker.go:func (l *InternalLinker) InjectLinks(ctx context.Context, articleID int64, htmlContent string) string {
./internal/infra/scraper/wrapper.go:type Wrapper struct {
./internal/infra/scraper/wrapper.go:func NewWrapper() (services.ScraperService, error) {
./internal/infra/scraper/wrapper.go:func (w *Wrapper) FetchRSSLinks(ctx context.Context, url string) ([]string, error) {
./internal/infra/scraper/wrapper.go:func (w *Wrapper) ScrapeArticle(ctx context.Context, url string) (title string, content string, image string, err error) {
./internal/infra/scraper/wrapper.go:func (w *Wrapper) Close() {
./internal/infra/scraper/universal.go:type UniversalScraper struct {
./internal/infra/scraper/universal.go:func NewUniversalScraper(b *rod.Browser) *UniversalScraper {
./internal/infra/scraper/universal.go:func (s *UniversalScraper) Scrape(ctx context.Context, rawURL string) (string, string, string, error) {
./internal/infra/scraper/client.go:type StealthClient struct {
./internal/infra/scraper/client.go:func GetProxyURL() *url.URL {
./internal/infra/scraper/client.go:func NewStealthClient() *StealthClient {
./internal/infra/scraper/client.go:func (s *StealthClient) Get(url string) (string, error) {
./internal/infra/scraper/rss.go:func FetchRSSLinks(ctx context.Context, browser *rod.Browser, url string) ([]string, error) {
./internal/infra/scraper/scraper_test.go:func TestURLValidation(t *testing.T) {
./internal/infra/scraper/stealth.go:func NewStealthBrowser() (*rod.Browser, error) {
./internal/infra/scraper/stealth.go:func PreparePage(b *rod.Browser) *rod.Page {
./internal/infra/ai/nvidia.go:type NvidiaProvider struct {
./internal/infra/ai/nvidia.go:func NewNvidiaProvider(apiKey, sysPrompt string) *NvidiaProvider {
./internal/infra/ai/nvidia.go:func (n *NvidiaProvider) ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error) {
./internal/infra/ai/nvidia.go:func (n *NvidiaProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error) {
./internal/infra/ai/nvidia.go:func (n *NvidiaProvider) runSimpleCompletion(ctx context.Context, url, model, prompt string) (string, error) {
./internal/infra/ai/nvidia.go:func (n *NvidiaProvider) tryModel(ctx context.Context, url, model, prompt string) (*services.AIResult, error) {
./internal/infra/ai/nvidia.go:func (n *NvidiaProvider) executeRequest(ctx context.Context, url string, bodyData interface{}) (string, error) {
./internal/infra/ai/nvidia.go:func (n *NvidiaProvider) Close() {}
./internal/infra/ai/modelslab.go:type ModelsLabProvider struct {
./internal/infra/ai/modelslab.go:func NewModelsLabProvider(apiKey string) *ModelsLabProvider {
./internal/infra/ai/modelslab.go:func (m *ModelsLabProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error) {
./internal/infra/ai/modelslab.go:func (m *ModelsLabProvider) ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error) {
./internal/infra/ai/factory.go:func NewProvider(ctx context.Context, tool, apiKey, sysPrompt string) (services.AIProvider, error) {
./internal/infra/factory.go:type DefaultProviderFactory struct{}
./internal/infra/factory.go:func NewDefaultProviderFactory() *DefaultProviderFactory {
./internal/infra/factory.go:func (f *DefaultProviderFactory) CreateAI(cfg config.Config) services.AIProvider {
./internal/infra/factory.go:func (f *DefaultProviderFactory) CreatePublisher(cfg config.Config) services.Publisher {
./internal/infra/factory.go:func (f *DefaultProviderFactory) CreateScraper() (services.ScraperService, error) {
./internal/usecase/article_usecase.go:type ProviderFactory interface {
./internal/usecase/article_usecase.go:type ArticleUseCase struct {
./internal/usecase/article_usecase.go:func NewArticleUseCase(repo repository.ArticleRepository, linker services.Linker, factory ProviderFactory) *ArticleUseCase {
./internal/usecase/article_usecase.go:func (u *ArticleUseCase) ExecuteScrapeCycle(ctx context.Context, cfg config.Config, startIndex int) int {
./internal/usecase/article_usecase.go:func (u *ArticleUseCase) ExecuteAICycle(ctx context.Context, cfg config.Config) {
./internal/usecase/article_usecase.go:func (u *ArticleUseCase) processWithAI(ctx context.Context, art *models.Article, aiProv services.AIProvider) {
./internal/usecase/article_usecase.go:func (u *ArticleUseCase) PublishSingle(ctx context.Context, cfg config.Config, art *models.Article) {
./internal/usecase/article_usecase.go:func (u *ArticleUseCase) ExecutePublishCycle(ctx context.Context, cfg config.Config, startIndex int) int {
./internal/usecase/article_usecase_test.go:func TestExecutePublishCycle_EmptyTopics(t *testing.T) {
./internal/usecase/article_usecase_test.go:func TestContextTimeout(t *testing.T) {
./cmd/bot/main.go:type BotRunner struct {
./cmd/bot/main.go:func NewBotRunner(uc *usecase.ArticleUseCase) *BotRunner {
./cmd/bot/main.go:func (r *BotRunner) SetServer(s *api.Server) {
./cmd/bot/main.go:func (r *BotRunner) Start() {
./cmd/bot/main.go:func (r *BotRunner) Stop() {
./cmd/bot/main.go:func main() {
./cmd/bot/main.go:func (r *BotRunner) runScrapeLoop(ctx context.Context) {
./cmd/bot/main.go:func (r *BotRunner) runAILoop(ctx context.Context) {
./cmd/bot/main.go:func (r *BotRunner) runPublishLoop(ctx context.Context) {
