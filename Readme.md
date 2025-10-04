# Task Management Application
## 🚀 Features

- Create, Read, Update, Delete (CRUD) tasks
- Change task status (Pending, In Progress, Completed)
- Input validation
- Responsive UI with TailwindCSS
- RESTful API backend

## 📋 Prerequisites

Before running this project, make sure you have installed:

- **Go** (v1.21 or higher) - [Download](https://golang.org/dl/)
- **Node.js** (v18 or higher) and **npm** - [Download](https://nodejs.org/)
- **PostgreSQL** (v13 or higher) - [Download](https://www.postgresql.org/download/)

## 🗄️ Database Setup

1. Create database:
```sql
CREATE DATABASE taskdb;
\q
```

3. The application will automatically create the `tasks` table when you run the backend.

## 🛠️ Backend Setup (Golang)

1. Navigate to backend directory:
```bash
cd go-backend
```

2. Install dependencies:
```bash
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go get github.com/lib/pq
go mod tidy
```

3. Configure database connection:
   - Edit the `.env` file in the backend directory
   - Update the database credentials if needed:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskdb
PORT=8080
```

4. Run the backend:
```bash
go run main.go
```

The backend server will start on `http://localhost:8080`

### Backend API Endpoints

- `GET /api/tasks` - Get all tasks
- `GET /api/tasks/:id` - Get single task
- `POST /api/tasks` - Create new task
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Delete task
- `GET /health` - Health check

## 🎨 Frontend Setup (React)

1. Navigate to frontend directory:
```bash
cd frontend-task
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm start
```

The frontend will start on `http://localhost:5173` and automatically open in your browser.
