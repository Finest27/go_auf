package services

import (
	"context"
)

type AIResult struct {
	RewrittenContent string `json:"rewritten"`
	MetaDescription  string `json:"description"`
	FocusKeywords    string `json:"keywords"`
	Slug             string `json:"slug"`
	Category         string `json:"category"`
	Tags             string `json:"tags"`
	ImageAlt         string `json:"image_alt"`
}

type AIProvider interface {
	ProcessArticle(ctx context.Context, title, content string) (*AIResult, error)
	ProcessImage(ctx context.Context, imageBytes []byte, prompt string) ([]byte, error)
}
