package main

import "fmt"

// ColorSet: colores para una tecla según su key_mode.
// Bg = fondo, FgPrimary = carácter principal, FgCorner = shift+altgr
// (unificados), FgSecondary = color por defecto del alfabeto secundario
// (se usa solo si el diseño no trae su propio secondary_color).
type ColorSet struct {
	Bg          string
	FgPrimary   string
	FgCorner    string
	FgSecondary string
}

var keyModePalette = map[string]ColorSet{
	"black": {Bg: "#1f1f1f", FgPrimary: "#f0f0f0", FgCorner: "#999999", FgSecondary: "#e0995e"},
	"white": {Bg: "#fdfcf9", FgPrimary: "#1d1d1f", FgCorner: "#999999", FgSecondary: "#2f6fed"},
	"gray":  {Bg: "#d6d6d8", FgPrimary: "#2a2a2a", FgCorner: "#6b6b6b", FgSecondary: "#2f9e44"},
}

func paletteForKeyMode(mode string) (ColorSet, error) {
	c, ok := keyModePalette[mode]
	if !ok {
		return ColorSet{}, fmt.Errorf("key_mode desconocido: %q", mode)
	}
	return c, nil
}