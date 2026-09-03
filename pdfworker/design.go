package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// DesignConfig es la forma real de la columna `config` (jsonb) en `designs`.
type DesignConfig struct {
	PrimaryAlphabet   string  `json:"primary_alphabet"`
	SecondaryAlphabet *string `json:"secondary_alphabet"`
	SecondaryColor    *string `json:"secondary_color"`
	KeyMode           string  `json:"keycap_mode"` // ojo: coincide con lo que manda el frontend, no "key_mode"
}

// fetchDesign trae el config de un design por id.
func fetchDesign(ctx context.Context, conn *pgx.Conn, id int) (*DesignConfig, error) {
	var raw []byte
	err := conn.QueryRow(ctx, "SELECT config FROM designs WHERE id = $1", id).Scan(&raw)
	if err != nil {
		return nil, fmt.Errorf("querying design %d: %w", id, err)
	}

	var cfg DesignConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

// LayoutKey es una tecla dentro de la definición de un alfabeto en `keyboard_layouts`.
type LayoutKey struct {
	KeyCode    string  `json:"key_code"`
	Base       string  `json:"base"`
	Shift      *string `json:"shift"`
	AltGr      *string `json:"altgr"`
	AltGrShift *string `json:"altgr_shift"`
}

// fetchKeyboardLayout trae la definición completa de teclas de un alfabeto.
func fetchKeyboardLayout(ctx context.Context, conn *pgx.Conn, alphabetName string) ([]LayoutKey, error) {
	var raw []byte
	err := conn.QueryRow(ctx,
		"SELECT keys FROM keyboard_layouts WHERE alphabet_name = $1", alphabetName,
	).Scan(&raw)
	if err != nil {
		return nil, fmt.Errorf("querying layout %q: %w", alphabetName, err)
	}

	var keys []LayoutKey
	if err := json.Unmarshal(raw, &keys); err != nil {
		return nil, fmt.Errorf("parsing keys: %w", err)
	}
	return keys, nil
}