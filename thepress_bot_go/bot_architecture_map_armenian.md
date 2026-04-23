# ThePress Bot Ultimate Ճարտարապետության Քարտեզ

Այս փաստաթուղթը տրամադրում է բոտի կոդի բազայի ամբողջական քարտեզը՝ օգնելով հասկանալ դրա կառուցվածքը, հիմնական գործառույթները և տվյալների բազան:

## Նախագծի Կառուցվածք (Clean Architecture)

Նախագիծը գրված է Go-ով (Golang) և հետևում է Clean Architecture-ի կանոններին։

*   **`cmd/bot/`**: Հավելվածի մուտքի կետը (entry point)։
    *   `main.go`: Գործարկում է տվյալների բազան, սերվերը և ֆոնային պրոցեսները (scraping, AI, publishing loops):
*   **`internal/config/`**: Կարգավորումների ղեկավարում։
    *   `config.go`: Կարդում/Պահպանում է բոտի կարգավորումները SQLite բազայից կամ `settings.json`-ից։
*   **`internal/domain/`**: Հիմնական մոդելները և ինտերֆեյսները։
    *   `models/article.go`: `Article` կառույցը (Title, Content, RewrittenContent, Status և այլն)։
    *   `models/feed.go`: `Feed` կառույցը։
    *   `repository/interfaces.go`: Տվյալների բազայի հետ աշխատելու ինտերֆեյսները։
    *   `services/interfaces.go`: Արտաքին ծառայությունների ինտերֆեյսները (Publisher, Linker, ScraperService, AIProvider)։
*   **`internal/usecase/`**: Հիմնական բիզնես տրամաբանությունը։
    *   `article_usecase.go`: Կառավարում է հոդվածների քաշելը (scrape), AI վերաշարադրումը (rewrite) և WordPress-ում հրապարակումը։
*   **`internal/infra/`**: Ինտերֆեյսների իրականացումներ (database, API, scrapers, AI)։
    *   **`ai/`**: `modelslab.go`, `nvidia.go` - AI պրովայդերները տեքստի և նկարների մշակման համար։
    *   **`api/`**: `server.go` - Fiber-ով գրված վեբ սերվեր և REST API/WebSocket ադմինիստրատիվ վահանակի (control panel) համար։ `static/` պապկայում գտնվում են HTML/JS/CSS ֆայլերը։
    *   **`database/`**: `sqlite.go` - Տվյալների բազայի միացում և աղյուսակների ստեղծում։
    *   **`publisher/`**: `wordpress.go`, `linker.go` - WordPress-ի հետ կապը (հոդվածների և նկարների վերբեռնում) և ներքին հղումների տեղադրումը։
    *   **`repository/`**: `sqlite_article_repo.go` - `Article` մոդելի համար SQL հարցումները։
    *   **`scraper/`**: `rss.go`, `universal.go`, `wrapper.go`, `stealth.go` - RSS ալիքների կարդալը և կայքերից հոդվածների տեքստերի առանձնացումը (go-rod, go-readability)՝ շրջանցելով պաշտպանությունները (bot protection):
    *   **`utils/`**: Օժանդակ ֆունկցիաներ (լոգավորում, նկարների ներբեռնում)։

## Տվյալների Բազա (SQLite)

Ֆայլի անունը հիմնականում `bot_ultimate.db` է, որը պարունակում է հետևյալ աղյուսակները.

1.  **`app_settings`**:
    *   Պահպանում է գլոբալ կարգավորումները (WP url/username/password, AI API keys, ժամանակային միջակայքերը, AI-ի prompt-ները)։
2.  **`rss_topics`**:
    *   Պահպանում է RSS ալիքների հղումները և դրանց կապված WordPress-ի Category ID-ները։
3.  **`articles`**:
    *   Պահպանում է բոլոր հոդվածները։
    *   Հիմնական դաշտերը՝ `source_url` (եզակի), `title`, `content` (օրիգինալ), `rewritten_content` (AI-ի կողմից մշակված), `status` (pending, rewritten, published, failed), `category_id`, SEO տվյալներ։

## Գլխավոր Ֆունկցիաներ և Պրոցեսներ

*   **Scraping Loop** (`runScrapeLoop`): Պարբերաբար կարդում է RSS-ները, քաշում նոր հոդվածները և պահպանում բազայում `pending` ստատուսով։
*   **AI Loop** (`runAILoop`): Վերցնում է `pending` կամ `failed` հոդվածները, ուղարկում NVIDIA կամ ModelsLab AI-ներին՝ վերաշարադրելու և թարգմանելու համար։ Հաջողության դեպքում դնում է `rewritten` ստատուս։
*   **Publish Loop** (`runPublishLoop`): Վերցնում է `rewritten` հոդվածները և հրապարակում WordPress-ում՝ տեղադրելով նաև նկարները։ Ստատուսը դառնում է `published`։

## Հետագա Փոփոխություններ Անելու Համար

*   Բազայում նոր դաշտ ավելացնելիս՝ թարմացնել `database/sqlite.go`-ի `updates` զանգվածը։
*   UI/դիզայն փոխելիս՝ խմբագրել `internal/infra/api/static/` ֆայլերը (հիշել, որ հավելվածը պետք է վերակոմպիլացվի՝ `go build`, քանի որ static ֆայլերը ներկառուցվում են)։
*   Արտաքին փաթեթներ օգտագործելիս նախընտրել Go-ի ստանդարտ գրադարանը և արդեն առկա գրադարանները: Բազայի համար պարտադիր օգտագործել `modernc.org/sqlite` (առանց CGO)։