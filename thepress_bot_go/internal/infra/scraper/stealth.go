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

	// Advanced fingerprint spoofing (WebGL, Canvas, AudioContext)
	spoofScript := `
		// Spoof WebGL
		const getParameter = WebGLRenderingContext.getParameter;
		WebGLRenderingContext.prototype.getParameter = function(parameter) {
			if (parameter === 37445) return 'Intel Inc.';
			if (parameter === 37446) return 'Intel Iris OpenGL Engine';
			return getParameter(parameter);
		};
		// Spoof AudioContext
		if (window.AudioContext) {
			const origGetChannelData = AudioBuffer.prototype.getChannelData;
			AudioBuffer.prototype.getChannelData = function() {
				const results = origGetChannelData.apply(this, arguments);
				for (let i = 0; i < results.length; i += 100) {
					results[i] = results[i] + (Math.random() * 0.0001 - 0.00005);
				}
				return results;
			};
		}
	`
	page.MustEvalOnNewDocument(spoofScript)

	return page
}
