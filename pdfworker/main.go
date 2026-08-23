package main

import (
	"context"
	"fmt"
	"image/color"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/pdf"
)

// Key es la tecla ya combinada: posición (del template) + carácter primario
// + carácter secundario opcional + colores (del key_mode elegido).
type Key struct {
	PrimaryChar   string
	SecondaryChar string // "" si no hay alfabeto secundario
	X             float64
	Y             float64
	Width         float64
	Height        float64
	Bg            string
	FgPrimary     string
	FgSecondary   string
}

// Design es lo que efectivamente se dibuja: hoja + lista de teclas ya resueltas.
type Design struct {
	SheetWidth        float64
	SheetHeight       float64
	Keys              []Key
	PrimaryAlphabet   string
	SecondaryAlphabet string // "" si no hay secundario
}

// buildDesign combina el template de layout físico con las teclas de los
// alfabetos elegidos (primario y opcionalmente secundario) y el key_mode.
func buildDesign(ctx context.Context, conn *pgx.Conn, cfg *DesignConfig) (*Design, error) {
	primaryKeys, err := fetchKeyboardLayout(ctx, conn, cfg.PrimaryAlphabet)
	if err != nil {
		return nil, err
	}

	var secondaryKeys []LayoutKey
	if cfg.SecondaryAlphabet != nil {
		secondaryKeys, err = fetchKeyboardLayout(ctx, conn, *cfg.SecondaryAlphabet)
		if err != nil {
			return nil, err
		}
	}

	palette, err := paletteForKeyMode(cfg.KeyMode)
	if err != nil {
		return nil, err
	}

	if len(primaryKeys) < len(testLayout) {
		return nil, fmt.Errorf("el alfabeto %q tiene menos teclas (%d) que el layout (%d)", cfg.PrimaryAlphabet, len(primaryKeys), len(testLayout))
	}

	keys := make([]Key, len(testLayout))
	for i, pos := range testLayout {
		k := Key{
			PrimaryChar: primaryKeys[i].Base,
			X:           pos.X,
			Y:           pos.Y,
			Width:       pos.Width,
			Height:      pos.Height,
			Bg:          palette.Bg,
			FgPrimary:   palette.FgPrimary,
			FgSecondary: palette.FgSecondary,
		}
		if secondaryKeys != nil && i < len(secondaryKeys) {
			k.SecondaryChar = secondaryKeys[i].Base
		}
		keys[i] = k
	}

	secondaryAlphabet := ""
	if cfg.SecondaryAlphabet != nil {
		secondaryAlphabet = *cfg.SecondaryAlphabet
	}

	return &Design{
		SheetWidth:        testSheetWidth,
		SheetHeight:       testSheetHeight,
		Keys:              keys,
		PrimaryAlphabet:   cfg.PrimaryAlphabet,
		SecondaryAlphabet: secondaryAlphabet,
	}, nil
}

func hexColor(hex string) color.RGBA {
	return canvas.Hex(hex)
}

func renderDesign(d *Design, outPath string) error {
	c := canvas.New(d.SheetWidth, d.SheetHeight)
	ctx := canvas.NewContext(c)

	primaryFamily := canvas.NewFontFamily("primary")
	if err := primaryFamily.LoadFontFile(fontPathFor(d.PrimaryAlphabet), canvas.FontRegular); err != nil {
		return fmt.Errorf("loading primary font (%s): %w", d.PrimaryAlphabet, err)
	}

	var secondaryFamily *canvas.FontFamily
	if d.SecondaryAlphabet != "" {
		secondaryFamily = canvas.NewFontFamily("secondary")
		if err := secondaryFamily.LoadFontFile(fontPathFor(d.SecondaryAlphabet), canvas.FontRegular); err != nil {
			return fmt.Errorf("loading secondary font (%s): %w", d.SecondaryAlphabet, err)
		}
	}

	for _, k := range d.Keys {
		// fondo de la tecla (rectángulo)
		ctx.SetFillColor(hexColor(k.Bg))
		rect := canvas.Rectangle(k.Width, k.Height)
		// canvas usa coordenadas cartesianas (Y crece hacia arriba);
		// convertimos desde Y "hacia abajo" tipo pantalla/imprenta.
		posY := d.SheetHeight - k.Y - k.Height
		ctx.DrawPath(k.X, posY, rect)

		// carácter primario, grande, centrado
		facePrimary := primaryFamily.Face(10.0, hexColor(k.FgPrimary), canvas.FontRegular, canvas.FontNormal)
		primaryText := canvas.NewTextLine(facePrimary, k.PrimaryChar, canvas.Center)
		ctx.DrawText(k.X+k.Width/2, posY+k.Height/2-3.5, primaryText)

		// carácter secundario, chico, esquina superior derecha (si existe)
		if k.SecondaryChar != "" && secondaryFamily != nil {
			faceSecondary := secondaryFamily.Face(5.0, hexColor(k.FgSecondary), canvas.FontRegular, canvas.FontNormal)
			secondaryText := canvas.NewTextLine(faceSecondary, k.SecondaryChar, canvas.Right)
			ctx.DrawText(k.X+k.Width-2, posY+k.Height-4, secondaryText)
		}
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	pdfRenderer := pdf.New(f, c.W, c.H, nil)
	c.RenderTo(pdfRenderer)
	return pdfRenderer.Close()
}

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		// si corrés esto desde la raíz del proyecto en vez de /pdfworker,
		// probá también ".env" sin el "../"
		if err2 := godotenv.Load(".env"); err2 != nil {
			panic("no se encontró .env: " + err.Error())
		}
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		panic("conectando a postgres: " + err.Error())
	}
	defer conn.Close(ctx)

	designID := 1
	if len(os.Args) > 1 {
		id, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic("id de design inválido: " + os.Args[1])
		}
		designID = id
	}

	cfg, err := fetchDesign(ctx, conn, designID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("design %d: primary=%s secondary=%v key_mode=%s\n", designID, cfg.PrimaryAlphabet, cfg.SecondaryAlphabet, cfg.KeyMode)

	design, err := buildDesign(ctx, conn, cfg)
	if err != nil {
		panic(err)
	}

	outPath := fmt.Sprintf("design_%d.pdf", designID)
	if err := renderDesign(design, outPath); err != nil {
		panic(err)
	}

	info, _ := os.Stat(outPath)
	fmt.Printf("PDF generado: %s (%d bytes)\n", outPath, info.Size())
}