# ThePressUSA Auto-Journalist Bot (Golang Edition)

## 📌 Ընդհանուր Նկարագիր (Overview)
Այս նախագիծը բարձր արտադրողականությամբ (high-performance) ավտոմատացված լրագրող-բոտ է՝ գրված **Go (Golang)** լեզվով։ Բոտն ավտոմատացնում է նորությունների հավաքագրման (scraping), արհեստական բանականությամբ վերամշակման (AI rewriting), և WordPress կայքում հրապարակման գործընթացը։

Նախագիծը կառուցված է **Clean Architecture** սկզբունքներով՝ ապահովելով կոդի մաքրություն, հեշտ սպասարկում և անվտանգություն։ Ունի ներկառուցված **Live Control Panel** (Web UI)՝ բոտը իրական ժամանակում կառավարելու համար։

---

## 🏗 Ճարտարապետություն և Պանակների Քարտեզ (Architecture & Folder Structure)

```text
thepress_bot_go/
│
├── cmd/
│   └── bot/
│       └── main.go                  # Ծրագրի գլխավոր մուտքի կետը (Entry point)
│
├── internal/                        # Ներքին տրամաբանություն
│   ├── config/
│   │   └── config.go                # Կարգավորումների և բազայի սինխրոնիզացիա (app_settings)
│   │
│   ├── domain/                      # Տվյալների մոդելներ և ինտերֆեյսներ
│   │   ├── models/                  # (Article, Feed)
│   │   ├── repository/              # Բազայի հետ աշխատելու ինտերֆեյսներ
│   │   └── services/                # AI պրովայդերների ինտերֆեյսներ
│   │
│   ├── usecase/                     # Բիզնես տրամաբանություն (Core Logic)
│   │   ├── article_usecase.go       # Հիմնական ցիկլերը՝ Scrape և Publish
│   │   └── article_usecase_test.go  # Unit թեստեր
│   │
│   └── infra/                       # Արտաքին գործիքներ և ինտեգրացիաներ
│       ├── ai/
│       │   ├── nvidia.go            # Nvidia NIM LLM ինտեգրացիա
│       │   ├── modelslab.go         # ModelsLab նկարների մշակում
│       │   └── factory.go           # AI մոդելների ընտրություն
│       │
│       ├── api/
│       │   ├── server.go            # Fiber Վեբ սերվեր (BasicAuth-ով)
│       │   └── static/              # HTML/CSS/JS (Live UI)
│       │
│       ├── database/
│       │   └── sqlite.go            # SQLite բազայի ստեղծում (bot_ultimate.db)
│       │
│       ├── publisher/
│       │   ├── wordpress.go         # WordPress REST API (Media + Posts)
│       │   └── linker.go            # Ինքնաշխատ ներքին հղումներ (Internal Linker / SEO)
│       │
│       ├── repository/
│       │   └── sqlite_article_repo.go # DB հարցումներ (CRUD)
│       │
│       ├── scraper/
│       │   ├── universal.go         # GoReadability (տեքստի մաքրում) + Բրաուզեր
│       │   ├── stealth.go           # Go-rod աննկատ (stealth) քրոուլինգ
│       │   ├── rss.go               # RSS ֆիդերի կարդացում
│       │   └── scraper_test.go      # Scraper-ի թեստեր
│       │
│       └── utils/
│           ├── logger.go            # WebSocket JSON Broadcasting
│           ├── downloader.go        # Ապահով (SSRF protected) ֆայլերի ներբեռնում
│           ├── sanitizer.go         # XSS-ի դեմ (bluemonday) HTML մաքրում
│           └── browser.go           # Բրաուզերի օժանդակ ֆունկցիաներ
```

---

## 🗄 Տվյալների Բազա (Database Schema)
Բոտն աշխատում է **SQLite**-ով (`bot_ultimate.db`): Ունի 3 հիմնական աղյուսակ.

1. **`app_settings`**: Պահում է բոլոր կարգավորումները (WP URL, App Password, Nvidia API Key, թայմերները, AI Prompts): UI-ը և բազան գործում են սինխրոն։
2. **`rss_topics`**: Պահում է RSS աղբյուրները և համապատասխան WordPress Category ID-ները։
3. **`articles`**: Հիմնական աշխատանքային աղյուսակն է.
   - `status`: Կարող է լինել `pending` (նոր), `rewritten` (AI-ի կողմից մշակված), `published` (հրապարակված), `failed` (խափանված)։
   - `retry_count` և `next_retry_at`: Օգտագործվում է ցանցային սխալների դեպքում (Exponential Backoff):

---

## 🔒 Անվտանգություն (Security Features)
Բոտը անցել է Senior QA և White-Hat Security Audit.
- **Basic Auth:** Վեբ վահանակը և API-ն պաշտպանված են գաղտնաբառով (լռելյայն՝ `admin` / `admin`)։
- **SSRF Պաշտպանություն:** `downloader.go`-ում արգելված է նկարներ քաշել ներքին IP-ներից (Localhost, Private IPs):
- **DOM XSS Պաշտպանություն:** `script.js`-ում բոլոր դինամիկ տվյալներն անցնում են `escapeHTML()` մաքրում։
- **Stored XSS Պաշտպանություն:** `universal.go`-ում կիրառված է `bluemonday.UGCPolicy()`, որպեսզի հոդվածների մեջ չպահվեն վտանգավոր `<script>` կամ `<iframe>` տեգեր։

---

## ⚡ Իրական Ժամանակի UI (Live Mode)
Կառավարման էջը (`http://127.0.0.1:8080`) գործում է **WebSocket** տեխնոլոգիայով։ 
Ամեն անգամ, երբ բազայում նոր հոդված է ավելանում, կամ հոդվածը հրապարակվում է, Բեքենդը (Go) ուղարկում է JSON `queue_update` իրադարձություն, և էկրանի աղյուսակը **ակնթարթորեն** թարմանում է՝ առանց էջը (refresh) անելու:

---

## 🚀 Գործարկման Ուղեցույց (How to Run)

### 1. Կոմպիլյացիա (Build)
Բացեք տերմինալը նախագծի պանակում և գրեք.
```bash
go build -o bot.exe ./cmd/bot/main.go
```

### 2. Միացում (Run)
```bash
.\bot.exe
```

### 3. Վահանակ (Dashboard)
Բրաուզերով մուտք գործեք.
👉 **http://127.0.0.1:8080**
- **Username:** `admin`
- **Password:** `admin`

*(Կարգավորումներ (Settings) բաժնում լրացրեք Ձեր WordPress և API բանալիները, պահպանեք և սեղմեք «Միացնել» գլխավոր էջում)*։
