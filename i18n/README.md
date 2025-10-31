# Internationalization (i18n)

This package provides internationalization support for El Poblador using gettext.

## Usage

To use translations in your code:

```go
import "el_poblador/i18n"

func example() {
    translatedString := i18n.T("Build")
}
```

## Language Selection

The language is automatically detected from the `LANG` environment variable.

To test Spanish translations:
```bash
LANG=es go run main.go new Player1 Player2 Player3
```

To use English (default):
```bash
go run main.go new Player1 Player2 Player3
```

## Translation Files

Translation files are located in `locales/<lang>/LC_MESSAGES/el_poblador.po`

Currently supported languages:
- English (default, no translation file needed)
- Spanish (`locales/es/LC_MESSAGES/el_poblador.po`)

## Embedded Translations

All translation files are embedded in the binary using Go's `embed` directive, so no external files are needed for deployment.
