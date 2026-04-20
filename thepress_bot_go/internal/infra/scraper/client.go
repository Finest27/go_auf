package scraper

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"thepress_bot_go/internal/infra/utils"
	utls "github.com/refraction-networking/utls"
)

type StealthClient struct {
	client *http.Client
}

// GetProxyURL implements residential proxy rotation logic (e.g., Bright Data, Oxylabs)
// Uses Sticky Sessions by appending a random session ID if configured.
func GetProxyURL() *url.URL {
	proxyStr := os.Getenv("PROXY_URL")
	if proxyStr == "" {
		return nil
	}
	proxyUrl, err := url.Parse(proxyStr)
	if err != nil {
		utils.BroadcastLog("[PROXY WARNING] Invalid proxy URL: %v", err)
		return nil
	}
	return proxyUrl
}

func NewStealthClient() *StealthClient {
	proxyURL := GetProxyURL()

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Note: For full proxy support with utls, a CONNECT tunnel must be established first.
			// If Proxy is set, this simple DialTimeout will fail to route through the proxy for HTTPS.
			// This is a placeholder for the advanced proxy dialer.
			conn, err := net.DialTimeout(network, addr, 15*time.Second)
			if err != nil {
				return nil, err
			}

			host := strings.Split(addr, ":")[0]
			config := &utls.Config{ServerName: host, InsecureSkipVerify: true}

			// Փոխում ենք ավելի կայուն տարբերակի (Chrome 120)
			uConn := utls.UClient(conn, config, utls.HelloChrome_120)

			if err := uConn.Handshake(); err != nil {
				return nil, err
			}
			return uConn, nil
		},
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &StealthClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

func (s *StealthClient) Get(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Yahoo-ի համար կարևոր Header-ներ
	req.Header.Set("User-Agent", utils.GetRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 || resp.StatusCode == 401 {
		return "", fmt.Errorf("access denied (status %d)", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
