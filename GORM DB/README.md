# GORM Auth API - Sistem Manajemen User dan Barang

## Deskripsi Project

Project ini adalah **RESTful API** yang dibangun menggunakan **Go (Golang)** dengan framework **Gin** dan ORM **GORM** untuk mengelola sistem autentikasi user dan manajemen barang. Project ini menerapkan arsitektur yang terstruktur dengan pemisahan tanggung jawab yang jelas antara handlers, models, middleware, dan database layer.

## Fitur Utama

### 🔐 Autentikasi & Autorisasi
- **Register** user baru dengan password hashing
- **Login** dengan JWT token (Access Token & Refresh Token)
- **Refresh Token** untuk memperbarui access token
- **Logout** untuk menghapus refresh token
- **Role-based access control** (Admin & User biasa)

### 👥 Manajemen User
- CRUD operations untuk users
- Get profile user yang sedang login
- Assign role ke user (Admin only)
- Many-to-many relationship dengan barang

### 📦 Manajemen Barang
- CRUD operations untuk barang
- Assign barang ke user (Admin only)
- Remove barang dari user (Admin only)
- Many-to-many relationship dengan users

### 🏷️ Manajemen Role
- CRUD operations untuk roles
- Default roles: Admin (ID: 1) dan User (ID: 2)

## Tech Stack

- **Backend**: Go 1.24.11
- **Web Framework**: Gin v1.11.0
- **ORM**: GORM v1.31.1
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt

## Struktur Project

```
GORM DB/
├── main.go                 # Entry point aplikasi
├── go.mod                  # Dependency management
├── database/
│   └── database.go         # Koneksi database
├── models/
│   ├── user.go            # Model User
│   ├── role.go            # Model Role
│   ├── barang.go          # Model Barang
│   └── refresh_token.go   # Model RefreshToken
├── handlers/
│   ├── auth_handler.go    # Handler autentikasi
│   ├── user_handler.go    # Handler user management
│   ├── barang_handler.go  # Handler barang management
│   ├── role_handler.go    # Handler role management
│   └── ownership_handler.go # Handler user-barang relationship
├── middlewares/
│   ├── auth.go            # JWT authentication middleware
│   └── admin.go           # Admin-only middleware
├── request/
│   └── *.go               # Request structs
├── respons/
│   └── *.go               # Response structs
├── helpers/
│   └── *.go               # Helper functions
└── seeders/
    └── *.go               # Database seeders
```

## Database Schema

### Users
- `id` (Primary Key)
- `name` (String)
- `email` (Unique)
- `password` (Hashed)
- `role_id` (Foreign Key ke roles)

### Roles
- `id` (Primary Key)
- `name` (String)

### Barangs
- `id` (Primary Key)
- `name` (String)

### User_Barangs (Many-to-Many)
- `user_id` (Foreign Key)
- `barang_id` (Foreign Key)

### Refresh Tokens
- `id` (Primary Key)
- `token` (String)
- `user_id` (Foreign Key)
- `expires_at` (Timestamp)

## API Endpoints

### Public Routes
- `POST /register` - Register user baru
- `POST /login` - Login user
- `POST /refresh-token` - Refresh access token
- `POST /logout` - Logout user
- `POST /dummy` - Run database seeders

### Authenticated Routes
- `GET /me` - Get profile user yang sedang login
- `GET /profile/:id` - Get user by ID
- `GET /users` - Get all users
- `GET /user/:id` - Get user with barangs by ID
- `GET /roles` - Get all roles
- `GET /role/:id` - Get role by ID
- `GET /barangs` - Get all barangs
- `GET /barang/:id` - Get barang by ID
- `GET /users/barangs` - Get all user-barang relationships
- `GET /user/barang` - Get user-barang pivot data

### Admin Only Routes
- `POST /roles` - Create role
- `PUT /role/:id` - Update role
- `DELETE /role/:id` - Delete role
- `PUT /user/:id` - Update user
- `DELETE /user/:id` - Delete user
- `POST /barangs` - Create barang
- `PUT /barang/:id` - Update barang
- `DELETE /barang/:id` - Delete barang
- `POST /user/:user_id/barang/:barang_id` - Assign barang to user
- `DELETE /user/:id/barang/:barang_id` - Remove barang from user

## Cara Menjalankan

### Prerequisites
- Go 1.24.11 atau lebih tinggi
- PostgreSQL database
- Git

### Setup

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd "GORM DB"
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup database**
   - Buat database PostgreSQL dengan nama `gorm_auth`
   - Update koneksi database di `database/database.go` jika perlu:
   ```go
   dsn := "host=localhost user=postgres password=360589 dbname=gorm_auth port=5432"
   ```

4. **Run application**
   ```bash
   go run main.go
   ```

5. **Run seeders** (Optional - untuk data awal)
   ```bash
   curl -X POST http://localhost:8080/dummy
   ```

## Default Credentials

Setelah menjalankan seeders, Anda akan memiliki:
- **Admin**: email: `admin@example.com`, password: `password123`
- **User**: email: `user@example.com`, password: `password123`

## Konfigurasi

### Database
Update konfigurasi database di `database/database.go`:
```go
dsn := "host=localhost user=postgres password=your_password dbname=your_db port=5432"
```

### JWT Secret
Update JWT secret di `helpers` package:
```go
var ACCESS_SECRET = []byte("your-access-secret")
var REFRESH_SECRET = []byte("your-refresh-secret")
```

## Contoh Penggunaan

### Register User
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get Profile (Authenticated)
```bash
curl -X GET http://localhost:8080/me \
  -H "Authorization: Bearer <access_token>"
```

## Keamanan

- Password di-hash menggunakan bcrypt
- JWT tokens dengan expiration time
- Role-based access control
- CORS middleware untuk cross-origin requests
- Input validation pada request bodies

## Development Notes

- Project ini menggunakan arsitektur modular dengan pemisahan tanggung jawad yang jelas
- Error handling yang konsisten di seluruh aplikasi
- Response format yang terstandardisasi
- Database migrations otomatis saat aplikasi dijalankan

## Author

Project ini dikembangkan sebagai pembelajaran implementasi GORM dengan PostgreSQL dan autentikasi JWT di Go.

## License

Project ini untuk tujuan pembelajaran dan pengembangan.