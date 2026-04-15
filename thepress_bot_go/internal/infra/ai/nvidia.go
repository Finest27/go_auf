package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	
	"thepress_bot_go/internal/domain/services"
	"thepress_bot_go/internal/infra/utils"
	"time"
)

type NvidiaProvider struct {
	APIKey       string
	SystemPrompt string
	TextModels   []string
	ImageModel   string
	Client       *http.Client
}

func NewNvidiaProvider(apiKey, sysPrompt string) *NvidiaProvider {
	return &NvidiaProvider{
		APIKey:       apiKey,
		SystemPrompt: sysPrompt,
		TextModels: []string{
			"meta/llama-3.1-70b-instruct",
			"meta/llama-3.1-8b-instruct",
			"nvidia/nemotron-4-340b-instruct",
		},
		ImageModel: "stabilityai/sdxl-turbo",
		Client:     &http.Client{Timeout: 120 * time.Second},
	}
}

func (n *NvidiaProvider) ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error) {
	encoded := base64.StdEncoding.EncodeToString(imageBytes)
	url := fmt.Sprintf("https://ai.api.nvidia.com/v1/genai/%s", n.ImageModel)
	bodyData := map[string]interface{}{
		"text_prompts":   []map[string]interface{}{{"text": prompt, "weight": 1.0}},
		"init_image":     encoded,
		"image_strength": 0.35,
		"cfg_scale":      7,
		"sampler":        "K_EULER_ANCESTRAL",
		"steps":          2,
	}
	reqBody, _ := json.Marshal(bodyData)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+n.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := n.Client.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	
	var result struct {
		Artifacts []struct {
			Base64 string `json:"base64"`
		} `json:"artifacts"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if len(result.Artifacts) == 0 { return nil, fmt.Errorf("no image returned") }
	return base64.StdEncoding.DecodeString(result.Artifacts[0].Base64)
}

func (n *NvidiaProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error) {
	url := "https://integrate.api.nvidia.com/v1/chat/completions"
	userPrompt := fmt.Sprintf("????????: %s\n???????????????: %s", title, content)

	for _, model := range n.TextModels {
		res, err := n.tryModel(ctx, url, model, userPrompt)
		if err == nil { return res, nil }
		
		// SCHEME 2: AI Self-Correction if JSON is invalid
		if err != nil && (res == nil) {
			utils.BroadcastLog("[AI] JSON error with %s, attempting self-correction...", model)
			retryPrompt := userPrompt + "\n\nCRITICAL: Your previous response was not a valid JSON. Output ONLY the JSON object, no markdown, no conversational text."
			res, err = n.tryModel(ctx, url, model, retryPrompt)
			if err == nil { return res, nil }
		}
		
		utils.BroadcastLog("[AI] %s failed: %v", model, err)
		time.Sleep(1 * time.Second)
	}
	return nil, fmt.Errorf("all models failed")
}

func (n *NvidiaProvider) tryModel(ctx context.Context, url, model, prompt string) (*services.AIResult, error) {
	bodyData := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": n.SystemPrompt},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
		"response_format": map[string]string{"type": "json_object"},
	}
	
	raw, err := n.executeRequest(ctx, url, bodyData)
	if err != nil { return nil, err }
	
	var aiRes services.AIResult
	clean := utils.CleanJSON(raw)
	if err := json.Unmarshal([]byte(clean), &aiRes); err != nil {
		return nil, err
	}
	return &aiRes, nil
}

func (n *NvidiaProvider) executeRequest(ctx context.Context, url string, bodyData interface{}) (string, error) {
	bodyBytes, _ := json.Marshal(bodyData)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+n.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := n.Client.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error %d", resp.StatusCode)
	}
	
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if len(result.Choices) == 0 { return "", fmt.Errorf("no response") }
	return result.Choices[0].Message.Content, nil
}

func (n *NvidiaProvider) Close() {}
