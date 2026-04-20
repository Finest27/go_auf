package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func isSafeURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid scheme: %s", u.Scheme)
	}

	ips, err := net.LookupIP(u.Hostname())
	if err != nil {
		return err
	}
	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
			return fmt.Errorf("access to private IP denied: %s", ip.String())
		}
	}
	return nil
}

func getSafeClient() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if err := isSafeURL(req.URL.String()); err != nil {
				return err
			}
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}
}

func ToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func DownloadFileToBytes(url string) ([]byte, error) {
	if err := isSafeURL(url); err != nil {
		return nil, fmt.Errorf("unsafe URL: %w", err)
	}

	client := getSafeClient()
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to download file: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func DownloadFileToPath(url, path string) error {
	if err := isSafeURL(url); err != nil {
		return fmt.Errorf("unsafe URL: %w", err)
	}

	client := getSafeClient()
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
