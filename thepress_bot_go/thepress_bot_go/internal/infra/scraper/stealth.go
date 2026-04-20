package scraper

import (
	"fmt"
	"thepress_bot_go/internal/infra/utils"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

func NewStealthBrowser() (*rod.Browser, error) {
	l := launcher.New().
		Headless(true).
		Leakless(false).
		NoSandbox(true).
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-dev-shm-usage").
		Set("disable-gpu")

	url, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch stealth browser: %w", err)
	}

	browser := rod.New().ControlURL(url).MustConnect()
	return browser, nil
}

func PreparePage(b *rod.Browser) *rod.Page {
	ua := utils.GetRandomUserAgent()

	// Simplified: Use the browser directly to avoid Incognito context leaks
	// since we already create a fresh browser for every cycle.
	page := stealth.MustPage(b)

	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent:      ua,
		AcceptLanguage: "en-US,en;q=0.9",
	})
	page.MustSetViewport(1920, 1080, 1, false)
	return page
}
