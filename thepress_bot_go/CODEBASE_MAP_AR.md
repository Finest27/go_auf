# Կոդերի և Ճարտարապետության Քարտեզ (ThePress Bot Ultimate)

Այս փաստաթուղթը տրամադրում է բոտի կոդի համապարփակ քարտեզ, որն ավելի հեշտ է դարձնում նավարկելը և փոփոխություններ կատարելը: Այն ուրվագծում է դիրեկտորիաների կառուցվածքը, հիմնական ֆայլերը, նրանց գործառույթները և տվյալների բազայի կառուցվածքը:

## Դիրեկտորիաների կառուցվածքը (Directory Structure)
*   **`cmd/bot/`**: Պարունակում է ծրագրի մուտքի կետը (entry point):
*   **`internal/config/`**: Կառավարում է հավելվածի կարգավորումները:
*   **`internal/domain/`**: Սահմանում է հիմնական մոդելները (models) և ինտերֆեյսները (interfaces), որոնք օգտագործվում են որպես կանոններ և պայմանագրեր համակարգի մյուս մասերի համար:
*   **`internal/infra/`**: Պարունակում է կոնկրետ իրականացումները (տվյալների բազա, արտաքին API-ներ, սքրեյփերներ / scrapers):
*   **`internal/usecase/`**: Իրականացնում է հիմնական բիզնես տրամաբանությունը (business logic), օրինակ՝ հոդվածների քաշելու, վերաշարադրելու և հրապարակելու ցիկլերը:

## Հիմնական ֆայլեր և ֆունկցիաներ (Key Files and Functions)

### 1. `cmd/bot/main.go`
Սա հավելվածի սկզբնակետն է: Այն նախաստեղծում է տվյալների բազան, բեռնում է կարգավորումները, միացնում կախվածությունները և սկսում ֆոնային ցիկլերը (background loops):
*   **`main()`**: Նախաստեղծում է բազան, կախվածությունները, API սերվերը և սկսում bot runner-ը:
*   **`BotRunner.Start()`**: Սկսում է ֆոնային scraping, AI և publishing ցիկլերը (`runScrapeLoop`, `runAILoop`, `runPublishLoop`):
*   **`BotRunner.Stop()`**: Կանգնեցնում է ֆոնային ցիկլերը:

### 2. `internal/config/config.go`
Կառավարում է հավելվածի կարգավորումները, որոնք հիմնականում պահվում են SQLite տվյալների բազայում, բայց կարող են անցնել JSON-ի:
*   **`Load()`**: Բեռնում է կարգավորումները բազայից հիշողության մեջ:
*   **`Save()`**: Պահպանում է ընթացիկ կարգավորումները բազայում:
*   **`Get()`**: Վերադարձնում է ընթացիկ կարգավորումների օբյեկտը (Config):

### 3. `internal/domain/models/`
Սահմանում է հիմնական տվյալների կառույցները (data structures), որոնք օգտագործվում են ամբողջ հավելվածում:
*   **`article.go`**: Սահմանում է `Article` կառույցը (SourceURL, Title, Content, RewrittenContent, Status, PublishDate, slug, tags, meta description և այլն):
*   **`feed.go`**: Սահմանում է RSS feed-ի մոդելները:

### 4. `internal/usecase/article_usecase.go`
Պարունակում է կենտրոնական բիզնես տրամաբանությունը, որը կառավարում է սքրեյփինգի (scraping), վերաշարադրման (rewriting) և հրապարակման (publishing) գործընթացները:
*   **`ArticleUseCase.ExecuteScrapeCycle()`**: Անցնում է կազմաձևված RSS թեմաների վրայով, գործարկում սքրեյփերը և նոր հոդվածները պահպանում բազայում՝ "pending" (սպասող) կարգավիճակով:
*   **`ArticleUseCase.ExecuteAICycle()`**: Գտնում է "pending" հոդվածները և գործարկում AI-ն բովանդակությունը վերաշարադրելու համար։
*   **`ArticleUseCase.ExecutePublishCycle()`**: Գտնում է արդեն վերաշարադրված հոդվածները և դրանք հրապարակում (publish) WordPress-ում:

### 5. `internal/infra/database/sqlite.go` (և կից ֆայլեր)
Կառավարում է SQLite տվյալների բազայի կապը և սխեման:
*   **`NewSQLiteDB(path)`**: Բացում է կապ SQLite տվյալների բազայի հետ և ստեղծում սխեման (աղյուսակները):

### 6. `internal/infra/` (Ենթադիրեկտորիաներ)
Պարունակում է արտաքին համակարգերի հետ փոխազդեցության իրականացումները:
*   **`api/` (`server.go`)**: Web սերվերի (REST/WebSocket) տրամաբանություն և HTTP handler-ներ բոտի dashboard/վահանակի համար՝ օգտագործելով Fiber framework (`github.com/gofiber/fiber/v2`):
*   **`repository/` (`sqlite_article_repo.go`)**: SQL իրականացումներ (database queries) բազայից հոդվածներ և կարգավորումներ պահպանելու/կարդալու համար:
*   **`scraper/`**: Տրամաբանություն RSS ֆիդերից XML քաշելու, որոնելու և վեբ էջերից հոդվածներ կարդալու համար: Օգտագործում է `github.com/go-rod/rod` վեբ զննարկիչի (browser) հիման վրա աշխատելու համար, և `gofeed`՝ RSS վերլուծելու համար:
*   **`publisher/` (`wordpress.go`, `linker.go`)**: Տրամաբանություն վավերացման և WordPress REST API-ի հետ փոխազդելու, հոդվածներն ու նկարները (images) ներբեռնելու/վերբեռնելու համար:
*   **`ai/` (`nvidia.go`, `modelslab.go`, `factory.go`)**: Իրականացումներ AI պրովայդերների հետ (Nvidia NIM կամ ModelsLab) շփվելու և հոդվածի տեքստը, նկարները վերամշակելու համար:
*   **`utils/`**: Օգնական ֆունկցիաներ (logger, sanitizer, downloader, browser utility), որոնք օգնում են ներբեռնել ֆայլեր, ստուգել անվտանգություն և այլն:

## Տվյալների բազայի սխեման / DB Schema (SQLite)

Բոտն օգտագործում է SQLite տվյալների բազա (հիմնականում `bot_ultimate.db`), հետևյալ հիմնական աղյուսակներով (tables):

1.  **`app_settings`**:
    *   Պահպանում է բոտի գլոբալ կարգավորումները (օր.՝ աշխատելու ինտերվալներ, WordPress-ի credentials-ներ, AI API բանալիներ):
    *   Ներառում է նաև AI պրոմպտները (`system_prompt_nvidia`, `system_prompt_modelslab`):
2.  **`rss_topics`**:
    *   Պահպանում է սքրեյփ արվող RSS ֆիդերի ցանկը:
    *   Յուրաքանչյուր ֆիդի URL կապում է կոնկրետ WordPress Category ID-ի հետ:
3.  **`articles`**:
    *   Պահպանում է բոլոր քաշված և մշակված հոդվածները:
    *   Հիմնական սյուներ (columns)՝
        *   `source_url` (Հիմնական նույնացուցիչ՝ կրկնօրինակներից խուսափելու համար)
        *   `title` (Բնօրինակ վերնագիր)
        *   `content` (Բնօրինակ բովանդակություն)
        *   `rewritten_content` (AI-ի կողմից գեներացված վերաշարադրված տեքստ)
        *   `status` (օրինակ՝ "pending", "published", "failed")
        *   `publish_date` և `created_at`
        *   SEO մետատվյալներ (`slug`, `tags`, `meta_description`, `image_alt`)
        *   `category_id` (WordPress կատեգորիայի ID)
        *   `retry_count` և `next_retry_at` (Անհաջողությունների դեպքում կրկնելու համար)

## Ինչպես կատարել փոփոխություններ (How to Make Changes)

*   **Նոր կարգավորում ավելացնելիս (Settings)**: Թարմացրեք `internal/config/config.go`, `app_settings` աղյուսակի սխեման (`sqlite.go`-ում), և API handler-ները `internal/infra/api/server.go`-ում:
*   **Scraping-ի տրամաբանությունը փոխելիս (Scraping)**: Փնտրեք և փոփոխեք կոդերը `internal/infra/scraper/` թղթապանակում (օրինակ՝ `universal.go` կամ `rss.go`):
*   **AI և վերաշարադրման տրամաբանությունը փոխելիս (AI / Rewriting)**: Թարմացրեք AI ինտեգրումը `internal/infra/ai/`-ում (Nvidia/Modelslab API կանչեր) և prompt-ի ու ցիկլի տրամաբանությունը `internal/usecase/article_usecase.go`-ում:
*   **WordPress հրապարակումը փոխելիս (Publishing)**: Փոփոխեք կոդը `internal/infra/publisher/wordpress.go`-ում:
*   **Նոր ֆոնային առաջադրանք (Background task) ավելացնելիս**: Ավելացրեք նոր ֆունկցիա `cmd/bot/main.go`-ում `BotRunner`-ի ներսում (որպես goroutine) և համապատասխան բիզնես լոգիկան՝ use case-ում:
