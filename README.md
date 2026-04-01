# 🚀 Task Tracker API

A scalable, production-ready **Task Management Backend** built with **Go and MongoDB**, designed for real-time collaboration, team workflows, and modern application architectures.

---

## ✨ Features

* 🔐 JWT Authentication & Authorization
* ⚡ Real-time updates via Server-Sent Events (SSE)
* 🧩 Clean layered architecture (Handler → Service → Repository)
* 🛡️ Rate limiting & request logging middleware
* 👥 Team & collaboration support
* 📊 Boards, Columns, and Tasks management
* 🔄 MongoDB transactions for data consistency
* 📡 RESTful API design

---

## 🏗️ Tech Stack

| Technology | Purpose             |
| ---------- | ------------------- |
| Go (1.21+) | Backend language    |
| MongoDB    | NoSQL database      |
| JWT        | Authentication      |
| SSE        | Real-time streaming |
| bcrypt     | Password hashing    |

---

## 📁 Project Structure

```
task-tracker-api/
│
├── config/        # Configuration & DB setup
├── models/        # Data models
├── repository/    # Database layer (MongoDB queries)
├── services/      # Business logic
├── handlers/      # HTTP layer
├── middleware/    # Auth, logging, rate limiting
├── sse/           # Real-time event system
├── routes/        # API routes
├── utils/         # Helpers (errors, JWT, responses)
│
└── main.go        # Application entry point
```

---

## 🧠 Architecture

The API follows a strict **layered architecture**:

```
Request → Middleware → Handler → Service → Repository → Database
```

### Key Principles

* Separation of concerns
* Dependency injection (no globals)
* Business logic isolated in services
* Database logic isolated in repositories

---

## 🔐 Authentication

Uses **JWT (JSON Web Tokens)** for secure access.

### Flow

1. User registers/logs in
2. Server returns JWT
3. Client sends token in header:

   ```
   Authorization: Bearer <token>
   ```
4. Middleware validates and injects user context

---

## 📡 Real-Time Updates (SSE)

The API supports **live updates without polling** using Server-Sent Events.

### Endpoint

```
GET /stream?board_id=<id>
```

### Events

* `task:created`
* `task:updated`
* `task:deleted`
* `column:updated`
* `board:updated`
* `member:added`
* `member:removed`

---

## 📚 API Endpoints

### Auth (Public)

```
POST   /auth/register
POST   /auth/login
```

### Tasks (Protected)

```
POST   /tasks
GET    /tasks?board_id=:id
GET    /tasks/:id
PATCH  /tasks/:id
DELETE /tasks/:id
```

### Boards

```
POST   /boards
GET    /boards
GET    /boards/:id
PATCH  /boards/:id
DELETE /boards/:id
```

### Teams

```
POST   /teams
GET    /teams
GET    /teams/:id
PATCH  /teams/:id
DELETE /teams/:id
POST   /teams/:id/members
DELETE /teams/:id/members/:memberID
```

### Columns

```
PATCH /columns/:id
```

---

## ⚙️ Environment Variables

Create a `.env` file:

```env
# Server
APP_PORT=8080
APP_ENV=development

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=task_tracker

# JWT
JWT_SECRET=your-super-secret-key
```

Generate a secure secret:

```bash
openssl rand -hex 32
```

---

## 🗄️ Database Design

### Collections

* `users`
* `tasks`
* `boards`
* `columns`
* `teams`
* `members`

### Highlights

* Unique index on `users.email`
* `$lookup` aggregation for dynamic joins
* Transactions for multi-collection operations

---

## 🛡️ Middleware

* **Rate Limiter** → Prevents abuse (2 req/sec, burst 5)
* **Logger** → Structured request logging
* **Auth Middleware** → JWT validation

---

## ❗ Error Handling

Standardized error responses:

```json
{
  "error": "message"
}
```

Validation errors:

```json
{
  "error": "validation failed",
  "fields": {
    "title": "required"
  }
}
```

---

## ▶️ Running the Project

### Prerequisites

* Go 1.21+
* MongoDB
* Git

### Setup

```bash
git clone https://github.com/your-username/task-tracker-api
cd task-tracker-api

go mod tidy

cp .env.example .env
# configure environment variables

go run main.go
```

Server runs at:

```
http://localhost:8080
```

---

## 📖 Swagger Docs

Generate docs:

```bash
swag init
```

Open:

```
http://localhost:8080/swagger/index.html
```

---

## 🎯 Design Highlights

* Clean architecture aligned with backend best practices
* Real-time system without WebSockets (SSE-based)
* Scalable MongoDB patterns (aggregation + transactions)
* Production-grade middleware stack

---

## 📌 Use Cases

* Task management platforms (like Trello/Jira)
* Team collaboration tools
* Real-time dashboards
* Workflow automation systems

---

## 👤 Author

**Kolawole Omopariola**
Frontend & Fullstack Engineer

