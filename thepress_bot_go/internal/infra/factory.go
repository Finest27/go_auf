package infra

import (
	"thepress_bot_go/internal/config"
	"thepress_bot_go/internal/domain/services"
	"thepress_bot_go/internal/infra/ai"
	"thepress_bot_go/internal/infra/publisher"
	"thepress_bot_go/internal/infra/scraper"
)

type DefaultProviderFactory struct{}

func NewDefaultProviderFactory() *DefaultProviderFactory {
	return &DefaultProviderFactory{}
}

func (f *DefaultProviderFactory) CreateAI(cfg config.Config) services.AIProvider {
	return ai.NewNvidiaProvider(cfg.AI.NvidiaAPIKey, cfg.Prompts.SystemPromptNvidia)
}

func (f *DefaultProviderFactory) CreatePublisher(cfg config.Config) services.Publisher {
	var imageProvider services.AIProvider
	if cfg.AI.ModelsLabKey != "" {
		imageProvider = ai.NewModelsLabProvider(cfg.AI.ModelsLabKey)
	} else {
		imageProvider = ai.NewNvidiaProvider(cfg.AI.NvidiaAPIKey, cfg.Prompts.SystemPromptNvidia)
	}
	return publisher.NewWPClient(cfg.WordPress.URL, cfg.WordPress.Username, cfg.WordPress.AppPassword, imageProvider)
}

func (f *DefaultProviderFactory) CreateScraper() (services.ScraperService, error) {
	return scraper.NewWrapper()
}
