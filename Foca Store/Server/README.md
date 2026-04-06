# Foca Store Backend API

Aplikasi backend e-commerce sederhana menggunakan Golang, Gin, GORM, PostgreSQL, dan JWT.

## Fitur Utama

- **Authentication**: Register, Login, Access Token (15 menit), Refresh Token (7 hari)
- **Role-based Authorization**: Admin dan User
- **Product Management**: CRUD produk (Admin only)
- **Category Management**: CRUD kategori produk (Admin only)
- **Address Management**: CRUD alamat pengiriman user
- **Cart System**: Add to cart, View cart, Remove item
- **Checkout System**: Checkout dengan status (pending, success, failed)
- **Coupon System**: Buat, klaim, dan gunakan kupon diskon (dengan masa berlaku dan status aktif)
- **System Management**: Endpoint untuk migrasi database, reset database, dan seeding data
- **Security**: Generator secret key untuk JWT dan proteksi endpoint sistem
- **Image Upload**: Cloudinary integration untuk product images dan profile images

---

## Instalasi

### Prasyarat

- Go 1.24+
- PostgreSQL 12+
- Git

### Langkah Instalasi

1. Clone repository ini:

```bash
git clone <repository-url>
cd "Foca Store Test"
```

2. Install dependencies:

```bash
go mod download
```

3. Buat file `.env` dari contoh:

```bash
cp .env.example .env
```

4. Edit file `.env` sesuai konfigurasi database Anda:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=foca_store
JWT_SECRET=your_super_secret_jwt_key_here_min_32_chars
SYSTEM_PASSWORD=your_system_password_here
CLOUDINARY_CLOUD_NAME=your_cloudinary_cloud_name
CLOUDINARY_API_KEY=your_cloudinary_api_key
CLOUDINARY_API_SECRET=your_cloudinary_api_secret
PORT=8080
```

---

## Setup Database

1. Buat database PostgreSQL:

```sql
CREATE DATABASE foca_store;
```

2. Pastikan user PostgreSQL memiliki akses ke database tersebut.

---

## Migrasi Database

Aplikasi menggunakan auto-migration GORM.  
Tabel akan dibuat otomatis saat aplikasi pertama kali dijalankan.  
Anda juga dapat melakukan migrasi manual melalui endpoint `/system/migrate`.

---

## Menjalankan Aplikasi

1. Jalankan aplikasi:

```bash
go run main.go
```

2. Aplikasi akan berjalan di:

```
http://localhost:8080
```

---

## System & Seeder

Semua endpoint sistem dan seeding dilindungi oleh `SYSTEM_PASSWORD` yang harus dikirim melalui `form-data` dengan key `password`.

### System Endpoints

- `POST /system/migrate` - Menjalankan migrasi database
- `POST /system/reset` - Reset database (drop semua tabel dan migrasi ulang)

### Seeder Endpoints

Untuk mengisi data awal, gunakan endpoint berikut:

```bash
# Seed roles (Admin, User)
GET http://localhost:8080/system/seed/roles

# Seed admin user
GET http://localhost:8080/system/seed/admin

# Seed sample users
GET http://localhost:8080/system/seed/users

# Seed sample products
GET http://localhost:8080/system/seed/products

# Seed products from assets folder
GET http://localhost:8080/system/seed/assets

# Seed sample coupons
GET http://localhost:8080/system/seed/coupons

# Sync asset products
PUT http://localhost:8080/system/seed/sync

# Seed all data (roles, admin, user, categories, coupons)
GET http://localhost:8080/system/seed/all
```

Data yang di-seed:

- **Roles**: Admin, User
- **Admin**: Email sesuai seeder (e.g., admin@foca-store.com), password: (lihat `seeders/seed.go`)
- **Users**: Beberapa sample users untuk testing
- **Products**: Berbagai kategori produk (Laptop, Smartphone, dsb)
- **Coupons**: Sample coupons dengan tipe percentage dan fixed

---

## API Endpoints

### Public Endpoints

- `GET /password` - Generate rekomendasi JWT secret key
- `POST /register` - Register user baru
- `POST /login` - Login user
- `POST /refresh` - Refresh access token
- `GET /category` - Lihat semua kategori
- `GET /category/:slug` - Lihat detail kategori berdasarkan slug
- `GET /products` - Lihat semua produk
- `GET /product/:slug` - Lihat detail produk berdasarkan slug
- `GET /coupons` - Lihat semua kupon yang tersedia

### User Profile (Protected)

- `GET /api/profile` - Lihat profil user
- `PUT /api/profile` - Update profil user (support upload gambar)

### Address Management (Protected)

- `POST /api/addresses` - Tambah alamat baru
- `GET /api/addresses` - Lihat semua alamat user
- `GET /api/addresses/:uid` - Lihat detail alamat berdasarkan UID
- `PUT /api/addresses/:uid` - Update alamat
- `DELETE /api/addresses/:uid` - Hapus alamat

### Cart (Protected)

- `GET /api/cart` - Lihat keranjang
- `POST /api/cart/items` - Tambah item ke keranjang
- `DELETE /api/cart/items/:id` - Hapus item dari keranjang

### Checkout (Protected)

- `POST /api/checkout` - Checkout keranjang
- `GET /api/checkout/me` - Lihat history checkout user

### Coupon Management (Protected)

- `POST /api/coupons/claim` - Klaim kupon ke akun user
- `GET /api/coupons/me` - Lihat kupon yang dimiliki user
- `DELETE /api/coupons/:id/remove` - Hapus kupon dari akun user

### Admin Management (Admin Only)

**Product Admin:**
- `POST /api/admin/products` - Buat produk baru
- `PUT /api/admin/products/:id` - Update produk
- `DELETE /api/admin/products/:id` - Hapus produk
- `DELETE /api/admin/products` - Hapus semua produk
- `DELETE /api/admin/products/assets` - Hapus semua gambar produk dari Cloudinary

**Category Admin:**
- `POST /api/admin/category` - Buat kategori baru
- `PUT /api/admin/category/:id` - Update kategori
- `DELETE /api/admin/category/:id` - Hapus kategori

**Coupon Admin:**
- `POST /api/admin/coupons` - Buat kupon baru
- `PUT /api/admin/coupon/:id` - Update kupon
- `DELETE /api/admin/coupon/:id` - Hapus kupon

**Checkout Admin:**
- `GET /api/admin/checkout` - Lihat semua history checkout
- `PATCH /api/admin/checkout/:id/approve` - Approve transaksi
- `PATCH /api/admin/checkout/:id/reject` - Reject transaksi

---

## Format Response

### Success (single)

```json
{
  "status": "success",
  "message": "Success message",
  "data": {}
}
```

### Success (list)

```json
{
  "status": "success",
  "message": "Success message",
  "data": {
    "entries": []
  }
}
```

### Error

```json
{
  "status": "error",
  "message": "Error message"
}
```

---

## Keamanan

- Password di-hash menggunakan bcrypt
- JWT untuk autentikasi (Access & Refresh tokens)
- Role-based authorization (Admin / User)
- System management protection menggunakan static password
- Cloudinary integration untuk secure image hosting

---

## Struktur Folder

```
Foca Strore Test/
├── main.go
├── database/
│   └── database.go
├── models/
│   ├── role.go, user.go, category.go, product.go, ...
│   └── coupon.go, user_coupon.go, address.go
├── request/
│   └── auth.go, product.go, checkout.go, address.go, ...
├── response/
│   └── response.go, auth.go, product.go, coupon.go, ...
├── handlers/
│   ├── auth.go, product.go, category.go, coupon.go, ...
│   ├── checkout.go, address.go, cart.go, seed.go, ...
│   └── Security.go
├── middleware/
│   ├── jwt.go, admin.go, system.go
├── helper/
│   ├── jwt.go, cloudinary.go, slug.go, password.go, context.go
├── seeders/
│   ├── seed.go, reset.go, migrate.go, drop.go
├── AssetPrivate/
└── README.md
```
