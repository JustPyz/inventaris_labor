# invela-be

Backend REST API untuk **inventaris laboratorium sekolah**: pengelolaan labor, perangkat, unit fisik (item instance), peminjaman, kerusakan, perbaikan, beserta audit log-nya.

Stack: **Go + Gin + GORM + SQLite**, autentikasi **JWT**, otorisasi berbasis **role**.

## Struktur project

- `main.go` — entry point & composition root (wiring repo → service → handler)
- `internal/config` — konfigurasi via environment variable
- `internal/database` — koneksi SQLite, AutoMigrate, seed role & admin
- `internal/router` — definisi route + RBAC per-endpoint + `/health`
- `internal/middlewares` — `RequireAuth` (JWT), `RequireRole`, `CORS`
- `internal/models` — struct GORM (User, Role, Jurusan, Kelas, Kategori, Labor, Perangkat, ItemInstance, Peminjaman, Penggunaan, Kerusakan, Perbaikan, RiwayatPerbaikan)
- `internal/repositories` — akses data per-entitas
- `internal/services` — validasi & business logic (termasuk transaksi)
- `internal/handlers` — parse request, panggil service, map error → HTTP
- `internal/jobs` — background job (penanda peminjaman melewati batas waktu)

## Menjalankan

```bash
go mod tidy
go run .
```

Server default berjalan di `http://localhost:8080`. Cek kesehatan:

```bash
curl http://localhost:8080/health
```

Saat pertama kali dijalankan, database `data/app.db` dibuat otomatis, seluruh tabel di-migrate, lalu role dan satu user `admin` di-seed.

## Konfigurasi (environment variable)

Salin `.env.example` menjadi `.env` untuk pemakaian lokal, lalu isi nilainya:

```bash
cp .env.example .env
```

> `.env` (berisi nilai asli) **tidak** ikut ter-commit. `.env.example` (tanpa nilai rahasia) **boleh** di-commit sebagai dokumentasi.

| Variabel     | Wajib?                     | Default        | Keterangan |
|--------------|----------------------------|----------------|------------|
| `APP_ENV`    | tidak                      | `development`  | `development` atau `production`. Di `production`, `JWT_SECRET` menjadi wajib. |
| `APP_PORT`   | tidak                      | `8080`         | Port HTTP server. |
| `DB_PATH`    | tidak                      | `data/app.db`  | Lokasi file database SQLite. |
| `JWT_SECRET` | **ya, di `production`**     | —              | Kunci penandatangan JWT. Di `development` boleh kosong (memakai secret default yang **tidak aman**, disertai peringatan di log). Di `production`, jika kosong aplikasi **berhenti saat start**. |

**Membuat nilai `JWT_SECRET`** (jalankan salah satu, lalu tempel hasilnya sebagai nilai `JWT_SECRET` di server — **jangan** menyimpan nilai aslinya di repo/README):

```bash
openssl rand -base64 48
# atau
head -c 48 /dev/urandom | base64
```

Contoh menjalankan di produksi:

```bash
export APP_ENV=production
export JWT_SECRET="<tempel-hasil-perintah-di-atas>"
go run .
```

## Autentikasi & role

- Login menghasilkan JWT; sertakan pada request berikutnya via header `Authorization: Bearer <token>`.
- Role yang tersedia: `admin`, `kabeng`, `kaprog`, `guru`, `sapras`.
- Otorisasi diterapkan per-endpoint melalui middleware `RequireAuth` + `RequireRole`.

User `admin` default di-seed saat inisialisasi database (kredensial awal ada di `internal/database`). **Ganti password admin default sebelum dipakai di produksi.**

## Endpoint (ringkasan)

Base path: `/api` (kecuali `/health`).

- `GET  /health` — health check (publik)
- `POST /api/auth/login` — login, mengembalikan JWT (publik)
- `/api/kelas`, `/api/jurusan` — data akademik
- `/api/user` — manajemen user
- `/api/kategori`, `/api/labor`, `/api/perangkat` — master inventaris
- `/api/item-instance` — unit fisik per perangkat (status: aktif / dipinjam / rusak)
- `/api/peminjaman`, `/api/penggunaan` — transaksi peminjaman & penggunaan
- `/api/kerusakan`, `/api/perbaikan`, `/api/riwayat-perbaikan` — siklus kerusakan → perbaikan

Detail method, payload, dan role yang diizinkan tiap endpoint dapat dilihat di `internal/router/router.go`.

## Background job

`PeminjamanOverdueJob` berjalan periodik (interval 1 jam) untuk menandai peminjaman yang melewati batas waktu pengembalian.
