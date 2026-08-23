package main

// fontPathForAlphabet: qué archivo de fuente usar según el alfabeto.
// Windows ya trae estas instaladas, no hace falta descargar nada:
//   - msgothic.ttc: japonés (hiragana, katakana, kanji)
//   - tahoma.ttf: buen soporte de árabe (y latín también)
//   - arialbd.ttf: latín (default para alfabetos con teclado físico normal)
var fontPathForAlphabet = map[string]string{
	"hiragana":  `C:\Windows\Fonts\msgothic.ttc`,
	"arabic":    `C:\Windows\Fonts\tahoma.ttf`,
	"azerty_fr": `C:\Windows\Fonts\arialbd.ttf`,
}

const defaultFontPath = `C:\Windows\Fonts\arialbd.ttf`

func fontPathFor(alphabetName string) string {
	if p, ok := fontPathForAlphabet[alphabetName]; ok {
		return p
	}
	return defaultFontPath
}