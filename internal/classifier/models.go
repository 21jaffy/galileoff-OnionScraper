package classifier

// ClassificationConfig YAML dosyasının ana yapısı
type ClassificationConfig struct {
	Categories []Category `yaml:"categories"`
}

// Category tek bir kategori tanımı (Market, Forum, Login vb.)
type Category struct {
	ID             string          `yaml:"id"`
	Name           string          `yaml:"name"`
	Tag            string          `yaml:"tag"`
	Color          string          `yaml:"color"` // red, green, blue, yellow, magenta, cyan, white
	Priority       int             `yaml:"priority"`
	Keywords       KeywordRules    `yaml:"keywords"`
	StructureRules []StructureRule `yaml:"structure_rules"`
	MaxLinks       int             `yaml:"max_links"` // Opsiyonel: belli bir link sayısından fazlaysa bu kategori olamaz
}

// KeywordRules kelime bazlı kurallar
type KeywordRules struct {
	High    []string `yaml:"high"`    // Yüksek puanlı (10 puan)
	Medium  []string `yaml:"medium"`  // Orta puanlı (5 puan)
	Exclude []string `yaml:"exclude"` // Negatif puanlı (-50 puan) - Bu kelime varsa kategori puanı düşer
}

// StructureRule yapısal analiz kuralları (DOM)
type StructureRule struct {
	Selector string `yaml:"selector"` // CSS Selector (örn: input[type='password'])
}

// Result analiz sonucu döndürülen yapı
type Result struct {
	CategoryID string
	Tag        string
	Color      string
	Score      int
	IsUnknown  bool
}
