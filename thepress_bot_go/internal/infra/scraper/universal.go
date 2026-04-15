package scraper
import (
	"context"
	"fmt"
	"strings"
	"time"
	"net/url"

	"github.com/go-rod/rod"
	"github.com/go-shiori/go-readability"
	"github.com/microcosm-cc/bluemonday"
)

type UniversalScraper struct {
	browser *rod.Browser
}

func NewUniversalScraper(b *rod.Browser) *UniversalScraper {
	return &UniversalScraper{browser: b}
}

func (s *UniversalScraper) Scrape(ctx context.Context, rawURL string) (string, string, string, error) {
	if s.browser == nil {
		return "", "", "", fmt.Errorf("stealth browser not initialized")
	}

	page := PreparePage(s.browser)
	// Ensure page is always closed to prevent memory leaks
	defer page.Close()

	// SAAS OPTIMIZATION: Reduced timeout for better throughput
	// 45 seconds is enough for most sites to render critical JS
	err := page.Context(ctx).Timeout(45 * time.Second).Navigate(rawURL)
	if err != nil {
		return "", "", "", fmt.Errorf("navigate failed: %w", err)
	}

	// Optimization: Only wait for the page to be ready, don't wait for all images/ads
	_ = page.WaitDOMStable(1*time.Second, 0.1)

	html, err := page.HTML()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get HTML: %w", err)
	}

	parsedURL, _ := url.Parse(rawURL)

	// Use readability to extract the CORE content (removes menus, footers, ads)
	article, err := readability.FromReader(strings.NewReader(html), parsedURL)
	if err != nil {
		return "", "", "", fmt.Errorf("readability failed: %w", err)
	}

	if len(article.TextContent) < 300 {
		return "", "", "", fmt.Errorf("extracted content is too thin")
	}

	// XSS Prevention: Sanitize the HTML to remove scripts/iframes but keep images and formatting
	p := bluemonday.UGCPolicy()
	safeContent := p.Sanitize(article.Content)

	return article.Title, safeContent, article.Image, nil
}
