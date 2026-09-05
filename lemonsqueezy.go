package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func CreateCheckout(orderID int, product string, amount float64, guestToken string) (string, error) {
	variantID := os.Getenv("LEMONSQUEEZY_VARIANT_PDF")
	if product == "print" {
		variantID = os.Getenv("LEMONSQUEEZY_VARIANT_PRINT")
	}

	payload := map[string]any{
		"data": map[string]any{
			"type": "checkouts",
			"attributes": map[string]any{
				"checkout_data": map[string]any{
					"custom": map[string]string{
						"order_id": fmt.Sprint(orderID),
					},
				},
			},
			"relationships": map[string]any{
				"store": map[string]any{
					"data": map[string]string{"type": "stores", "id": os.Getenv("LEMONSQUEEZY_STORE_ID")},
				},
				"variant": map[string]any{
					"data": map[string]string{"type": "variants", "id": variantID},
				},
			},
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.lemonsqueezy.com/v1/checkouts", bytes.NewReader(body))
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("LEMONSQUEEZY_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("lemon squeezy respondió %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Attributes struct {
				URL string `json:"url"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Data.Attributes.URL, nil
}

func lemonSqueezyWebhook(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	mac := hmac.New(sha256.New, []byte(os.Getenv("LEMONSQUEEZY_WEBHOOK_SECRET")))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	signature := r.Header.Get("X-Signature")

	if !hmac.Equal([]byte(expected), []byte(signature)) {
		http.Error(w, "firma inválida", http.StatusUnauthorized)
		return
	}

	var event struct {
		Meta struct {
			EventName  string `json:"event_name"`
			CustomData struct {
				OrderID string `json:"order_id"`
			} `json:"custom_data"`
		} `json:"meta"`
		Data struct {
			Attributes struct {
				UserEmail string `json:"user_email"`
				TestMode  bool   `json:"test_mode"`
				Total     int    `json:"total"`
				Currency  string `json:"currency"`
				Urls      struct {
					Receipt string `json:"receipt"`
				} `json:"urls"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if event.Meta.EventName != "order_created" {
		w.WriteHeader(http.StatusOK)
		return
	}

	orderID, err := strconv.Atoi(event.Meta.CustomData.OrderID)
	if err != nil {
		http.Error(w, "order_id inválido", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	var designID int
	err = tx.QueryRow(ctx,
		`UPDATE orders 
		 SET status='paid', customer_email=$2, test_mode=$3, receipt_url=$4, ls_total=$5, ls_currency=$6
		 WHERE id=$1 RETURNING design_id`,
		orderID,
		event.Data.Attributes.UserEmail,
		event.Data.Attributes.TestMode,
		event.Data.Attributes.Urls.Receipt,
		event.Data.Attributes.Total,
		event.Data.Attributes.Currency,
	).Scan(&designID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var config json.RawMessage
	tx.QueryRow(ctx, "SELECT config FROM designs WHERE id=$1", designID).Scan(&config)

	_, err = tx.Exec(ctx,
		"INSERT INTO pdf_jobs (order_id, status, config) VALUES ($1, 'pending', $2)",
		orderID, config,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rdb.LPush(ctx, "pdf_jobs_queue", orderID)
	w.WriteHeader(http.StatusOK)
}