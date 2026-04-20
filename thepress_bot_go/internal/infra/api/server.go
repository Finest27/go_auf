package api

import (
	"context"
	"embed"
	"net/http"
	"os"
	"strconv"
	"sync"
	"thepress_bot_go/internal/config"
	"thepress_bot_go/internal/infra/repository"
	"thepress_bot_go/internal/infra/utils"
	"thepress_bot_go/internal/usecase"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
	app          *fiber.App
	isBotRunning bool
	mu           sync.Mutex
	onStart      func()
	onStop       func()
	articleRepo  *repository.SQLiteArticleRepository
	useCase      *usecase.ArticleUseCase
}

func NewServer(onStart, onStop func(), repo *repository.SQLiteArticleRepository, uc *usecase.ArticleUseCase) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(recover.New())

	// Implement Basic Auth with environment variables
	adminUser := os.Getenv("ADMIN_USER")
	if adminUser == "" {
		adminUser = "admin"
	}
	adminPass := os.Getenv("ADMIN_PASS")
	if adminPass == "" {
		adminPass = "admin"
	}

	app.Use(basicauth.New(basicauth.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/ws/logs"
		},
		Users: map[string]string{
			adminUser: adminPass,
		},
		Realm: "ThePressUSA Control Panel",
	}))

	s := &Server{
		app:         app,
		onStart:     onStart,
		onStop:      onStop,
		articleRepo: repo,
		useCase:     uc,
	}

	api := app.Group("/api")
	api.Get("/bot/status", s.handleStatus)
	api.Post("/bot/toggle", s.handleToggle)
	api.Get("/settings", s.handleGetSettings)
	api.Post("/settings", s.handlePostSettings)
	api.Get("/analytics", s.handleGetAnalytics)

	api.Get("/queue", s.handleGetQueue)
	api.Get("/queue/:id", s.handleGetQueueItem)
	api.Put("/queue/:id", s.handleUpdateQueueItem)
	api.Post("/queue/publish", s.handlePublishItem)
	api.Post("/queue/delete", s.handleDeleteItem)
	api.Post("/queue/clear", s.handleClearQueue)

	app.Get("/ws/logs", websocket.New(func(c *websocket.Conn) {
		logChan := utils.Hub.Register()
		defer utils.Hub.Unregister(logChan)
		for msg := range logChan {
			if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				break
			}
		}
	}))

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(staticFiles),
		PathPrefix: "static",
		Browse:     false,
	}))

	return s
}

func (s *Server) SetBotRunning(running bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isBotRunning = running
}

func (s *Server) handleGetQueueItem(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	art, err := s.articleRepo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Article not found"})
	}
	return c.JSON(art)
}

func (s *Server) handleUpdateQueueItem(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	art, err := s.articleRepo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Article not found"})
	}

	var body struct {
		Title            string `json:"title"`
		RewrittenContent string `json:"rewritten_content"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	art.Title = body.Title
	art.RewrittenContent.String = body.RewrittenContent
	art.RewrittenContent.Valid = true

	if err := s.articleRepo.Update(c.Context(), art); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func (s *Server) handlePublishItem(c *fiber.Ctx) error {
	var body struct {
		ID int64 `json:"id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	art, err := s.articleRepo.GetByID(c.Context(), body.ID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Article not found"})
	}

	cfg := config.Get()
	go func() {
		bgCtx := context.Background()
		utils.BroadcastLog("[MANUAL] Publishing article #%d: %s", art.ID, art.Title)
		s.useCase.PublishSingle(bgCtx, cfg, art)
	}()

	return c.JSON(fiber.Map{"status": "processing"})
}

func (s *Server) handleGetAnalytics(c *fiber.Ctx) error {
	pub, pend, errs, err := s.articleRepo.GetStats(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"total_published": pub,
		"pending_queue":   pend,
		"errors":          errs,
	})
}

func (s *Server) handleGetQueue(c *fiber.Ctx) error {
	articles, err := s.articleRepo.GetPending(c.Context(), 50)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(articles)
}

func (s *Server) handleDeleteItem(c *fiber.Ctx) error {
	var body struct {
		ID int64 `json:"id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}
	err := s.articleRepo.Delete(c.Context(), body.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success"})
}

func (s *Server) handleClearQueue(c *fiber.Ctx) error {
	err := s.articleRepo.ClearQueue(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(200)
}

func (s *Server) handleStatus(c *fiber.Ctx) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return c.JSON(fiber.Map{"running": s.isBotRunning})
}

func (s *Server) handleToggle(c *fiber.Ctx) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isBotRunning {
		s.isBotRunning = false
		s.onStop()
		return c.JSON(fiber.Map{"status": "stopped"})
	} else {
		s.isBotRunning = true
		go s.onStart()
		return c.JSON(fiber.Map{"status": "started"})
	}
}

func (s *Server) handleGetSettings(c *fiber.Ctx) error {
	if err := config.Load(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(config.Get())
}

func (s *Server) handlePostSettings(c *fiber.Ctx) error {
	var cfg config.Config
	if err := c.BodyParser(&cfg); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	if _, err := url.ParseRequestURI(cfg.WordPress.URL); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid WordPress URL"})
	}
	if err := config.Save(cfg); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success"})
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}
