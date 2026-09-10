package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
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

	ctx := context.Background()
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("cargando config aws: %w", err)
	}
	client := sesv2.NewFromConfig(cfg)

	from := os.Getenv("SES_FROM_EMAIL")
	if from == "" {
		return fmt.Errorf("SES_FROM_EMAIL no está seteado")
	}

	_, err = client.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(from),
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(subject)},
				Body:    &types.Body{Text: &types.Content{Data: aws.String(body)}},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("enviando email via SES: %w", err)
	}
	return nil
}