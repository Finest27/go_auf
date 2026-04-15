package scraper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/mmcdole/gofeed"
	"thepress_bot_go/internal/infra/utils"
)

func FetchRSSLinks(ctx context.Context, browser *rod.Browser, url string) ([]string, error) {
	utils.BroadcastLog("[RSS] Fetching feed via Stealth Browser: %s", url)

	page := PreparePage(browser)
	defer page.Close() // Prevent memory leaks for RSS pages too

	err := page.Timeout(30 * time.Second).Navigate(url)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to RSS: %w", err)
	}

	consentButtons := []string{
		"button[name='agree']",
		"button.accept-all",
		"button[value='agree']",
		".con-wizard button[type='submit']",
	}

	for _, sel := range consentButtons {
		has, _, _ := page.Has(sel)
		if has {
			if el, err := page.Element(sel); err == nil {
				utils.BroadcastLog("[RSS] Found consent wall, bypassing...")
				_ = el.Click("left", 1)
				time.Sleep(3 * time.Second)
				break
			}
		}
	}

	time.Sleep(4 * time.Second)

	content, err := page.HTML()
	if err != nil {
		return nil, fmt.Errorf("failed to get RSS content: %w", err)
	}

	// FIX: Robust XML Extraction (Ignore Chrome's HTML wrapping)
	if rssStart := strings.Index(content, "<rss"); rssStart != -1 {
		if rssEnd := strings.LastIndex(content, "</rss>"); rssEnd != -1 {
			content = content[rssStart : rssEnd+6]
		}
	} else if feedStart := strings.Index(content, "<feed"); feedStart != -1 {
		if feedEnd := strings.LastIndex(content, "</feed>"); feedEnd != -1 {
			content = content[feedStart : feedEnd+7]
		}
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(strings.NewReader(content))
	if err != nil {
		// Fallback
		feed, err = fp.ParseString(content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
		}
	}

	var links []string
	for _, item := range feed.Items {
		if item.Link != "" {
			links = append(links, item.Link)
		}
	}

	utils.BroadcastLog("[RSS] Successfully extracted %d links", len(links))
	return links, nil
}
