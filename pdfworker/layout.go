package main

// KeyPos define solo la posición y tamaño de una tecla dentro del layout físico.
// Esto NO cambia según el diseño del usuario; es fijo por tipo de teclado.
type KeyPos struct {
	X, Y, Width, Height float64
}

// testLayout: 4 teclas en fila, layout de prueba (placeholder hasta tener
// el layout real armado en Figma/HTML). Coordenadas en mm.
var testLayout = []KeyPos{
	{X: 5, Y: 5, Width: 18, Height: 18},
	{X: 26, Y: 5, Width: 18, Height: 18},
	{X: 47, Y: 5, Width: 18, Height: 18},
	{X: 68, Y: 5, Width: 18, Height: 18},
}

const (
	testSheetWidth  = 100.0
	testSheetHeight = 60.0
)