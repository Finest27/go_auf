# Կոդերի Բազայի Քարտեզ - ThePress Bot Ultimate

Այս փաստաթուղթը պարունակում է բոտի կոդերի բազայի ամբողջական քարտեզը՝ հայերեն լեզվով։

## Տեղեկատուների (Դիրեկտորիաների) Կառուցվածք
- `cmd/bot/` - Ծրագրի մուտքի կետը (main.go)
- `internal/config/` - Կարգավորումների ղեկավարում (config.go)
- `internal/domain/` - Հիմնական մոդելները և ինտերֆեյսները
  - `models/` - Տվյալների մոդելներ (օր.՝ Article)
  - `repository/` - Տվյալների բազայի ինտերֆեյսներ
  - `services/` - Ծառայությունների ինտերֆեյսներ
- `internal/infra/` - Ենթակառուցվածքների իրականացումներ (տվյալների բազա, արտաքին API-ներ)
  - `ai/` - AI պրովայդերների իրականացում (Nvidia, ModelsLab)
  - `api/` - Վեբ սերվերի իրականացում և UI
  - `database/` - SQLite բազայի միացում և սխեմա
  - `publisher/` - WordPress-ում հրապարակելու տրամաբանություն
  - `repository/` - Տվյալների բազայի պահոցների (repository) իրականացում
  - `scraper/` - Հոդվածների քերման (scraping) տրամաբանություն
  - `utils/` - Օժանդակ ֆունկցիաներ (լոգեր, բրաուզեր և այլն)
- `internal/usecase/` - Բիզնես տրամաբանություն (օր.՝ ArticleUseCase)

## Հիմնական Ֆայլեր և Ֆունկցիաներ
- `cmd/bot/main.go`: `main()` - Գլխավոր ֆունկցիան, `runScrapeLoop()`, `runAILoop()`, `runPublishLoop()` - Ֆոնային ցիկլեր
- `internal/usecase/article_usecase.go`: `ExecuteScrapeCycle()`, `ExecuteAICycle()`, `ExecutePublishCycle()`
- `internal/infra/database/sqlite.go`: `NewSQLiteDB()` - Տվյալների բազայի ստեղծում և կառուցվածք
- `internal/infra/repository/sqlite_article_repo.go`: Հոդվածների պահպանում և ստացում տվյալների բազայից

## Տվյալների Բազա (SQLite)
Հիմնական ֆայլ: `bot_ultimate.db`
Աղյուսակներ:
1. `app_settings` - Բոտի կարգավորումներ (API բանալիներ, ինտերվալներ և այլն)
2. `rss_topics` - RSS թեմաների ցանկ և կապակցված WordPress կատեգորիաներ
3. `articles` - Պահպանված հոդվածներ (source_url, title, content, rewritten_content, status, publish_date և այլն)
