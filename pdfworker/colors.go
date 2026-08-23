package main

import "fmt"

// ColorSet: colores para una tecla según su key_mode.
// Bg = fondo, FgPrimary = color del carácter principal (grande, centrado),
// FgSecondary = color del carácter secundario (chico, esquina).
type ColorSet struct {
	Bg          string
	FgPrimary   string
	FgSecondary string
}

var keyModePalette = map[string]ColorSet{
	"black": {Bg: "#1d1d1f", FgPrimary: "#ffffff", FgSecondary: "#39ff14"},
	"white": {Bg: "#f5f5f7", FgPrimary: "#1d1d1f", FgSecondary: "#16a34a"},
	"gray":  {Bg: "#8e8e93", FgPrimary: "#1d1d1f", FgSecondary: "#39ff14"},
}

func paletteForKeyMode(mode string) (ColorSet, error) {
	c, ok := keyModePalette[mode]
	if !ok {
		return ColorSet{}, fmt.Errorf("key_mode desconocido: %q", mode)
	}
	return c, nil
}