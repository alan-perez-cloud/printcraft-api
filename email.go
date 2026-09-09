package main

import (
	"fmt"
	"os"
)

type OrderEmailData struct {
	CustomerEmail     string
	OrderID           int
	FulfillmentType   string
	PrimaryAlphabet   string
	SecondaryAlphabet string
	KeycapMode        string
	DownloadURL       string // solo para digital
}

func sendCustomerEmail(data OrderEmailData) error {
	var subject, body string

	if data.FulfillmentType == "digital" {
		subject = "Tu pedido de LePrintCraft está listo 🎉"
		body = fmt.Sprintf(
			"Hola,\n\nGracias por tu compra. Aquí está tu diseño:\n\nAlfabeto principal: %s\nAlfabeto secundario: %s\nEstilo de tecla: %s\n\nDescarga tu PDF aquí: %s\n\nLePrintCraft",
			data.PrimaryAlphabet, data.SecondaryAlphabet, data.KeycapMode, data.DownloadURL,
		)
	} else {
		subject = "Tu pedido de LePrintCraft está confirmado"
		body = fmt.Sprintf(
			"Hola,\n\nTu pedido fue confirmado y está en preparación:\n\nAlfabeto principal: %s\nAlfabeto secundario: %s\nEstilo de tecla: %s\n\nTe avisaremos cuando esté en camino.\n\nLePrintCraft",
			data.PrimaryAlphabet, data.SecondaryAlphabet, data.KeycapMode,
		)
	}

	return dispatchEmail(data.CustomerEmail, subject, body)
}

func sendAdminNotification(data OrderEmailData, amount int, currency string) error {
	subject := fmt.Sprintf("🔔 Nueva venta - Orden #%d", data.OrderID)
	body := fmt.Sprintf(
		"Cliente: %s\nTipo: %s\nMonto: %d %s\nAlfabetos: %s + %s",
		data.CustomerEmail, data.FulfillmentType, amount, currency, data.PrimaryAlphabet, data.SecondaryAlphabet,
	)
	return dispatchEmail(os.Getenv("ADMIN_EMAIL"), subject, body)
}

func dispatchEmail(to, subject, body string) error {
	if os.Getenv("ENV") != "production" {
		fmt.Printf("(dev) Email a %s\nAsunto: %s\n%s\n---\n", to, subject, body)
		return nil
	}
	// TODO: aquí va la llamada real a SES
	return nil
}