package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var db *pgxpool.Pool
var rdb *redis.Client

type Design struct {
	ID         int             `json:"id"`
	Config     json.RawMessage `json:"config"`
	GuestToken string          `json:"guest_token"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No se pudo cargar .env:", err)
		os.Exit(1)
	}

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		fmt.Println("Error al conectar:", err)
		os.Exit(1)
	}
	defer db.Close()

	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		fmt.Println("Error al conectar a Redis:", err)
		os.Exit(1)
	}
	fmt.Println("Conectado exitosamente a Redis")
	fmt.Println("Conectado exitosamente a PostgreSQL")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/designs", createDesign)
	mux.HandleFunc("GET /api/v1/designs/{id}", getDesign)
	mux.HandleFunc("PUT /api/v1/designs/{id}", updateDesign)
	mux.HandleFunc("DELETE /api/v1/designs/{id}", deleteDesign)
	mux.HandleFunc("POST /api/v1/orders", createOrder)
	mux.HandleFunc("POST /api/v1/auth/register", register)
	mux.HandleFunc("POST /api/v1/auth/login", login)
	mux.HandleFunc("POST /api/v1/webhooks/lemonsqueezy", lemonSqueezyWebhook)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", corsMiddleware(mux))
}

func createDesign(w http.ResponseWriter, r *http.Request) {
	var d Design
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	err := db.QueryRow(context.Background(),
		"INSERT INTO designs (config, guest_token) VALUES ($1, $2) RETURNING id",
		d.Config, d.GuestToken,
	).Scan(&d.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func getDesign(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var d Design

	err := db.QueryRow(context.Background(),
		"SELECT id, config, guest_token FROM designs WHERE id=$1", id,
	).Scan(&d.ID, &d.Config, &d.GuestToken)

	if err != nil {
		http.Error(w, "No encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

type Order struct {
	ID         int     `json:"id"`
	DesignID   int     `json:"design_id"`
	GuestToken string  `json:"guest_token"`
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
}

func updateDesign(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var d Design
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(context.Background(),
		"UPDATE designs SET config=$1 WHERE id=$2", d.Config, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Actualizado")
}

func deleteDesign(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	_, err := db.Exec(context.Background(), "DELETE FROM designs WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Eliminado")
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var o Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := db.QueryRow(ctx,
		"INSERT INTO orders (design_id, guest_token, status, amount) VALUES ($1, $2, 'pending', $3) RETURNING id",
		o.DesignID, o.GuestToken, o.Amount,
	).Scan(&o.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	checkoutURL, err := CreateCheckout(o.ID, o.Amount, o.GuestToken)
	if err != nil {
		http.Error(w, "no se pudo generar el pago", http.StatusBadGateway)
		return
	}

	o.Status = "pending"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":           o.ID,
		"status":       o.Status,
		"checkout_url": checkoutURL,
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}