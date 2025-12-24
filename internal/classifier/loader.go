package classifier

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var GlobalConfig ClassificationConfig

// LoadRules belirtilen dosya yolundan kuralları yükler
func LoadRules(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("kural dosyası okunamadı: %v", err)
	}

	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		return fmt.Errorf("YAML parse hatası: %v", err)
	}

	return nil
}

// GetCategoryByID ID'ye göre kategori bilgisini döndürür
func GetCategoryByID(id string) *Category {
	for i := range GlobalConfig.Categories {
		if GlobalConfig.Categories[i].ID == id {
			return &GlobalConfig.Categories[i]
		}
	}
	return nil
}
