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

	re := regexp.MustCompile(`\n{3,}`)
	clean = re.ReplaceAllString(clean, "\n\n")

	return strings.TrimSpace(clean)
}

func CleanJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	
	// Find the first { and the last }
	startIdx := strings.Index(raw, "{")
	endIdx := strings.LastIndex(raw, "}")
	
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		return raw[startIdx : endIdx+1]
	}

	// Fallback to basic cleaning
	if strings.HasPrefix(raw, "```json") {
		raw = strings.TrimPrefix(raw, "```json")
	} else if strings.HasPrefix(raw, "```") {
		raw = strings.TrimPrefix(raw, "```")
	}
	raw = strings.TrimSuffix(raw, "```")
	return strings.TrimSpace(raw)
}
