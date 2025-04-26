# 📚 Book Store API (Golang + Gin + PostgreSQL)

This project is a simple RESTful API for managing books.  
Built with **Go**, **Gin Web Framework**, and **PostgreSQL**.

---

## 🚀 Features
- Create a new book 📖
- Get a book by ID 🔍
- Get all books 🗂
- Update a book 📝
- Delete a book 🗑
- Input validation ✅
- Duplicate title prevention 🔒

---

## 📦 Project Structure
```
/Book_Store_Api
├── main.go
├── go.mod
├── go.sum
├── .env
└── README.md
```

---

## ⚙️ Setup Instructions

1. **Clone the repository**
```bash
git clone https://github.com/D-singh121/Book_Store_Api.git
cd Book_Store_Api
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Set environment variables**
Create a `.env` file:
```
DB_HOST=localhost
DB_USER=your_postgres_user
DB_PASSWORD=your_postgres_password
DB_NAME=your_database_name
DB_PORT=5432
```

4. **Create the `books` table**
```sql
CREATE TABLE IF NOT EXISTS books (
  id SERIAL PRIMARY KEY,
  title_name VARCHAR(255) UNIQUE NOT NULL,
  author_name VARCHAR(255) NOT NULL,
  published_at TIMESTAMP
);
```

5. **Run the application**
```bash
go run main.go
```
Server will start at: `http://localhost:8000`

---

## 📜 API Endpoints

All endpoints are prefixed with `/api/v1`

| Method | Endpoint                  | Description             |
|:------:|:----------------------------|:-------------------------|
| GET    | `/api/v1/`                  | Welcome Route            |
| GET    | `/api/v1/health`             | Health Check             |
| POST   | `/api/v1/book`               | Create a new book        |
| GET    | `/api/v1/book/:id`           | Get a single book by ID  |
| GET    | `/api/v1/books`              | Get all books            |
| PUT    | `/api/v1/book/:id`           | Update a book            |
| DELETE | `/api/v1/book/:id`           | Delete a book            |

---

## 🛡 Validations and Business Rules
- **Title must be unique** (No two books can have the same title).
- **JSON body validation** for required fields.
- **Context timeout** of 2 seconds for every DB operation to avoid hanging requests.

---

## 🛠 Tech Stack
- **Golang** (Go 1.22+)
- **Gin** (HTTP web framework)
- **PostgreSQL** (Database)
- **GoDotEnv** (Load .env files)

---

## ❤️ Special Notes
- Always run the server in `release mode` for production.
- Never trust all proxies in production (`gin.SetTrustedProxies` should be properly set).

---

## 🙋‍♂️ Author
**Devesh Singh**  
Built with 💙 and Go.

GitHub Repository: [Book_Store_Api](https://github.com/D-singh121/Book_Store_Api.git)

---

## 📸 Example Payloads

### Create Book (POST `/api/v1/book`)
```json
{
  "title_name": "The Go Programming Language",
  "author_name": "Alan A. A. Donovan",
  "published_at": "2025-04-26T12:00:00Z"
}
```

### Update Book (PUT `/api/v1/book/:id`)
```json
{
  "title_name": "The Advanced Go Programming",
  "author_name": "Brian Kernighan",
  "published_at": "2025-04-30T00:00:00Z"
}
```

---
