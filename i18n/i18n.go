package i18n

import (
	"embed"
	"os"
	"strings"

	"github.com/leonelquinteros/gotext"
)

//go:embed locales
var localeFS embed.FS

var locale *gotext.Locale

func init() {
	// Detect language from LANG environment variable
	lang := detectLanguage()

	// Create locale from embedded filesystem
	locale = gotext.NewLocale(localeFS, lang)
	locale.SetDomain("el_poblador")
}

// detectLanguage detects the language from the environment
func detectLanguage() string {
	// Check LANG environment variable (e.g., "es_ES.UTF-8" -> "es")
	langEnv := os.Getenv("LANG")
	if langEnv != "" {
		// Extract language code (first two characters before _ or .)
		parts := strings.SplitN(langEnv, "_", 2)
		if len(parts[0]) >= 2 {
			return parts[0][:2]
		}
	}
	return "en"
}

// T translates a string to the current language
func T(key string) string {
	return locale.Get(key)
}
