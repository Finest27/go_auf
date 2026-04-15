package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"thepress_bot_go/internal/domain/models"
	"thepress_bot_go/internal/domain/services"
	"thepress_bot_go/internal/infra/utils"
)

type WPClient struct {
	BaseURL  string
	Username string
	Password string
	Client   *http.Client
	AI       services.AIProvider
}

func NewWPClient(url, user, pass string, ai services.AIProvider) *WPClient {
	return &WPClient{
		BaseURL:  strings.TrimSuffix(url, "/"),
		Username: user,
		Password: pass,
		Client:   &http.Client{Timeout: 120 * time.Second},
		AI:       ai,
	}
}

func (wp *WPClient) doRequest(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(wp.Username, wp.Password)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	return wp.Client.Do(req)
}

func (wp *WPClient) UploadImageFromURL(imageURL, altText string) (int, error) {
	if imageURL == "" {
		return 0, nil
	}
	temp := filepath.Join(os.TempDir(), fmt.Sprintf("botfeat_%d.jpg", time.Now().UnixNano()))

	var err error
	for i := 0; i < 2; i++ {
		if err = utils.DownloadFileToPath(imageURL, temp); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return 0, err
	}
	defer os.Remove(temp)

	if wp.AI != nil {
		imgBytes, _ := os.ReadFile(temp)
		utils.BroadcastLog("[AI] Մշակվում է նկարը...")
		processed, err := wp.AI.ProcessImage(context.Background(), imgBytes, "Clean news photo")
		if err == nil {
			os.WriteFile(temp, processed, 0644)
		}
	}

	id, _, err := wp.UploadImageToMediaLibrary(temp, altText)
	return id, err
}

func (wp *WPClient) UploadImageToMediaLibrary(imagePath, altText string) (int, string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(imagePath))
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest("POST", wp.BaseURL+"/wp-json/wp/v2/media", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(imagePath)))

	resp, err := wp.doRequest(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	var res struct {
		ID        int    `json:"id"`
		SourceURL string `json:"source_url"`
	}
	json.NewDecoder(resp.Body).Decode(&res)

	if altText != "" {
		buf := &bytes.Buffer{}
		json.NewEncoder(buf).Encode(map[string]string{"alt_text": altText})
		reqUpdate, _ := http.NewRequest("POST", fmt.Sprintf("%s/wp-json/wp/v2/media/%d", wp.BaseURL, res.ID), buf)
		reqUpdate.Header.Set("Content-Type", "application/json")
		wp.doRequest(reqUpdate)
	}

	return res.ID, res.SourceURL, nil
}

func (wp *WPClient) ProcessInlineImages(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html, err
	}

	type uploadResult struct {
		selection *goquery.Selection
		newURL    string
	}
	results := make(chan uploadResult, 10)
	var wg sync.WaitGroup

	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if src == "" || strings.Contains(src, wp.BaseURL) {
			return
		}
		wg.Add(1)
		go func(sel *goquery.Selection, u string) {
			defer wg.Done()
			temp := filepath.Join(os.TempDir(), fmt.Sprintf("bimg_%d.jpg", time.Now().UnixNano()))
			if err := utils.DownloadFileToPath(u, temp); err == nil {
				defer os.Remove(temp)
				if id, nu, err := wp.UploadImageToMediaLibrary(temp, ""); err == nil && id > 0 {
					results <- uploadResult{sel, nu}
				}
			}
		}(s, src)
	})

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		res.selection.SetAttr("src", res.newURL)
	}

	return doc.Html()
}

func (wp *WPClient) Publish(article *models.Article, catID int) (string, error) {
	utils.BroadcastLog("[PUBLISHER] Ուղարկվում է WordPress -> %s", article.Title)
	content, _ := wp.ProcessInlineImages(article.RewrittenContent.String)
	
	featID := 0
	if article.ImageURL.Valid {
		featID, _ = wp.UploadImageFromURL(article.ImageURL.String, article.ImageAlt.String)
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"title":          article.Title,
		"content":        content,
		"status":         "publish",
		"categories":     []int{catID},
		"excerpt":        article.MetaDescription.String,
		"featured_media": featID,
		"slug":           article.Slug.String,
	})

	req, _ := http.NewRequest("POST", wp.BaseURL+"/wp-json/wp/v2/posts", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := wp.doRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Link string `json:"link"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Link, nil
}
