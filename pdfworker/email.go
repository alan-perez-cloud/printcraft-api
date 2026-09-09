
package main

import (
	"fmt"
	"os"
)

type OrderEmailData struct {
	CustomerEmail     string
	PrimaryAlphabet   string
	SecondaryAlphabet string
	KeycapMode        string
	DownloadURL       string
}

func sendDigitalEmail(data OrderEmailData) error {
	subject := "Tu pedido de LePrintCraft está listo 🎉"
	body := fmt.Sprintf(
		"Hola,\n\nGracias por tu compra. Aquí está tu diseño:\n\nAlfabeto principal: %s\nAlfabeto secundario: %s\nEstilo de tecla: %s\n\nDescarga tu PDF aquí: %s\n\nLePrintCraft",
		data.PrimaryAlphabet, data.SecondaryAlphabet, data.KeycapMode, data.DownloadURL,
	)
	return dispatchEmail(data.CustomerEmail, subject, body)
}

func dispatchEmail(to, subject, body string) error {
	if os.Getenv("ENV") != "production" {
		fmt.Printf("(dev) Email a %s\nAsunto: %s\n%s\n---\n", to, subject, body)
		return nil
	}
	// TODO: llamada real a SES
	return nil
}