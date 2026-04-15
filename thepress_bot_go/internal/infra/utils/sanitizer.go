package utils

import (
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func SanitizeHTML(html string) string {
	p := bluemonday.StrictPolicy()
	p.AllowElements("p", "h2", "h3", "br")

	clean := p.Sanitize(html)

	re := regexp.MustCompile(`\n{3,}`)
	clean = re.ReplaceAllString(clean, "\n\n")

	return strings.TrimSpace(clean)
}

func CleanJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "```json") {
		raw = strings.TrimPrefix(raw, "```json")
	} else if strings.HasPrefix(raw, "```") {
		raw = strings.TrimPrefix(raw, "```")
	}
	if strings.HasSuffix(raw, "```") {
		raw = strings.TrimSuffix(raw, "```")
	}
	return strings.TrimSpace(raw)
}