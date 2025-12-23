package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// UserAgentProfile User-Agent ve ilgili başlıklarla birlikte tarayıcı profili tanımlar
type UserAgentProfile struct {
	Name      string            `json:"name"`
	UserAgent string            `json:"user_agent"`
	Headers   map[string]string `json:"headers"`
}

// profiles Varsayılan UA için liste
var profiles = []UserAgentProfile{
	{
		Name:      "Default (Chrome 120)",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		Headers: map[string]string{
			"Accept":             "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
			"Accept-Language":    "en-US,en;q=0.9",
			"Sec-Ch-Ua":          `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": `"Windows"`,
		},
	},
}

// init rastgele sayı üreticisini başlatır (.jsondaki UA için)
func init() {
	rand.Seed(time.Now().UnixNano())
}

// LoadProfiles .json dosyasından User-Agent yükler
func LoadProfiles(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var loadedProfiles []UserAgentProfile
	if err := json.Unmarshal(data, &loadedProfiles); err != nil {
		return fmt.Errorf("JSON parse hatası: %v", err)
	}

	if len(loadedProfiles) > 0 {
		profiles = loadedProfiles
	}
	return nil
}

// GetRandomProfile rastgele bir UserAgentProfile döndürür
func GetRandomProfile() UserAgentProfile {
	if len(profiles) == 0 {
		return UserAgentProfile{
			Name:      "Fallback",
			UserAgent: "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			Headers:   map[string]string{},
		}
	}
	return profiles[rand.Intn(len(profiles))]
}
