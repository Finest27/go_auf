package ai

import (
	"context"
	"thepress_bot_go/internal/domain/services"
)

func NewProvider(ctx context.Context, tool, apiKey, sysPrompt string) (services.AIProvider, error) {
	return NewNvidiaProvider(apiKey, sysPrompt), nil
}
