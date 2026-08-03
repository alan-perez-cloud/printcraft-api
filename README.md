# PrintCraft API

Backend in Go for PrintCraft, a SaaS platform for personalized printable products. First product: **KeyStickers** — customizable keyboard stickers for language learners.

## Stack

- **Go** 1.26 + `net/http` (native router, no external frameworks)
- **PostgreSQL** 18 (driver `pgx/v5`)
- Environment variables via `godotenv`

## Architecture

A reusable engine designed to support multiple printable products (KeyStickers, coloring pages, crosswords, etc.) on top of a shared core: auth, checkout, payments, queues, and emails.

### Current flow

1. Client creates a `design` (personalized configuration as JSON)
2. On payment, an `order` is created (transactional)
3. A `pdf_job` is automatically generated with a snapshot of the design, ready to be processed
4. *(Coming soon)* A worker consumes the queue, generates the PDF, uploads it to storage, and emails the download link

## Current Endpoints

| Method | Route                        | Description                                     |
|--------|-------------------------------|--------------------------------------------------|
| POST   | `/api/v1/designs`             | Create a design                                   |
| GET    | `/api/v1/designs/{id}`        | Get a design                                      |
| PUT    | `/api/v1/designs/{id}`        | Update a design                                   |
| DELETE | `/api/v1/designs/{id}`        | Delete a design                                   |
| POST   | `/api/v1/orders`               | Create a paid order (simulated) and trigger the `pdf_job` |

## Running locally

```bash
# 1. Install PostgreSQL and create the database
createdb printcraft

# 2. Copy .env.example to .env and fill in your credentials
cp .env.example .env

# 3. Install dependencies
go mod download

# 4. Run the server
go run main.go
```

Server available at `http://localhost:8080`.

## Roadmap

- [ ] Redis-based job queue
- [ ] PDF generation worker
- [ ] Email delivery (Resend)
- [ ] Stripe integration
- [ ] Auth + guest checkout