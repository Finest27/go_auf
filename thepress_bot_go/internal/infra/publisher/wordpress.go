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
	"regexp"
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

	utils.BroadcastLog("[MEDIA PIPELINE] Downloading featured image from: %s", imageURL)
	var err error
	for i := 0; i < 2; i++ {
		if err = utils.DownloadFileToPath(imageURL, temp); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		utils.BroadcastLog("[MEDIA PIPELINE WARNING] Failed to download image: %v", err)
		return 0, err
	}
	defer os.Remove(temp)

	if wp.AI != nil {
		imgBytes, err := os.ReadFile(temp)
		if err == nil {
			utils.BroadcastLog("[MEDIA PIPELINE] AI processing image (inpainting/cleaning)...")
			processed, err := wp.AI.ProcessImage(context.Background(), imgBytes, "Clean news photo")
			if err == nil && len(processed) > 0 {
				utils.BroadcastLog("[MEDIA PIPELINE] AI processing successful, saving enhanced image.")
				_ = os.WriteFile(temp, processed, 0644)
			} else {
				utils.BroadcastLog("[MEDIA PIPELINE WARNING] AI image processing failed or returned empty: %v", err)
			}
		}
	}

	utils.BroadcastLog("[MEDIA PIPELINE] Uploading processed image to WordPress Media Library...")
	id, _, err := wp.UploadImageToMediaLibrary(temp, altText)
	if err == nil {
		utils.BroadcastLog("[MEDIA PIPELINE] Successfully uploaded image. Media ID: %d", id)
	}
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
	part, err := writer.CreateFormFile("file", filepath.Base(imagePath))
	if err != nil {
		return 0, "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return 0, "", err
	}
	writer.Close()

	req, err := http.NewRequest("POST", wp.BaseURL+"/wp-json/wp/v2/media", body)
	if err != nil {
		return 0, "", err
	}
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
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, "", err
	}

	if altText != "" {
		buf := &bytes.Buffer{}
		_ = json.NewEncoder(buf).Encode(map[string]string{"alt_text": altText})
		reqUpdate, err := http.NewRequest("POST", fmt.Sprintf("%s/wp-json/wp/v2/media/%d", wp.BaseURL, res.ID), buf)
		if err == nil {
			reqUpdate.Header.Set("Content-Type", "application/json")
			updateResp, err := wp.doRequest(reqUpdate)
			if err == nil {
				updateResp.Body.Close()
			}
		}
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
		src, exists := s.Attr("src")
		if !exists || src == "" || strings.Contains(src, wp.BaseURL) {
			return
		}
		wg.Add(1)
		go func(sel *goquery.Selection, u string) {
			defer wg.Done()
			utils.BroadcastLog("[MEDIA PIPELINE] Processing inline image: %s", u)
			temp := filepath.Join(os.TempDir(), fmt.Sprintf("bimg_%d.jpg", time.Now().UnixNano()))
			if err := utils.DownloadFileToPath(u, temp); err == nil {
				defer os.Remove(temp)
				if id, nu, err := wp.UploadImageToMediaLibrary(temp, ""); err == nil && id > 0 {
					utils.BroadcastLog("[MEDIA PIPELINE] Inline image uploaded successfully. New URL: %s", nu)
					results <- uploadResult{sel, nu}
				} else {
					utils.BroadcastLog("[MEDIA PIPELINE WARNING] Failed to upload inline image: %v", err)
				}
			} else {
				utils.BroadcastLog("[MEDIA PIPELINE WARNING] Failed to download inline image: %v", err)
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
        finalTitle := article.Title
        finalContent := article.RewrittenContent.String
        
        re := regexp.MustCompile("(?i)<h1>(.*?)</h1>")
        matches := re.FindStringSubmatch(finalContent)
        if len(matches) > 1 {
                finalTitle = matches[1]
                finalContent = re.ReplaceAllString(finalContent, "")
        }
        article.Title = finalTitle
        article.RewrittenContent.String = finalContent

	utils.BroadcastLog("[PUBLISHER] Ուղարկվում է WordPress -> %s", article.Title)
	content, err := wp.ProcessInlineImages(article.RewrittenContent.String)
	if err != nil {
		// Log but continue with original rewritten content if processing inline images fails
		utils.BroadcastLog("[PUBLISHER WARNING] Failed to process inline images: %v", err)
		content = article.RewrittenContent.String
	}
	
	featID := 0
	if article.ImageURL.Valid {
		id, err := wp.UploadImageFromURL(article.ImageURL.String, article.ImageAlt.String)
		if err == nil {
			featID = id
		} else {
			utils.BroadcastLog("[PUBLISHER WARNING] Failed to upload featured image: %v", err)
		}
	}

	payload, err := json.Marshal(map[string]interface{}{
		"title":          article.Title,
		"content":        content,
		"status":         "publish",
		"categories":     []int{catID},
		"excerpt":        article.MetaDescription.String,
		"featured_media": featID,
		"slug":           article.Slug.String,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", wp.BaseURL+"/wp-json/wp/v2/posts", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := wp.doRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("wordpress returned status %d: %s", resp.StatusCode, string(b))
	}

	var result struct {
		Link string `json:"link"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Link, nil
}

