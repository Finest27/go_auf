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
}

func NewBotRunner(uc *usecase.ArticleUseCase) *BotRunner {
	return &BotRunner{useCase: uc}
}

func (r *BotRunner) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cancelFunc != nil {
		r.cancelFunc()
	}

	r.ctx, r.cancelFunc = context.WithCancel(context.Background())
	
	// Broadcast status update
	utils.BroadcastEvent("status_update", map[string]bool{"running": true})
	
go r.runScrapeLoop(r.ctx)
go r.runPublishLoop(r.ctx)
}

func (r *BotRunner) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cancelFunc != nil {
		r.cancelFunc()
		r.cancelFunc = nil
	}
	
	// Broadcast status update
	utils.BroadcastEvent("status_update", map[string]bool{"running": false})
}

func main() {
	utils.BroadcastLog("[KERNEL] Initializing Ultimate Engine v4.5 (Live Mode)...")

	// Make sure the path is dynamic based on the executable or current dir
	dbPath := "bot_ultimate.db"
	db := database.NewSQLiteDB(dbPath)
	config.InitDB(db)

	if err := config.Load(); err != nil {
		utils.BroadcastLog("[SYSTEM] Failed to load config, using defaults.")
	}

	articleRepo := repository.NewSQLiteArticleRepository(db)
	linker := publisher.NewInternalLinker(articleRepo)
	articleUC := usecase.NewArticleUseCase(articleRepo, linker)
	runner := NewBotRunner(articleUC)

	server := api.NewServer(runner.Start, runner.Stop, articleRepo, articleUC)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		utils.BroadcastLog("[SYSTEM] Shutting down...")
		runner.Stop()
		os.Exit(0)
	}()

	port := 8080
	url := fmt.Sprintf("http://127.0.0.1:%d", port)

	go func() {
		time.Sleep(1500 * time.Millisecond)
		utils.BroadcastLog("[SYSTEM] Launching Interface at %s", url)
		_ = utils.OpenBrowser(url)
	}()

	fmt.Printf("\n--- ThePress Bot Ultimate ---\n")
	fmt.Printf("Interface: %s\n", url)
	fmt.Printf("---------------------------\n\n")

	log.Fatal(server.Listen(fmt.Sprintf(":%d", port)))
}

func (r *BotRunner) runScrapeLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			config.Load()
			cfg := config.Get()
			utils.BroadcastLog("--- [SCRAPE CYCLE START] ---")
			r.useCase.ExecuteScrapeCycle(ctx, cfg)

			interval := time.Duration(cfg.Bot.RunIntervalMinutes) * time.Minute
			if interval < 5*time.Minute {
				interval = 5 * time.Minute
			}
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

func (r *BotRunner) runPublishLoop(ctx context.Context) {
	time.Sleep(15 * time.Second) // Stagger start
	currentTopicIndex := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
			config.Load()
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
