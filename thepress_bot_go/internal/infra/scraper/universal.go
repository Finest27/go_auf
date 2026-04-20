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
	"thepress_bot_go/internal/infra/utils"
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

	// SAAS OPTIMIZATION: Increased timeout for better reliability
	err := page.Context(ctx).Timeout(60 * time.Second).Navigate(rawURL)
	if err != nil {
		return "", "", "", fmt.Errorf("navigate failed: %w", err)
	}

	// Advanced Evasion: Polling for main content selectors instead of arbitrary wait
	// This reduces scraping time significantly for fast pages and ensures we wait for React/Vue hydration
	selectorsToPoll := []string{"article", "main", ".post-content", ".article-body", ".entry-content"}
	pollCtx, cancelPoll := context.WithTimeout(ctx, 15*time.Second)
	defer cancelPoll()
	
	for _, sel := range selectorsToPoll {
		if _, err := page.Context(pollCtx).Element(sel); err == nil {
			utils.BroadcastLog("[SCRAPER] Selective element polling found: %s", sel)
			break
		}
	}
	
	// Slight buffer for image loading inside the found element
	time.Sleep(1 * time.Second)

	html, err := page.HTML()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get HTML: %w", err)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse url: %w", err)
	}

	// Use readability to extract the CORE content (removes menus, footers, ads)
	article, err := readability.FromReader(strings.NewReader(html), parsedURL)
	if err != nil {
		return "", "", "", fmt.Errorf("readability failed: %w", err)
	}

	if len(article.TextContent) < 100 {
		return "", "", "", fmt.Errorf("extracted content is too thin")
	}

	// Security/Quality Gate: Block Bot Protection & Access Denied pages
	lowerContent := strings.ToLower(article.TextContent)
	lowerTitle := strings.ToLower(article.Title)
	
	blockedPhrases := []string{
		"access to this page has been denied",
		"please enable javascript and cookies",
		"checking your browser before accessing",
		"cloudflare ray id",
		"you have been blocked",
		"security by cloudflare",
		"verify you are human",
		"enable javascript",
		"access denied",
		"attention required!",
	}

	for _, phrase := range blockedPhrases {
		if strings.Contains(lowerContent, phrase) || strings.Contains(lowerTitle, phrase) {
			return "", "", "", fmt.Errorf("blocked by anti-bot protection (detected phrase: %s)", phrase)
		}
	}

	// XSS Prevention: Sanitize the HTML to remove scripts/iframes but keep images and formatting
	p := bluemonday.UGCPolicy()
	safeContent := p.Sanitize(article.Content)

	return article.Title, safeContent, article.Image, nil
}
