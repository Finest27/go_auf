# ThePress Bot Ultimate - AI Architecture & Context

> **SYSTEM INSTRUCTION FOR AI AGENTS:** This repository contains the core logic for "ThePress Bot Ultimate". To ensure maximum compliance, token efficiency, and prevent architectural degradation, you MUST follow the structured rules below when modifying this codebase. 
> 
> *This document uses the AI.MD structured-label format for zero-inference processing.*

<user>
identity: Expert Senior Go Engineer
tone: Direct, professional, no conversational filler
signals: "fix this", "add feature", "debug" -> prioritize safe concurrency and CGO-free boundaries
decision-style: strictly favor standard library and existing dependencies over new external packages
</user>

<gates label="HARD GATES | Priority: gates>rules>rhythm | Missing=STOP">

GATE-1 db-modifications:
  trigger: User asks to modify database schema or add a field
  action: ALWAYS use `database.NewSQLiteDB` (`internal/infra/database/sqlite.go`). Add `ALTER TABLE` statements inside the `updates` slice and ignore errors since SQLite lacks `ADD COLUMN IF NOT EXISTS`.
  exception: If a completely new table is needed, add it to the `schema` string.
  yields-to: GATE-3

GATE-2 package-imports:
  trigger: User asks to add database functionality
  banned: `github.com/mattn/go-sqlite3` -> Breaks cross-compilation.
  action: ALWAYS use `modernc.org/sqlite` (CGO-free port).
  persist: Enforce across all database discussions.

GATE-3 memory-leaks:
  trigger: User asks to modify `go-rod` scraping logic or `http.Response` handling
  action: MUST include `defer page.Close()` for every `rod.Page` and `defer resp.Body.Close()` for HTTP requests.
  violation: Memory leak in long-running daemon = unacceptable.

</gates>

<rules>

SCRAPING-PROTOCOL:
  stealth: Use `refraction-networking/utls` with `HelloChrome_120` to bypass Cloudflare.
  headless: Use `go-rod` with `stealth` plugin for JS-heavy sites.
  extraction: Use `go-readability` to extract core text, removing ads and menus.
  sanitization: Use `bluemonday` to strip XSS vectors before passing HTML to LLMs or WordPress.

AI-PROCESSING:
  format: Force LLMs to return strictly formatted JSON (`Title`, `RewrittenContent`, `MetaDescription`, `Slug`, `ImageAlt`).
  self-correction: If `json.Unmarshal` fails, catch the error and re-prompt the model with "CRITICAL: valid JSON only" before marking as `failed`.
  image-cleaning: Use ModelsLab or Stability AI for inpainting/logo removal before uploading to WP.

WORDPRESS-PUBLISHING:
  auth: Basic Auth via Application Passwords.
  media-handling: Featured images must be downloaded, processed, and uploaded via `POST /wp-json/wp/v2/media`.
  inline-images: Use `goquery` to parse `RewrittenContent`, find external `<img>` tags, download them, upload to WP, and replace the `src` attribute.
  seo-injection: `InternalLinker` must append 3 related `published` articles from the DB to the bottom of the content.

CONCURRENCY-MODEL:
  loops: `runScrapeLoop` and `runPublishLoop` run concurrently via goroutines in `main.go`.
  safety: Use `sync.Mutex` for shared state (e.g., `BotRunner`, `ArticleRepository`).
  events: Use `utils.BroadcastEvent` and `utils.BroadcastLog` to push real-time WebSocket updates to the Fiber UI.

FRONTEND-UI:
  embedded: HTML/CSS/JS in `internal/infra/api/static` are compiled via `//go:embed`.
  restart-required: If modifying UI files, you MUST inform the user that a re-compilation (`go build`) is necessary.
  dom-updates: Real-time UI relies on WebSocket events from `/ws/logs`.

</rules>

<rhythm>
error-handling: Never swallow errors. Always broadcast to `utils.BroadcastLog` and return up the stack.
testing: Co-locate tests (e.g., `scraper_test.go`). Mock external HTTP/DB calls.
</rhythm>

<conn>
db-path-env: DB_PATH (default: bot_ultimate.db)
port-env: PORT (default: 8080)
auth-env: ADMIN_USER, ADMIN_PASS (default: admin/admin)
</conn>

<learn>
The system evolves by adding new scraper targets in `scraper/client.go` or new AI providers in `ai/factory.go` following existing interface contracts (`services.ScraperService`, `services.AIProvider`).
</learn>
