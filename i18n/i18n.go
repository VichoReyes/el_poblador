package i18n

import (
	"embed"
	"os"

	"github.com/leonelquinteros/gotext"
)

//go:embed locales
var localesFS embed.FS

var locale *gotext.Locale

func init() {
	// Determine language from environment, default to English
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "en_US"
	}

	// Extract language code (e.g., "es" from "es_ES.UTF-8")
	if len(lang) >= 2 {
		lang = lang[:2]
	}

	// Initialize locale with embedded filesystem
	locale = gotext.NewLocale("", lang)
	locale.AddDomain("el_poblador")

	// Load translations from embedded filesystem
	if lang != "en" {
		// Try to load .po file from embedded filesystem
		poFile := "locales/" + lang + "/LC_MESSAGES/el_poblador.po"
		data, err := localesFS.ReadFile(poFile)
		if err == nil {
			locale.ParsePO(data)
		}
	}
}

// T translates a string using gettext
func T(msgid string) string {
	return locale.Get(msgid)
}
