package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"thepress_bot_go/internal/domain/services"
	"thepress_bot_go/internal/infra/utils"
	"time"
)

type ModelsLabProvider struct {
	APIKey string
	Client *http.Client
}

func NewModelsLabProvider(apiKey string) *ModelsLabProvider {
	return &ModelsLabProvider{
		APIKey: apiKey,
		Client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (m *ModelsLabProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error) {
	return nil, fmt.Errorf("ModelsLab not used for text processing")
}

func (m *ModelsLabProvider) ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error) {
	utils.BroadcastLog("[MODELSLAB] Starting image cleaning (Object Removal)...")

	// ModelsLab often requires a URL for the image or a direct upload.
	// For their "Magic Erase" or "Inpainting" API:
	url := "https://modelslab.com/api/v3/magic_erase" // Example endpoint for 2026

	// Note: Actual ModelsLab API might require multipart upload for raw bytes.
	// This is a standard implementation for their v3 API.
	bodyData := map[string]interface{}{
		"key":             m.APIKey,
		"image":           fmt.Sprintf("data:image/jpeg;base64,%s", utils.ToBase64(imageBytes)),
		"prompt":          prompt,
		"negative_prompt": "logos, text, watermarks, distorted, blurry",
		"steps":           20,
		"guidance_scale":  7.5,
	}

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modelslab request: %w", err)
	}
	resp, err := m.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("ModelsLab error %d: could not read body: %v", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("ModelsLab error %d: %s", resp.StatusCode, string(b))
	}

	var res struct {
		Status string   `json:"status"`
		Output []string `json:"output"` // They usually return a URL to the result
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if res.Status != "success" || len(res.Output) == 0 {
		return nil, fmt.Errorf("ModelsLab failed or still processing: %s", res.Status)
	}

	// Download the result from the URL provided by ModelsLab
	utils.BroadcastLog("[MODELSLAB] Image processed, downloading result...")
	return utils.DownloadFileToBytes(res.Output[0])
}
