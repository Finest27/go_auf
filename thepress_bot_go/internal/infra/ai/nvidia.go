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
		Client:     &http.Client{Timeout: 300 * time.Second},
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
	reqBody, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
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
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if len(result.Artifacts) == 0 { return nil, fmt.Errorf("no image returned") }
	return base64.StdEncoding.DecodeString(result.Artifacts[0].Base64)
}

func (n *NvidiaProvider) ProcessArticle(ctx context.Context, title, content string) (*services.AIResult, error) {
	url := "https://integrate.api.nvidia.com/v1/chat/completions"

	utils.BroadcastLog("[AI] Starting Multi-Stage Rewriting Pipeline for: %s", title)

	// Multi-Stage Rewriting Pipeline
	// Stage 1: Extraction & Summarization
	stage1Prompt := fmt.Sprintf("Extract the core facts, entities, and primary narrative from this article.\nTitle: %s\nContent: %s", title, content)
	facts, err := n.runSimpleCompletion(ctx, url, n.TextModels[0], stage1Prompt)
	if err != nil {
		utils.BroadcastLog("[AI WARNING] Stage 1 (Extraction) failed, falling back to direct rewrite.")
		facts = content // fallback
	}

	// Stage 2: Localization & Stylistic Adaptation (Armenian) and JSON output
	userPrompt := fmt.Sprintf("Based on these facts, write a high-quality journalistic article in Armenian.\nFacts:\n%s", facts)
	
	var res *services.AIResult
	var lastErr error
	
	// Circuit Breaker / Fallback Logic
	for _, model := range n.TextModels {
		res, lastErr = n.tryModel(ctx, url, model, userPrompt)
		if lastErr == nil && res != nil { 
			utils.BroadcastLog("[AI SUCCESS] Article rewritten successfully using %s", model)
			return res, nil 
		}
		
		// SCHEME 2: AI Self-Correction if JSON is invalid
		if lastErr != nil && (res == nil) {
			utils.BroadcastLog("[AI] JSON error with %s, attempting self-correction...", model)
			retryPrompt := userPrompt + "\n\nCRITICAL: Your previous response was not a valid JSON. Output ONLY the JSON object, no markdown, no conversational text. Include fields: rewritten, description, keywords, slug, category, tags, image_alt."
			res, lastErr = n.tryModel(ctx, url, model, retryPrompt)
			if lastErr == nil && res != nil { 
				utils.BroadcastLog("[AI SUCCESS] Article rewritten successfully after correction using %s", model)
				return res, nil 
			}
		}
		
		utils.BroadcastLog("[AI] Model %s failed: %v", model, lastErr)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("circuit breaker triggered: all models failed. Last error: %w", lastErr)
}

func (n *NvidiaProvider) runSimpleCompletion(ctx context.Context, url, model, prompt string) (string, error) {
	bodyData := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are an expert journalist. Extract facts objectively."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.1,
	}
	return n.executeRequest(ctx, url, bodyData)
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
	bodyBytes, err := json.Marshal(bodyData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
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
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	if len(result.Choices) == 0 { return "", fmt.Errorf("no response") }
	return result.Choices[0].Message.Content, nil
}

func (n *NvidiaProvider) Close() {}
