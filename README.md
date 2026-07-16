# Invela Go

Backend API untuk pengelolaan data `role` menggunakan Gin, Gorm, dan SQLite.

## Fitur awal

- Health check
- Bootstrap Gin + Gorm + SQLite
- CRUD `role` dengan kolom `id` dan `role`

## Struktur project

- `cmd/server`: entry point aplikasi
- `internal/config`: konfigurasi aplikasi
- `internal/database`: koneksi SQLite
- `internal/database`: koneksi SQLite
- `internal/router`: routing aplikasi
- `internal/models/role.go`: model tabel `role`

## Menjalankan

```bash
go mod tidy
go run ./cmd/server
```

## Endpoint

- `GET /health`
- `GET /api/v1/roles`
- `POST /api/v1/roles`
- `GET /api/v1/roles/:id`
- `PUT /api/v1/roles/:id`
- `DELETE /api/v1/roles/:id`

## Environment

- `APP_PORT` default: `8080`
- `DB_PATH` default: `data/app.db`
