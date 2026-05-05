# ThePress Bot Ultimate Ճարտարապետության և Կոդի Քարտեզ

Այս փաստաթուղթը տրամադրում է բոտի կոդի բազայի ամբողջական քարտեզը՝ հեշտացնելով նավարկումը և փոփոխությունների կիրառումը: Այն նկարագրում է թղթապանակների կառուցվածքը, հիմնական ֆայլերը, դրանց ֆունկցիաները և տվյալների բազայի (database) սխեման:

## Թղթապանակների Կառուցվածք (Directory Structure)

*   **`cmd/bot/`**: Պարունակում է ծրագրի մուտքային կետը (entry point):
*   **`internal/config/`**: Կառավարում է հավելվածի կարգավորումները:
*   **`internal/domain/`**: Սահմանում է հիմնական էությունները (մոդելներ) և պայմանագրերը (ինտերֆեյսներ):
*   **`internal/infra/`**: Պարունակում է արտաքին համակարգերի հետ աշխատող կոդը (տվյալների բազա, արտաքին API-ներ, սքրեյփերներ):
*   **`internal/usecase/`**: Իրականացնում է հիմնական բիզնես տրամաբանությունը:

## Հիմնական Ֆայլեր և Ֆունկցիաներ

### 1. `cmd/bot/main.go`
Ծրագրի մեկնարկի կետն է: Այն սկզբնավորում է տվյալների բազան, բեռնում է կարգավորումները, ստեղծում է կախվածությունները և սկսում է ֆոնային պրոցեսները:
*   **`main()`**: Գլխավոր ֆունկցիան, որը միացնում է բոլոր բաղադրիչները:
*   **`BotRunner.Start()`**: Սկսում է ֆոնային scraping-ի, AI-ի և publishing-ի ցիկլերը:
*   **`runScrapeLoop(ctx)`**, **`runAILoop(ctx)`**, **`runPublishLoop(ctx)`**: Ֆոնային անվերջ ցիկլեր, որոնք պարբերաբար կանչում են համապատասխան գործողությունները:

### 2. `internal/config/config.go`
Աշխատում է հավելվածի կարգավորումների հետ, որոնք հիմնականում պահվում են SQLite տվյալների բազայում:
*   **`Config` struct**: Սահմանում է կարգավորումների կառուցվածքը:
*   **`Load()`**, **`Save()`**, **`Get()`**: Կարդում, պահպանում և վերադարձնում են կարգավորումները:

### 3. `internal/domain/models/`
*   **`article.go`**: Սահմանում է `Article` կառուցվածքը (Title, Content, RewrittenContent, Status, URL և այլն):
*   **`feed.go`**: Սահմանում է RSS հոսքերի մոդելները:

### 4. `internal/usecase/article_usecase.go`
Պարունակում է գլխավոր բիզնես տրամաբանությունը:
*   **`ExecuteScrapeCycle()`**: Ստուգում է RSS աղբյուրները և նոր հոդվածները պահպանում որպես "pending":
*   **`ExecuteAICycle()`**: Վերցնում է "pending" կամ "failed" հոդվածները և դիմում AI-ին (Nvidia/ModelsLab)՝ տեքստը վերաշարադրելու համար (դառնում է "rewritten"):
*   **`ExecutePublishCycle()`**: "rewritten" հոդվածները հրապարակում է WordPress-ում և նշում որպես "published":

### 5. `internal/infra/database/sqlite.go`
*   **`NewSQLiteDB(path)`**: Բացում է SQLite տվյալների բազայի կապը և ապահովում է աղյուսակների ստեղծումը:

### 6. `internal/infra/` (Ենթաթղթապանակներ)
*   **`api/`**: HTTP սերվեր և API (Fiber framework) dashboard-ի համար (օր. `server.go`):
*   **`repository/`**: `SQLiteArticleRepository`՝ հոդվածների հետ տվյալների բազայում աշխատելու համար (`sqlite_article_repo.go`):
*   **`scraper/`**: RSS և HTML էջերից տեղեկություն քաշող ֆունկցիաներ (օգտագործում է `go-rod` և `gofeed`):
*   **`publisher/`**: WordPress API-ի հետ աշխատող կոդ (`wordpress.go`) և ներքին հղումներ ավելացնող (`linker.go`):
*   **`ai/`**: Nvidia և ModelsLab AI գործիքների ինտեգրացիաներ:
*   **`utils/`**: Օժանդակ ֆունկցիաներ (լոգեր, բրաուզերի բացում, HTML մաքրում և այլն):

## Տվյալների Բազայի Սխեմա (SQLite Database: `bot_ultimate.db`)

Հիմնական աղյուսակները՝

1.  **`app_settings`**:
    *   Պահում է բոտի գլոբալ կարգավորումները:
    *   Օրինակ՝ WordPress URL, գաղտնաբառեր, AI API բանալիներ, աշխատանքի ինտերվալներ և AI հրահանգներ (prompts):
2.  **`rss_topics`**:
    *   Պահում է RSS աղբյուրների ցանկը:
    *   Կապում է յուրաքանչյուր RSS հղում WordPress-ի կոնկրետ Category ID-ի հետ:
3.  **`articles`**:
    *   Պահում է բոլոր քաշված և մշակված հոդվածները:
    *   Հիմնական սյուներ՝
        *   `source_url` (Օրիգինալ հղումը)
        *   `title`, `content` (Օրիգինալ վերնագիր և բովանդակություն)
        *   `rewritten_content` (AI-ի կողմից վերաշարադրված տեքստ)
        *   `status` (Կարգավիճակ՝ "pending", "rewritten", "published", "failed")
        *   `retry_count`, `next_retry_at` (Կրկնակի փորձերի համար)
        *   `category_id`, `image_url` և այլն:

## Ինչպես Կատարել Փոփոխություններ

*   **Նոր կարգավորում ավելացնելու համար:** Փոփոխեք `internal/config/config.go` ֆայլը, `app_settings` աղյուսակը և dashboard-ի API հարցումները (`internal/infra/api/server.go`):
*   **Scraping-ը փոխելու համար:** Նայեք `internal/infra/scraper/` թղթապանակը:
*   **AI-ի վերաշարադրման տրամաբանությունը փոխելու համար:** Փոփոխեք `internal/infra/ai/` ֆայլերը և `processWithAI` ֆունկցիան `internal/usecase/article_usecase.go`-ում:
*   **WordPress հրապարակումը փոխելու համար:** Փոփոխեք `internal/infra/publisher/wordpress.go`:
*   **Նոր ֆոնային պրոցես ավելացնելու համար:** Ավելացրեք նոր ցիկլ `cmd/bot/main.go` ֆայլում գտնվող `BotRunner`-ում:
