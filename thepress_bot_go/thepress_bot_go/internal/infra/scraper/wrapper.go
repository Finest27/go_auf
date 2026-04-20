package scraper

import (
	"context"
	"thepress_bot_go/internal/domain/services"

	"github.com/go-rod/rod"
)

// Wrapper implements services.ScraperService
type Wrapper struct {
	browser *rod.Browser
	univ    *UniversalScraper
}

func NewWrapper() (services.ScraperService, error) {
	browser, err := NewStealthBrowser()
	if err != nil {
		return nil, err
	}
	return &Wrapper{
		browser: browser,
		univ:    NewUniversalScraper(browser),
	}, nil
}

func (w *Wrapper) FetchRSSLinks(ctx context.Context, url string) ([]string, error) {
	return FetchRSSLinks(ctx, w.browser, url)
}

func (w *Wrapper) ScrapeArticle(ctx context.Context, url string) (title string, content string, image string, err error) {
	return w.univ.Scrape(ctx, url)
}

func (w *Wrapper) Close() {
	if w.browser != nil {
		w.browser.Close()
	}
}
