package main

import "os"
import "path/filepath"

var fontFiles = map[string]string{
	"hiragana":  "msgothic.ttc",
	"arabic":    "tahoma.ttf",
	"azerty_fr": "arialbd.ttf",
	"hangul":    "malgunbd.ttf",
	"cangjie":   "msyhbd.ttc",
}

func fontPathFor(alphabetName string) string {
	dir := os.Getenv("FONT_DIR")
	if dir == "" {
		dir = `C:\Windows\Fonts`
	}

	file, ok := fontFiles[alphabetName]
	if !ok {
		file = "arialbd.ttf"
	}

	return filepath.Join(dir, file)
}