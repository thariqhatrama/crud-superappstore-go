````markdown
# FinalTask Back-End (Rakamin × Evermos)

A clean-architecture REST API built with Go, Fiber, GORM and MySQL. Supports:

- Authentication (JWT) & role-based access (admin vs user)  
- User profile & address management  
- Store (toko) management  
- Category management (admin only)  
- Product CRUD + image upload + pagination & filtering  
- Transaction handling + audit snapshots (`log_produk`)  
- Public region lookup (Province & Regency) via Emsifa API  

---

## 🛠️ Prerequisites

- Go 1.20+  
- MySQL 5.7+ (or compatible)  
- [Postman](https://www.postman.com/) (for testing)  

---

## ⚙️ Setup

1. **Clone repository**  
   ```bash
   git clone https://github.com/thariqhatarama/crud-superappstore-go.git
   cd finaltask-backend
````

2. **Create database**

   ```sql
   CREATE DATABASE finaltaskrakamin;
   ```

3. **Configure `.env`**
   Copy `.env.example` → `.env` and edit:

   ```dotenv
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_USER=root
   DB_PASS=123456
   DB_NAME=finaltaskrakamin

   JWT_SECRET=your_jwt_secret_here
   ```

4. **Run migrations**
   The app auto-migrates on startup (via GORM). Just start the server.

5. **Start server**

   ```bash
   cd cmd/app
   go run main.go
   ```

   > Server listens on `:8000` by default.

---

## 📡 API Endpoints

All routes are prefixed with `/api/v1`

### Auth

| Method | Path             | Auth | Body                                 |
| ------ | ---------------- | ---- | ------------------------------------ |
| POST   | `/auth/register` | ❌    | `{ nama, email, no_telp, password }` |
| POST   | `/auth/login`    | ❌    | `{ email, password }`                |

### Users

| Method | Path        | Auth | Description    |
| ------ | ----------- | ---- | -------------- |
| GET    | `/users/me` | ✅    | Get my profile |
| PUT    | `/users/me` | ✅    | Update profile |

### Store (Toko)

| Method | Path     | Auth | Description     |
| ------ | -------- | ---- | --------------- |
| GET    | `/store` | ✅    | Get my store    |
| PUT    | `/store` | ✅    | Update my store |

### Addresses

| Method | Path                              | Auth | Description                           |
| ------ | --------------------------------- | ---- | ------------------------------------- |
| GET    | `/addresses`                      | ✅    | List my addresses                     |
| POST   | `/addresses`                      | ✅    | Create new address                    |
| GET    | `/addresses/:id`                  | ✅    | Get one of my addresses               |
| PUT    | `/addresses/:id`                  | ✅    | Update my address                     |
| DELETE | `/addresses/:id`                  | ✅    | Delete my address                     |
| GET    | `/addresses/provinces`            | ❌    | List all provinces (Emsifa API proxy) |
| GET    | `/addresses/provinces/:id`        | ❌    | Get one province by ID                |
| GET    | `/addresses/regencies/:prov_id`   | ❌    | List regencies by province            |
| GET    | `/addresses/regencies/detail/:id` | ❌    | Get one regency by ID                 |

### Categories (Admin only)

| Method | Path              | Auth    | Body                |
| ------ | ----------------- | ------- | ------------------- |
| GET    | `/categories`     | ✅ Admin | —                   |
| POST   | `/categories`     | ✅ Admin | `{ nama_category }` |
| PUT    | `/categories/:id` | ✅ Admin | `{ nama_category }` |
| DELETE | `/categories/:id` | ✅ Admin | —                   |

### Products

| Method | Path                   | Auth | Query                        | Body (JSON / FormData)                                                                 |
| ------ | ---------------------- | ---- | ---------------------------- | -------------------------------------------------------------------------------------- |
| GET    | `/products`            | ❌    | `?page=&limit=&id_category=` | —                                                                                      |
| GET    | `/products/:id`        | ❌    | —                            | —                                                                                      |
| POST   | `/products`            | ✅    | —                            | `{ nama_produk, slug?, harga_reseller, harga_konsumen, stok, deskripsi, id_category }` |
| PUT    | `/products/:id`        | ✅    | —                            | `{ nama_produk, slug?, harga_reseller, harga_konsumen, stok, deskripsi, id_category }` |
| DELETE | `/products/:id`        | ✅    | —                            | —                                                                                      |
| POST   | `/products/:id/upload` | ✅    | —                            | FormData `file` field (image)                                                          |

### Transactions

| Method | Path                | Auth | Body                                                                         |
| ------ | ------------------- | ---- | ---------------------------------------------------------------------------- |
| GET    | `/transactions`     | ✅    | `?page=&limit=`                                                              |
| GET    | `/transactions/:id` | ✅    | —                                                                            |
| POST   | `/transactions`     | ✅    | `{ alamat_pengiriman, items: [{ log_produk_id, kuantitas }], method_bayar }` |

---

## 🗂️ Project Structure

```
FinalTask/
├─ cmd/app/main.go         # Entry point
├─ config/config.go        # Load .env & DB init
├─ internal/
│  ├─ models/              # GORM models
│  ├─ repository/          # DB queries
│  ├─ service/             # Business logic
│  ├─ handler/             # HTTP handlers
│  └─ middleware/          # JWT & role checks
├─ router/router.go        # Route setup
├─ utils/                  # JWT & hashing
├─ .env
├─ go.mod
└─ README.md
```

---

## 📦 Postman Collection

Impor file `Rakamin Evermos Virtual Internship.postman_collection.json` ke Postman, lalu:

1. Set environment variable `base_url = http://localhost:8000/api/v1`
2. Set variable `jwt_token` setelah login
3. Jalankan request sesuai urutan:

   * Auth → Register & Login
   * User → Profile CRUD
   * Store → Get/Update
   * Addresses → CRUD & region lookup
   * Categories (Admin) → CRUD
   * Products → CRUD + upload
   * Transactions → Create & Get

---

## 🙋‍♂️ Q\&A

Jika ada pertanyaan, silakan hubungi saya pada email : thariqhatrama@gmail.com. Terima kasih!

