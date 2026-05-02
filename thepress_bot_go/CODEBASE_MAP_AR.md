# ThePress Bot Ultimate - Ճարտարապետական Քարտեզ

Այս փաստաթուղթը ներկայացնում է բոտի կոդերի բազայի (codebase) ամբողջական կառուցվածքը: Այն օգնում է հասկանալ պանակների (directories) և ֆայլերի նշանակությունը, ինչպես նաև հիմնական ֆունկցիաներն ու տվյալների բազայի կառուցվածքը: Կոդը գրված է Գո (Go) ծրագրավորման լեզվով՝ հետևելով "Clean Architecture" սկզբունքներին:

## Պանակների կառուցվածք (Directory Structure)

*   **`cmd/bot/`**: Ծրագրի մուտքի կետը (entry point): Պարունակում է `main.go` ֆայլը, որտեղից սկսվում է բոտի աշխատանքը:
*   **`internal/`**: Հիմնական լոգիկան, որը բաժանված է շերտերի.
    *   **`config/`**: Կարգավորումների (settings) ղեկավարում:
    *   **`domain/`**: Դոմենային շերտ. մոդելներ (entities) և ինտերֆեյսներ:
    *   **`usecase/`**: Բիզնես լոգիկայի շերտ (Business Logic / Orchestration):
    *   **`infra/`**: Ենթակառուցվածքներ (Infrastructure). տվյալների բազա, արտաքին API-ներ, սքրեյփերներ:

## Հիմնական ֆայլեր և դրանց նշանակությունը

### 1. Գլխավոր Մուտք և Ցիկլեր (Entry point)

*   **`cmd/bot/main.go`**: Սկզբնավորում է տվյալների բազան (SQLite), կարդում է կարգավորումները, միացնում է վեբ-սերվերը և գործարկում է գլխավոր ցիկլերը (`runScrapeLoop`, `runAILoop`, `runPublishLoop`):

### 2. Կարգավորումներ (Configuration)

*   **`internal/config/config.go`**: Բեռնում և պահպանում է ծրագրի կարգավորումները տվյալների բազայի `app_settings` աղյուսակում (WordPress API, AI Keys, ինտերվալներ և այլն):

### 3. Դոմեն (Domain)

*   **`internal/domain/models/article.go`**: `Article` կառույցը, որը ներկայացնում է հոդվածի մոդելը բազայում:
*   **`internal/domain/models/feed.go`**: RSS հոսքի մոդելը:
*   **`internal/domain/repository/interfaces.go`**: Տվյալների բազայի հետ աշխատելու ինտերֆեյսները:
*   **`internal/domain/services/interfaces.go` & `ai_interface.go`**: Արտաքին ծառայությունների ինտերֆեյսներ (AI, Scraper, Publisher, Linker):

### 4. Բիզնես Լոգիկա (Use Case)

*   **`internal/usecase/article_usecase.go`**: Ղեկավարում է գործընթացները:
    *   `ExecuteScrapeCycle()`: Կարդում է RSS-ը, գտնում նոր հղումներ, սկրեյփ է անում (scrape) և պահպանում բազայում որպես `pending`:
    *   `ExecuteAICycle()`: Գտնում է չմշակված հոդվածները, ուղարկում է AI-ին վերաշարադրման և թարմացնում ստատուսը `rewritten`:
    *   `ExecutePublishCycle()`: Վերաշարադրված հոդվածները ուղարկում է WordPress-ին և հրապարակում:

### 5. Ենթակառուցվածք (Infrastructure)

*   **`internal/infra/database/sqlite.go`**: SQLite տվյալների բազայի միացումը և աղյուսակների ստեղծումը (Migrations):
*   **`internal/infra/repository/sqlite_article_repo.go`**: Հոդվածների պահպանում, թարմացում և ընթերցում բազայից (SQL հարցումներ):
*   **`internal/infra/scraper/`**: Հոդվածների և RSS-ի սկրեյփինգ: Օգտագործում է `go-rod` գրադարանը բրաուզերի ավտոմատացման համար:
    *   `universal.go`: Հոդվածի տեքստի և նկարի առանձնացման լոգիկան (օգտագործելով `go-readability`):
    *   `rss.go`: RSS հոսքից հղումների առանձնացում:
*   **`internal/infra/ai/`**: AI ծառայությունների ինտեգրացիաներ:
    *   `nvidia.go`: Nvidia NIM ինտեգրացիա (օգտագործում է Llama մոդելներ) տեքստի վերաշարադրման և թարգմանության համար:
    *   `modelslab.go`: ModelsLab ինտեգրացիա նկարների մշակման (Image processing) համար:
    *   `factory.go`: AI ծառայությունների և Publisher-ի ֆաբրիկան:
*   **`internal/infra/publisher/`**:
    *   `wordpress.go`: WordPress REST API-ի հետ աշխատանք, նկարների վերբեռնում և հոդվածների հրապարակում (Publishing):
    *   `linker.go`: Ներքին հղումների ավելացում հոդվածի վերջնական տեքստում:
*   **`internal/infra/api/`**:
    *   `server.go`: Վեբ դաշբորդի և հեռակառավարման API-ի ապահովում (օգտագործում է Fiber վեբ-ֆրեյմվորքը և WebSocket լոգերի համար):
    *   `static/`: HTML, CSS և JS ֆայլեր դաշբորդի (Control Panel) ինտերֆեյսի համար (`index.html`, `script.js`):
*   **`internal/infra/utils/`**: Օժանդակ ֆունկցիաներ:
    *   `logger.go`: Լոգերի գեներացում և փոխանցում WebSocket-ով վեբ ինտերֆեյսին:
    *   `sanitizer.go`: HTML և JSON մաքրման ֆունկցիաներ (օրինակ՝ AI-ից եկած պատասխանների համար):
    *   `browser.go`: Բրաուզերի բացում օպերացիոն համակարգում (Dashboard-ը բացելու համար):

## Տվյալների Բազայի Կառուցվածք (Database Schema)

Բոտը օգտագործում է տեղական SQLite բազա՝ հիմնականում `bot_ultimate.db` ֆայլը: Հիմնական աղյուսակները հետևյալն են.

1.  **`app_settings`**:
    *   Պահում է բոտի կարգավորումները (միայն 1 տողով):
    *   Դաշտեր՝ `wp_url`, `wp_username`, `wp_app_password`, `ai_tool`, `nvidia_api_key`, `modelslab_api_key`, `system_prompt_nvidia`, ինտերվալներ և այլն:

2.  **`rss_topics`**:
    *   Պահում է հավաքագրվող RSS հոսքերի ցանկը:
    *   Դաշտեր՝ `id`, `name`, `wp_category_id` (WordPress կատեգորիայի ID-ն, որտեղ պետք է հրապարակվի), `rss_url`:

3.  **`articles`**:
    *   Պահում է բոլոր ներբեռնված, մշակված և հրապարակված հոդվածները:
    *   Կարևոր դաշտեր.
        *   `source_url`: Օրիգինալ հղումը (եզակի - unique որպեսզի կրկնօրինակներ չհավաքի):
        *   `title`, `content`: Օրիգինալ վերնագիր և տեքստ:
        *   `rewritten_content`: AI-ի կողմից ստեղծված և վերաշարադրված տեքստը:
        *   `status`: Հոդվածի կարգավիճակը (`pending`, `rewritten`, `published`, `failed`):
        *   `image_url`: Գլխավոր նկարի հղումը (կամ դրա լոկալ պահպանված տարբերակը):
        *   SEO դաշտեր՝ `meta_description`, `focus_keywords`, `slug`, `tags`, `image_alt`:
        *   `category_id`: Կապված WordPress կատեգորիան:
        *   `retry_count`, `next_retry_at`: Սխալների դեպքում կրկնելու քանակը և հաջորդ փորձի ժամանակը:
        *   `publish_date`, `created_at`: Հրապարակման և ստեղծման ամսաթվեր:

## Հիմնական Աշխատանքային Հոսք (Data Flow)

1.  **Scraper Loop**: Ծրագիրը ըստ սահմանված ինտերվալի կարդում է `rss_topics`-ի հղումները $\\rightarrow$ վերցնում է նոր հոդվածների հղումները $\\rightarrow$ ներբեռնում է հոդվածի տեքստն ու նկարը $\\rightarrow$ պահպանում է `articles` աղյուսակում **`pending`** ստատուսով:
2.  **AI Loop**: Անընդհատ ստուգում է `pending` կամ `failed` հոդվածների առկայությունը $\\rightarrow$ դրանք ուղարկում է արտաքին AI API (օր. Nvidia NIM) $\\rightarrow$ AI-ը թարգմանում, վերաշարադրում և ստեղծում է SEO մետադատա $\\rightarrow$ պահպանվում է բազայում, և ստատուսը դառնում է **`rewritten`**:
3.  **Publish Loop**: Եթե միացված է ավտոմատ հրապարակումը, վերցնում է **`rewritten`** ստատուսով հոդվածները $\\rightarrow$ մշակում է նկարները, ավելացնում է ներքին հղումները (linker) $\\rightarrow$ API-ով ներբեռնում է WordPress $\\rightarrow$ հրապարակում է հոդվածը $\\rightarrow$ բազայում ստատուսը փոխում է **`published`**: