package classifier

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Analyze verilen HTML içeriğini ve URLi analiz ederek en uygun kategoriyi belirler.
func Analyze(htmlContent string, url string, linkCount int) Result {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		// HTML bozuksa bile en azından URL analizi yap
		return simpleAnalyze()
	}

	bestScore := 0
	var bestCategory *Category

	// Tüm kategorileri gez ve puanla
	for i := range GlobalConfig.Categories {
		cat := &GlobalConfig.Categories[i]
		score := calculateScore(cat, doc, htmlContent, linkCount)

		// Eğer öncelikli bir yapısal eşleşme varsa (Priority >= 100) ve skor > 0 ise direkt etiketlenir
		if cat.Priority >= 100 && score > 0 {
			return Result{
				CategoryID: cat.ID,
				Tag:        cat.Tag,
				Color:      cat.Color,
				Score:      score,
				IsUnknown:  false,
			}
		}

		if score > bestScore {
			bestScore = score
			bestCategory = cat
		}
	}

	// Hiçbir kategori eşiği geçemediyse (Eşik = 10 puan - En az 1 güçlü eşleşme lazım)
	if bestScore < 10 || bestCategory == nil {
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

// AnalyzeLinkContext henüz gidilmemiş bir linki (URL ve text) analiz eder
func AnalyzeLinkContext(url, anchorText string) Result {
	bestScore := 0
	var bestCategory *Category

	textLower := strings.ToLower(anchorText)
	urlLower := strings.ToLower(url)

	for i := range GlobalConfig.Categories {
		cat := &GlobalConfig.Categories[i]
		score := 0

		// URL de geçiyor mu (örn: forum)
		if strings.Contains(urlLower, cat.ID) {
			score += 5
		}

		// Anchor Text analizi
		for _, kw := range cat.Keywords.High {
			if strings.Contains(textLower, strings.ToLower(kw)) {
				score += 5
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
			Tag:        bestCategory.Tag, // örn: [MARKET]
			Color:      bestCategory.Color,
			Score:      bestScore,
			IsUnknown:  false,
		}
	}

	// Bilinmeyen durum
	return Result{
		CategoryID: "unknown",
		Tag:        "[?]",
		Color:      "gray",
		Score:      0,
		IsUnknown:  true,
	}
}

func calculateScore(cat *Category, doc *goquery.Document, rawHTML string, linkCount int) int {
	score := 0
	lowerHTML := strings.ToLower(rawHTML)

	// Max link kontrolü (login sayfaları için)
	// Eğer kategori max_links sınırı koymuşsa ve link sayısı bunu aşıyorsa, bu kategori olamaz (ceza puanı: -15)
	if cat.MaxLinks > 0 && linkCount > cat.MaxLinks {
		return -15
	}

	// Yapısal Analiz (Selector)
	for _, rule := range cat.StructureRules {
		if doc.Find(rule.Selector).Length() > 0 {
			score += 50 // Yapısal eşleşme yüksek puan
		}
	}

	// Kelime Analizi (Yüksek)
	for _, kw := range cat.Keywords.High {
		if strings.Contains(lowerHTML, strings.ToLower(kw)) {
			score += 10
		}
	}

	// 4. Kelime Analizi (Orta)
	for _, kw := range cat.Keywords.Medium {
		if strings.Contains(lowerHTML, strings.ToLower(kw)) {
			score += 5
		}
	}

	// 5. Negatif Kelimeler (Hariç Tut)
	for _, kw := range cat.Keywords.Exclude {
		if strings.Contains(lowerHTML, strings.ToLower(kw)) {
			score -= 50
		}
	}

	return score
}

// simpleAnalyze goquery çökerse fallback olarak çalışır
func simpleAnalyze() Result {
	return Result{
		CategoryID: "unknown",
		Tag:        "[BİLİNMEYEN]",
		Color:      "gray",
		IsUnknown:  true,
	}
}
