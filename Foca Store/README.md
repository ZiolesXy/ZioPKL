# Voca Plane

Voca Plane adalah aplikasi full-stack untuk booking tiket pesawat. Repository ini terdiri dari:

- `client`: frontend Next.js untuk user dan admin
- `server`: backend Golang + Gin + GORM + PostgreSQL

Fitur utama yang sudah terlihat dari kode:

- registrasi dan login
- pencarian flight
- detail flight dan pemilihan kursi
- pembuatan transaksi booking
- riwayat booking user
- dashboard admin
- manajemen flight, airline, airport, promo, transaksi, dan user

## Struktur Project

```text
LSP/
|- client/   # Next.js 16 + React 19
`- server/   # Go + Gin + GORM + PostgreSQL
```

## Tech Stack

### Frontend

- Next.js 16
- React 19
- TypeScript
- Tailwind CSS 4
- Axios
- TanStack Query
- shadcn/ui

### Backend

- Go
- Gin
- GORM
- PostgreSQL
- JWT
- Midtrans
- Cloudinary

## Prasyarat

Sebelum menjalankan project, pastikan sudah terpasang:

- Node.js 20+
- npm
- Go 1.22+ atau versi yang kompatibel dengan `go.mod`
- PostgreSQL

## Menjalankan Project

Karena project dipisah menjadi frontend dan backend, jalankan keduanya di terminal terpisah.

### 1. Setup Backend

Masuk ke folder server:

```powershell
cd server
```

Copy environment file:

```powershell
Copy-Item .env.example .env
```

Isi minimal konfigurasi berikut di `server/.env`:

```env
APP_PORT=8000
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=flight_booking
DB_SSLMODE=disable

JWT_SECRET=supersecretkey_production_change_me
ACCESS_TOKEN_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=168h

ALLOWED_ORIGINS=http://localhost:3000
APP_PASSWORD=isi-password-internal

MIDTRANS_SERVER_KEY=SB-Mid-server-xxxxxxxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxxxxxx
MIDTRANS_IS_PRODUCTION=false

CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
```

Catatan penting:

- Secara default file `.env.example` backend memakai `APP_PORT=8080`.
- Frontend saat ini hardcoded ke `http://localhost:8000/api/v1` di `client/lib/axios.ts`.
- Supaya langsung jalan tanpa ubah kode frontend, paling aman set `APP_PORT=8000` di backend.

Jalankan backend:

```powershell
go run cmd/main.go
```

Atau:

```powershell
make start
```

Jika sukses, API akan tersedia di:

```text
http://localhost:8000/api/v1
```

Health check:

```text
GET http://localhost:8000/api/v1/health
```

### 2. Setup Frontend

Masuk ke folder client:

```powershell
cd client
```

Install dependency:

```powershell
npm install
```

Jalankan frontend:

```powershell
npm run dev
```

Buka aplikasi:

```text
http://localhost:3000
```

## Bootstrap Database dan Data Sample

Saat server pertama kali dijalankan, aplikasi hanya melakukan seed akun user awal. Data seperti airline, airport, flight, dan promo belum otomatis dibuat.

### Akun bawaan hasil seed

- Super Admin: `superadmin@flightbooking.com` / `admin123`
- Admin: `admin@flightbooking.com` / `admin123`
- User: `user@flightbooking.com` / `user123`

### Cara isi data sample lengkap

1. Login dulu sebagai `superadmin`.
2. Ambil `access_token` dari response login.
3. Panggil endpoint system seed.

Contoh login:

```http
POST http://localhost:8000/api/v1/auth/login
Content-Type: application/json

{
  "email": "superadmin@flightbooking.com",
  "password": "admin123"
}
```

Contoh response sukses:

```json
{
  "success": true,
  "message": "login successfully",
  "data": {
    "access_token": "...",
    "refresh_token": "...",
    "expires_in": 86400,
    "token_type": "Bearer"
  }
}
```

Lalu seed full database:

```http
POST http://localhost:8000/api/v1/system/seed
Authorization: Bearer <ACCESS_TOKEN>
password-app: <APP_PASSWORD>
```

Endpoint ini akan reset database lalu mengisi ulang:

- user
- airline
- airport
- flight
- promo

## Tutorial Penggunaan

### Flow User

1. Buka halaman utama di `http://localhost:3000`.
2. Isi `Origin` dan `Destination`.
3. Klik `Book Now` atau cari lewat halaman `/flight`.
4. Pilih salah satu flight dari hasil pencarian.
5. Login atau register jika belum punya akun.
6. Di halaman detail flight, isi form booking dan pilih kursi.
7. Buat transaksi.
8. Jika transaksi masih `PENDING`, buka link pembayaran yang tersedia.
9. Cek status booking di menu `My Bookings`.

### Flow Admin

1. Login menggunakan akun admin atau super admin.
2. Akses halaman admin seperti:
   - `/dashboard`
   - `/airlines`
   - `/airports`
   - `/flights-schedule`
   - `/transactions`
   - `/users-monitoring`
3. Kelola data master seperti airline, airport, dan flight.
4. Pantau transaksi dan user dari dashboard admin.

## Endpoint Penting

### Public

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `GET /api/v1/flights`
- `GET /api/v1/flight/search`
- `GET /api/v1/flights/:id`
- `GET /api/v1/flights/:id/seats`
- `GET /api/v1/airports`
- `GET /api/v1/airlines`

### User Protected

- `GET /api/v1/user/profile`
- `PATCH /api/v1/user/profile`
- `GET /api/v1/transactions`
- `GET /api/v1/transactions/:code`
- `POST /api/v1/transactions`
- `DELETE /api/v1/transactions/:code`

### Admin Protected

- `GET /api/v1/admin/dashboard`
- `GET /api/v1/admin/users`
- `GET /api/v1/admin/transactions`
- `GET /api/v1/admin/flights`
- `POST /api/v1/admin/flights`
- `GET /api/v1/admin/airlines`
- `GET /api/v1/admin/airports`
- `GET /api/v1/admin/promos`

## Catatan Konfigurasi

### Midtrans

Project ini sudah menyiapkan callback:

```text
POST /api/v1/transactions/midtrans/callback
```

Jika ingin menguji flow pembayaran, isi kredensial sandbox Midtrans di `.env`.

### Cloudinary

Cloudinary dipakai untuk upload logo airline. Jika fitur upload logo ingin dipakai, isi konfigurasi Cloudinary di `.env`.

### CORS

Pastikan `ALLOWED_ORIGINS` di backend mengizinkan origin frontend lokal:

```env
ALLOWED_ORIGINS=http://localhost:3000
```

## Troubleshooting

### Frontend tidak bisa connect ke backend

Cek file `client/lib/axios.ts`. Saat ini base URL diset ke:

```ts
http://localhost:8000/api/v1
```

Jika backend berjalan di port lain, sesuaikan salah satu:

- ubah `APP_PORT` backend ke `8000`
- atau ubah `baseURL` frontend

### Login berhasil tapi route tetap redirect

Frontend menyimpan token di cookie `token` dan `access_token`, lalu middleware Next.js membaca cookie tersebut untuk proteksi route. Jika role atau token belum terbaca:

- logout lalu login ulang
- pastikan cookie berhasil tersimpan di browser
- pastikan backend mengembalikan token valid

### Data flight kosong

Jalankan endpoint `/api/v1/system/seed` menggunakan akun `superadmin` dan header `password-app`.

## Saran Pengembangan Lanjutan

- pindahkan `baseURL` frontend ke environment variable
- tambahkan Docker Compose untuk PostgreSQL + server + client
- tambahkan dokumentasi API Swagger/OpenAPI
- tambahkan testing backend dan frontend
- buat seed command CLI agar tidak perlu trigger via endpoint internal

## Lisensi

Belum ditentukan.
