package publisher

import (
	"context"
	"fmt"
	"thepress_bot_go/internal/domain/repository"
	"thepress_bot_go/internal/config"
	"strings"
)

type InternalLinker struct {
	repo repository.ArticleRepository
}

func NewInternalLinker(repo repository.ArticleRepository) *InternalLinker {
	return &InternalLinker{repo: repo}
}

func (l *InternalLinker) InjectLinks(ctx context.Context, articleID int64, htmlContent string) string {
	related, err := l.repo.GetRelated(ctx, articleID, 3)
	if err != nil || len(related) == 0 {
		return htmlContent
	}

	cfg := config.Get()
	baseURL := strings.TrimSuffix(cfg.WordPress.URL, "/")

	linksHTML := "\n\n<div class=\"pro-related-section\" style=\"margin-top: 20px; padding: 15px; background: #f9f9f9; border-left: 4px solid #3b82f6;\">"
	linksHTML += "<h4 style=\"margin-top: 0;\">Կարդացեք նաև՝</h4><ul style=\"list-style: none; padding: 0;\">"

	for _, art := range related {
		// Generate the local URL instead of the source URL
		localURL := fmt.Sprintf("%s/%s", baseURL, art.Slug.String)
		linksHTML += fmt.Sprintf("<li style=\"margin-bottom: 8px;\"><a href=\"%s\" style=\"color: #2563eb; text-decoration: none; font-weight: 500;\">Կարդալ՝ %s</a></li>", localURL, art.Title)
	}

	linksHTML += "</ul></div>"

	return htmlContent + linksHTML
}
