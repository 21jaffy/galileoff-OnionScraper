package utils

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// LinkData linkin URLini ve görünen metnini (anchor text) tutar
type LinkData struct {
	URL  string
	Text string
}

// ExtractLinks HTML içeriğinden tüm linkleri ve metinlerini çeker
func ExtractLinks(htmlContent string) []LinkData {
	var links []LinkData
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return links
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			val := strings.TrimSpace(href)
			if val != "" && !strings.HasPrefix(val, "#") && !strings.HasPrefix(val, "javascript:") && !strings.HasPrefix(val, "mailto:") {
				text := strings.TrimSpace(s.Text())
				links = append(links, LinkData{URL: val, Text: text})
			}
		}
	})

	return links
}
