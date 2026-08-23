package main

// KeyPos define solo la posición y tamaño de una tecla dentro del layout físico.
// Esto NO cambia según el diseño del usuario; es fijo por tipo de teclado.
type KeyPos struct {
	X, Y, Width, Height float64
}

// testLayoutGrid genera una grilla simple de N teclas (placeholder hasta
// tener el layout real armado en Figma/HTML). Coordenadas en mm.
func testLayoutGrid(count int, cols int) []KeyPos {
	const keyW, keyH, gap = 12.0, 12.0, 2.0
	keys := make([]KeyPos, count)
	for i := 0; i < count; i++ {
		row := i / cols
		col := i % cols
		keys[i] = KeyPos{
			X:      5 + float64(col)*(keyW+gap),
			Y:      5 + float64(row)*(keyH+gap),
			Width:  keyW,
			Height: keyH,
		}
	}
	return keys
}

// testLayout: grilla de 47 teclas (12 columnas), suficiente para el
// AZERTY completo cargado en la DB. Placeholder de posición, no diseño final.
var testLayout = testLayoutGrid(47, 12)

var (
	testSheetWidth  = 5 + 12.0*12 + 2.0*11 + 5
	testSheetHeight = 5 + 12.0*4 + 2.0*3 + 5
)
