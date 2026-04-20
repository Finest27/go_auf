package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"thepress_bot_go/internal/config"
	"thepress_bot_go/internal/infra"
	"thepress_bot_go/internal/infra/api"
	"thepress_bot_go/internal/infra/database"
	"thepress_bot_go/internal/infra/publisher"
	"thepress_bot_go/internal/infra/repository"
	"thepress_bot_go/internal/infra/utils"
	"thepress_bot_go/internal/usecase"
)

type BotRunner struct {
	mu         sync.Mutex
	cancelFunc context.CancelFunc
	ctx        context.Context
	useCase    *usecase.ArticleUseCase
	server     *api.Server
}

func NewBotRunner(uc *usecase.ArticleUseCase) *BotRunner {
	return &BotRunner{useCase: uc}
}

func (r *BotRunner) SetServer(s *api.Server) {
	r.server = s
}

func (r *BotRunner) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cancelFunc != nil {
		r.cancelFunc()
	}

	r.ctx, r.cancelFunc = context.WithCancel(context.Background())

	if r.server != nil {
		r.server.SetBotRunning(true)
	}

	// Broadcast status update
	utils.BroadcastEvent("status_update", map[string]bool{"running": true})

	go r.runScrapeLoop(r.ctx)
	go r.runAILoop(r.ctx)
	go r.runPublishLoop(r.ctx)
}

func (r *BotRunner) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cancelFunc != nil {
		r.cancelFunc()
		r.cancelFunc = nil
	}

	if r.server != nil {
		r.server.SetBotRunning(false)
	}

	// Broadcast status update
	utils.BroadcastEvent("status_update", map[string]bool{"running": false})
}

func main() {
	utils.BroadcastLog("[KERNEL] Initializing Ultimate Engine v4.5 (Live Mode)...")

	// Make sure the path is dynamic based on the executable or current dir
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "bot_ultimate.db"
	}
	db := database.NewSQLiteDB(dbPath)
	config.InitDB(db)

	if err := config.Load(); err != nil {
		utils.BroadcastLog("[SYSTEM] Failed to load config: %v", err)
	}

	articleRepo := repository.NewSQLiteArticleRepository(db)
	linker := publisher.NewInternalLinker(articleRepo)
	providerFactory := infra.NewDefaultProviderFactory()
	articleUC := usecase.NewArticleUseCase(articleRepo, linker, providerFactory)
	runner := NewBotRunner(articleUC)

	server := api.NewServer(runner.Start, runner.Stop, articleRepo, articleUC)
	// auto-start is disabled via user request
	runner.SetServer(server)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		utils.BroadcastLog("[SYSTEM] Shutting down...")
		runner.Stop()
		os.Exit(0)
	}()

	portStr := os.Getenv("PORT")
	port := 8080
	if portStr != "" {
		fmt.Sscanf(portStr, "%d", &port)
	}
	url := fmt.Sprintf("http://127.0.0.1:%d", port)

	go func() {
		time.Sleep(1500 * time.Millisecond)
		utils.BroadcastLog("[SYSTEM] Launching Interface at %s", url)
		if err := utils.OpenBrowser(url); err != nil {
			utils.BroadcastLog("[SYSTEM] Failed to open browser: %v", err)
		}
	}()

	fmt.Printf("\n--- ThePress Bot Ultimate ---\n")
	fmt.Printf("Interface: %s\n", url)
	fmt.Printf("---------------------------\n\n")

	log.Fatal(server.Listen(fmt.Sprintf(":%d", port)))
}

func (r *BotRunner) runScrapeLoop(ctx context.Context) {
        currentTopicIndex := 0
        for {
                select {
                case <-ctx.Done():
                        return
                default:
                        // Optimization: only load config if needed or with some interval     
                        cfg := config.Get()
                        utils.BroadcastLog("--- [SCRAPE CYCLE START] ---")
                        currentTopicIndex = r.useCase.ExecuteScrapeCycle(ctx, cfg, currentTopicIndex)

			intervalMinutes := cfg.Bot.RunIntervalMinutes
			if intervalMinutes < 5 {
				intervalMinutes = 5
			}
			interval := time.Duration(intervalMinutes) * time.Minute
			utils.BroadcastLog("--- [SCRAPE CYCLE END] Sleeping %v ---", interval)

			timer := time.NewTimer(interval)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}
		}
	}
}

func (r *BotRunner) runAILoop(ctx context.Context) {
	time.Sleep(5 * time.Second) // Stagger start

	for {
		select {
		case <-ctx.Done():
			return
		default:
			cfg := config.Get()
			r.useCase.ExecuteAICycle(ctx, cfg)

			// Polling interval for AI processing
			interval := 15 * time.Second
			timer := time.NewTimer(interval)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}
		}
	}
}

func (r *BotRunner) runPublishLoop(ctx context.Context) {
	time.Sleep(15 * time.Second) // Stagger start
	currentTopicIndex := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
			cfg := config.Get()

			if cfg.Bot.AutoPublish {
				utils.BroadcastLog("--- [PUBLISH CYCLE START] ---")
				currentTopicIndex = r.useCase.ExecutePublishCycle(ctx, cfg, currentTopicIndex)
			}

			pubInt := cfg.Bot.PublishIntervalMinutes
			if pubInt < 1 {
				pubInt = 15
			}

			interval := time.Duration(pubInt) * time.Minute
			utils.BroadcastLog("--- [PUBLISH CYCLE END] Next article in %v ---", interval)

			timer := time.NewTimer(interval)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}
		}
	}
}

