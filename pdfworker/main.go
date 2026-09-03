package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/pdf"
)

type Key struct {
	Base           string
	Shift          string
	AltGr          string
	SecondaryChar  string
	SecondaryShift string
	SecondaryAltGr string
	X              float64
	Y              float64
	Width          float64
	Height         float64
	Bg             string
	FgPrimary      string
	FgCorner       string
	FgSecondary    string
}

type Design struct {
	SheetWidth        float64
	SheetHeight       float64
	Keys              []Key
	PrimaryAlphabet   string
	SecondaryAlphabet string
}

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

	secondaryColor := palette.FgSecondary
	if cfg.SecondaryColor != nil && *cfg.SecondaryColor != "" {
		secondaryColor = *cfg.SecondaryColor
	}

	if len(primaryKeys) < len(testLayout) {
		return nil, fmt.Errorf("el alfabeto %q tiene menos teclas (%d) que el layout (%d)", cfg.PrimaryAlphabet, len(primaryKeys), len(testLayout))
	}

	keys := make([]Key, len(testLayout))
	for i, pos := range testLayout {
		lk := primaryKeys[i]
		k := Key{
			Base:        lk.Base,
			X:           pos.X,
			Y:           pos.Y,
			Width:       pos.Width,
			Height:      pos.Height,
			Bg:          palette.Bg,
			FgPrimary:   palette.FgPrimary,
			FgCorner:    palette.FgCorner,
			FgSecondary: secondaryColor,
		}
		if lk.Shift != nil {
			k.Shift = *lk.Shift
		}
		if lk.AltGr != nil {
			k.AltGr = *lk.AltGr
		}
		if secondaryKeys != nil && i < len(secondaryKeys) {
			sk := secondaryKeys[i]
			k.SecondaryChar = sk.Base
			if sk.Shift != nil {
				k.SecondaryShift = *sk.Shift
			}
			if sk.AltGr != nil {
				k.SecondaryAltGr = *sk.AltGr
			}
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
		ctx.SetFillColor(hexColor(k.Bg))
		rect := canvas.Rectangle(k.Width, k.Height)
		posY := d.SheetHeight - k.Y - k.Height
		ctx.DrawPath(k.X, posY, rect)

		faceBase := primaryFamily.Face(9.0, hexColor(k.FgPrimary), canvas.FontRegular, canvas.FontNormal)
		baseText := canvas.NewTextLine(faceBase, k.Base, canvas.Center)
		ctx.DrawText(k.X+k.Width/2, posY+k.Height/2-3, baseText)

		if k.Shift != "" {
			faceShift := primaryFamily.Face(4.5, hexColor(k.FgCorner), canvas.FontRegular, canvas.FontNormal)
			shiftText := canvas.NewTextLine(faceShift, k.Shift, canvas.Left)
			ctx.DrawText(k.X+1.5, posY+k.Height-3.5, shiftText)
		}

		if k.SecondaryShift != "" && secondaryFamily != nil {
			faceSecShift := secondaryFamily.Face(4.5, hexColor(k.FgSecondary), canvas.FontRegular, canvas.FontNormal)
			secShiftText := canvas.NewTextLine(faceSecShift, k.SecondaryShift, canvas.Left)
			ctx.DrawText(k.X+6, posY+k.Height-3.5, secShiftText)
		}

		if k.AltGr != "" {
			faceAltGr := primaryFamily.Face(4.5, hexColor(k.FgCorner), canvas.FontRegular, canvas.FontNormal)
			altGrText := canvas.NewTextLine(faceAltGr, k.AltGr, canvas.Right)
			ctx.DrawText(k.X+k.Width-5.5, posY+2, altGrText)
		}

		if k.SecondaryAltGr != "" && secondaryFamily != nil {
			faceSecAltGr := secondaryFamily.Face(4.5, hexColor(k.FgSecondary), canvas.FontRegular, canvas.FontNormal)
			secAltGrText := canvas.NewTextLine(faceSecAltGr, k.SecondaryAltGr, canvas.Right)
			ctx.DrawText(k.X+k.Width-1.5, posY+2, secAltGrText)
		}

		if k.SecondaryChar != "" && secondaryFamily != nil {
			faceSecondary := secondaryFamily.Face(6.5, hexColor(k.FgSecondary), canvas.FontRegular, canvas.FontBold)
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

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("conectando a redis: " + err.Error())
	}

	os.MkdirAll("output", 0755)

	fmt.Println("Worker escuchando pdf_jobs_queue...")

	for {
		result, err := rdb.BRPop(ctx, 0, "pdf_jobs_queue").Result()
		if err != nil {
			fmt.Println("Error leyendo de la cola:", err)
			continue
		}

		orderID, err := strconv.Atoi(result[1])
		if err != nil {
			fmt.Println("order_id inválido en la cola:", result[1])
			continue
		}

		fmt.Printf("Procesando order %d...\n", orderID)
		if err := processJob(ctx, conn, orderID); err != nil {
			fmt.Printf("Error procesando order %d: %v\n", orderID, err)
			conn.Exec(ctx, "UPDATE pdf_jobs SET status='failed' WHERE order_id=$1", orderID)
			continue
		}
		fmt.Printf("Order %d: PDF generado\n", orderID)
	}
}

func processJob(ctx context.Context, conn *pgx.Conn, orderID int) error {
	var rawConfig []byte
	err := conn.QueryRow(ctx,
		"SELECT config FROM pdf_jobs WHERE order_id=$1", orderID,
	).Scan(&rawConfig)
	if err != nil {
		return fmt.Errorf("buscando pdf_job: %w", err)
	}

	var cfg DesignConfig
	if err := json.Unmarshal(rawConfig, &cfg); err != nil {
		return fmt.Errorf("parseando config: %w", err)
	}

	design, err := buildDesign(ctx, conn, &cfg)
	if err != nil {
		return fmt.Errorf("construyendo diseño: %w", err)
	}

	outPath := fmt.Sprintf("output/order_%d.pdf", orderID)
	if err := renderDesign(design, outPath); err != nil {
		return fmt.Errorf("renderizando PDF: %w", err)
	}

	_, err = conn.Exec(ctx,
		"UPDATE pdf_jobs SET status='done', file_url=$1 WHERE order_id=$2",
		outPath, orderID,
	)
	return err
}