package classifier

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Analyze HTML + URL analiz eder
func Analyze(htmlContent string, url string, linkCount int) Result {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return simpleAnalyze()
	}

	// İçerik ve Kod Ayrımı
	visibleText := extractVisibleText(doc)

	// Gelişmiş Meta Analizi
	title := doc.Find("title").Text()
	metaDesc, _ := doc.Find("meta[name='description']").Attr("content")
	metaKw, _ := doc.Find("meta[name='keywords']").Attr("content")
	ogTitle, _ := doc.Find("meta[property='og:title']").Attr("content")
	ogDesc, _ := doc.Find("meta[property='og:description']").Attr("content")

	// Meta verileri birleştir
	combinedMeta := strings.Join([]string{title, metaDesc, metaKw, ogTitle, ogDesc}, " ")

	bestScore := 0
	var bestCategory *Category

	for i := range GlobalConfig.Categories {
		cat := &GlobalConfig.Categories[i]
		score := calculateScore(cat, doc, htmlContent, visibleText, combinedMeta, linkCount)

		if score > bestScore {
			bestScore = score
			bestCategory = cat
		}
	}

	// Login override: başka ciddi kategori varsa login ezilmesin
	if bestCategory != nil && bestCategory.ID == "login" {
		for i := range GlobalConfig.Categories {
			cat := &GlobalConfig.Categories[i]
			if cat.ID == "login" {
				continue
			}
			altScore := calculateScore(cat, doc, htmlContent, visibleText, combinedMeta, linkCount)
			if altScore >= bestScore-10 {
				bestScore = altScore
				bestCategory = cat
			}
		}
	}

	if bestScore < 20 || bestCategory == nil {
		return Result{
			CategoryID: "unknown",
			Tag:        "[BİLİNMEYEN]",
			Color:      "gray",
			Score:      0,
			IsUnknown:  true,
		}
	}

	return Result{
		CategoryID: bestCategory.ID,
		Tag:        bestCategory.Tag,
		Color:      bestCategory.Color,
		Score:      bestScore,
		IsUnknown:  false,
	}
}

// AnalyzeLinkContext henüz girilmemiş linkleri analiz eder
func AnalyzeLinkContext(url, anchorText string) Result {
	bestScore := 0
	var bestCategory *Category

	textLower := strings.ToLower(anchorText)
	urlLower := strings.ToLower(url)

	for i := range GlobalConfig.Categories {
		cat := &GlobalConfig.Categories[i]
		score := 0

		if strings.Contains(urlLower, cat.ID) {
			score += 5
		}

		for _, kw := range cat.Keywords.High {
			if strings.Contains(textLower, strings.ToLower(kw)) {
				score += 5
			}
		}

		for _, kw := range cat.Keywords.Medium {
			if strings.Contains(textLower, strings.ToLower(kw)) {
				score += 2
			}
		}

		for _, kw := range cat.Keywords.Exclude {
			if strings.Contains(textLower, strings.ToLower(kw)) {
				score -= 10
			}
		}

		if score > bestScore {
			bestScore = score
			bestCategory = cat
		}
	}

	if bestScore >= 5 && bestCategory != nil {
		return Result{
			CategoryID: bestCategory.ID,
			Tag:        bestCategory.Tag,
			Color:      bestCategory.Color,
			Score:      bestScore,
			IsUnknown:  false,
		}
	}

	return Result{
		CategoryID: "unknown",
		Tag:        "[?]",
		Color:      "gray",
		Score:      0,
		IsUnknown:  true,
	}
}

// calculateScore kategori skorunu hesaplar
func calculateScore(cat *Category, doc *goquery.Document, rawHTML, visibleText, metaText string, linkCount int) int {
	score := 0
	lowerHTML := strings.ToLower(rawHTML)
	lowerVisible := strings.ToLower(visibleText)
	lowerMeta := strings.ToLower(metaText)

	// Max link kontrolü (ceza)
	if cat.MaxLinks > 0 && linkCount > cat.MaxLinks {
		score -= 15
	}

	// Yapısal analiz
	for _, rule := range cat.StructureRules {
		if doc.Find(rule.Selector).Length() > 0 {
			score += 20
		}
	}

	highHit := 0
	for _, kw := range cat.Keywords.High {
		k := strings.ToLower(kw)
		matched := false

		// Görünen metinde varsa +15
		if strings.Contains(lowerVisible, k) {
			score += 15
			matched = true
		}
		// Meta verilerde varsa +15
		if strings.Contains(lowerMeta, k) {
			score += 15
			matched = true
		}
		// Sadece HTML içinde varsa ama yukarıdakilerde yoksa +5
		if !matched && strings.Contains(lowerHTML, k) {
			score += 5
		}

		if matched {
			highHit++
		}
		if highHit > 5 {
		}
	}

	// Medium keyword
	medHit := 0
	for _, kw := range cat.Keywords.Medium {
		k := strings.ToLower(kw)
		if strings.Contains(lowerVisible, k) || strings.Contains(lowerMeta, k) {
			medHit++
			if medHit <= 7 {
				score += 5
			}
		} else if strings.Contains(lowerHTML, k) {
			// Sadece kod içinde varsa düşük puan
			score += 2
		}
	}

	// Exclude kelimeler
	for _, kw := range cat.Keywords.Exclude {
		k := strings.ToLower(kw)
		if strings.Contains(lowerVisible, k) || strings.Contains(lowerMeta, k) {
			score -= 50
		} else if strings.Contains(lowerHTML, k) {
			score -= 20
		}
	}

	return score
}

// extractVisibleText sadece sayfanın görünen metnini çeker
func extractVisibleText(doc *goquery.Document) string {
	// Bodynin bir kopyasını al
	selection := doc.Find("body").Clone()

	// İstenmeyen tagleri temizle
	selection.Find("script, style, noscript, iframe, svg").Remove()

	// Metni al ve boşlukları temizle
	text := selection.Text()
	return strings.Join(strings.Fields(text), " ")
}

// fallback
func simpleAnalyze() Result {
	return Result{
		CategoryID: "unknown",
		Tag:        "[BİLİNMEYEN]",
		Color:      "gray",
		IsUnknown:  true,
	}
}
