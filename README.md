# Inventory API

Small inventory-management REST API for backend live-coding interviews.

The code is split into `handler -> service -> repository`. `cmd/server` is the composition root: it opens SQLite, constructs the repository, injects it into the service, and injects the service into the HTTP handler.

## Run

Requires Go 1.24+.

```sh
go run ./cmd/server
```

The server listens on `http://localhost:8080` and creates `data/inventory.db`.
Set `DB_PATH` to use another SQLite file.

## API

```sh
curl localhost:8080/healthz

curl -X POST localhost:8080/items \
  -H 'Content-Type: application/json' \
  -d '{"sku":"KB-001","name":"Keyboard","quantity":10}'

curl localhost:8080/items
curl localhost:8080/items/1

curl -X PATCH localhost:8080/items/1/stock \
  -H 'Content-Type: application/json' \
  -d '{"delta":-2}'
```

Versioned migrations in `cmd/server/migrations` are embedded and applied automatically at startup using SQLite's native `user_version`.
