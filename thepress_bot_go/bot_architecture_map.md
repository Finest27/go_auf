# ThePress Bot Ultimate Ճարտարապետության Քարտեզ

Այս փաստաթուղթը տրամադրում է բոտի կոդի բազայի համապարփակ քարտեզը, ինչը հեշտացնում է նավարկելը և փոփոխություններ կատարելը: Այն նկարագրում է դիրեկտորիաների կառուցվածքը, հիմնական ֆայլերը, դրանց գործառույթները և տվյալների բազայի կառուցվածքը:

## Դիրեկտորիաների Կառուցվածքը

*   **`cmd/bot/`**: Պարունակում է հավելվածի մուտքի կետը (entry point):
*   **`internal/config/`**: Կառավարում է հավելվածի կարգավորումները:
*   **`internal/domain/`**: Սահմանում է հիմնական մոդելները (entities) և ինտերֆեյսները (contracts):
*   **`internal/infra/`**: Պարունակում է կոնկրետ իրականացումները (տվյալների բազա, արտաքին API-ներ, սքրեյփերներ):
*   **`internal/usecase/`**: Իրականացնում է հիմնական բիզնես տրամաբանությունը:

## Հիմնական Ֆայլերը և Ֆունկցիաները

### 1. `cmd/bot/main.go`
Սա հավելվածի սկզբնակետն է։ Այն սկզբնավորում է տվյալների բազան, բեռնում կարգավորումները, կարգավորում կախվածությունները և սկսում ֆոնային ցիկլերը:
*   **`main()`**: Սկզբնավորում է տվյալների բազան, API սերվերը և գործարկում բոտի runner-ը:
*   **`BotRunner.Start()`**: Սկսում է ֆոնային գործընթացները (`runScrapeLoop`, `runAILoop`, `runPublishLoop`):
*   **`BotRunner.Stop()`**: Կանգնեցնում է ֆոնային գործընթացները:
*   **`runScrapeLoop(ctx)`**: Պարբերաբար կանչում է `ArticleUseCase.ExecuteScrapeCycle()`:
*   **`runAILoop(ctx)`**: Պարբերաբար կանչում է `ArticleUseCase.ExecuteAICycle()`:
*   **`runPublishLoop(ctx)`**: Պարբերաբար կանչում է `ArticleUseCase.ExecutePublishCycle()`:

### 2. `internal/config/config.go`
Կառավարում է հավելվածի կարգավորումները, որոնք պահվում են SQLite տվյալների բազայում, բայց կարող են վերցվել նաև JSON-ից։
*   **`Load()`**: Բեռնում է կարգավորումները տվյալների բազայից հիշողության մեջ:
*   **`Save(cfg)`**: Պահպանում է նոր կարգավորումները տվյալների բազայում:
*   **`Get()`**: Վերադարձնում է ընթացիկ կարգավորումների օբյեկտը:

### 3. `internal/domain/models/`
Սահմանում է հիմնական տվյալների կառուցվածքները:
*   **`article.go`**: Սահմանում է `Article` կառուցվածքը (SourceURL, Title, Content, RewrittenContent, Status, PublishDate և այլն):
*   **`feed.go`**: Սահմանում է `Feed` կառուցվածքը:

### 4. `internal/domain/repository/` և `internal/domain/services/`
Սահմանում են ինտերֆեյսներ, որոնց միջոցով usecase-ը շփվում է ինֆրաստրուկտուրայի հետ (`ArticleRepository`, `AIProvider`, `Publisher`, `ScraperService`):

### 5. `internal/usecase/article_usecase.go`
Պարունակում է կենտրոնական բիզնես տրամաբանությունը, որը կառավարում է հոդվածների հավաքագրման, վերաշարադրման և հրապարակման գործընթացները:
*   **`ExecuteScrapeCycle()`**: Ստուգում է RSS ալիքները, հավաքագրում նոր հոդվածները և պահպանում բազայում որպես "pending":
*   **`ExecuteAICycle()`**: Վերցնում է չմշակված կամ ձախողված հոդվածները և ուղարկում AI-ին վերաշարադրման: Կարգավիճակը փոխվում է "rewritten":
*   **`ExecutePublishCycle()`**: Գտնում է վերաշարադրված հոդվածները և հրապարակում WordPress-ում: Կարգավիճակը փոխվում է "published":

### 6. `internal/infra/` (Ենթադիրեկտորիաներ)
*   **`database/sqlite.go`**: Կապակցում է SQLite բազայի հետ և ստեղծում սխեմաները (`NewSQLiteDB`):
*   **`repository/sqlite_article_repo.go`**: Իրականացնում է հոդվածների պահպանումը, թարմացումը և ստացումը բազայից:
*   **`ai/`**: Պարունակում է AI պրովայդերների իրականացումները (`nvidia.go`, `modelslab.go`):
*   **`publisher/`**: Պարունակում է WordPress-ի հետ կապը (`wordpress.go`, `linker.go`):
*   **`scraper/`**: Իրականացնում է կայքերից տվյալների հավաքագրումը `go-rod` և `gofeed` գրադարանների միջոցով (`wrapper.go`, `universal.go`):
*   **`api/server.go`**: Իրականացնում է Fiber վեբ սերվերը և HTTP հարցումները բոտի կառավարման վահանակի համար:

## Տվյալների Բազայի Սխեման (SQLite)

Բոտն օգտագործում է SQLite բազա (սովորաբար `bot_ultimate.db`), որն ունի հետևյալ հիմնական աղյուսակները.

### 1. `app_settings`
*   Պահպանում է բոտի գլոբալ կարգավորումները:
*   **Սյուներ**: `wp_url`, `wp_username`, `wp_app_password`, `ai_tool`, `nvidia_api_key`, `modelslab_api_key`, `min_article_len`, ինտերվալներ, ավտոհրապարակման կարգավորումներ և AI պրոմպտներ:

### 2. `rss_topics`
*   Պահպանում է RSS աղբյուրների ցանկը:
*   **Սյուներ**: `name` (Անուն), `wp_category_id` (WordPress կատեգորիայի ID), `rss_url` (RSS հասցեն):

### 3. `articles`
*   Պահպանում է բոլոր հավաքագրված և մշակված հոդվածները:
*   **Հիմնական սյուներ**:
    *   `id`, `source_url` (Բնօրինակ հղում, դուպլիկատներից խուսափելու համար)
    *   `title` (Վերնագիր), `content` (Բնօրինակ տեքստ)
    *   `rewritten_content` (AI-ի կողմից վերաշարադրված տեքստ)
    *   `status` (օրինակ՝ "pending", "rewritten", "published", "failed")
    *   `image_url`, `meta_description`, `focus_keywords`, `slug`, `tags`, `image_alt`
    *   `retry_count`, `next_retry_at` (Կրկնակի փորձերի համար)
    *   `publish_date`, `created_at`

## Ինչպես կատարել փոփոխություններ

*   **Նոր կարգավորում ավելացնելիս**: Թարմացրեք `internal/config/config.go`, `app_settings` աղյուսակը (`internal/infra/database/sqlite.go`) և `internal/infra/api/server.go` կառավարման վահանակը:
*   **Սքրեյփինգի տրամաբանությունը փոխելիս**: Աշխատեք `internal/infra/scraper/` ֆայլերի հետ:
*   **AI վերաշարադրումը փոխելիս**: Թարմացրեք `internal/infra/ai/` պրովայդերները կամ պրոմպտների տրամաբանությունը `internal/usecase/article_usecase.go`-ում:
*   **WordPress հրապարակումը փոխելիս**: Խմբագրեք `internal/infra/publisher/wordpress.go`:
*   **Նոր ֆոնային առաջադրանք ավելացնելիս**: Ավելացրեք նոր ցիկլ (loop) `cmd/bot/main.go` ֆայլում և համապատասխան ֆունկցիա `internal/usecase/`-ում: