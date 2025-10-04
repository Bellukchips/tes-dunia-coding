package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=3"`
	Description string `json:"description"`
	DueDate     string `json:"due_date" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=pending in-progress completed"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"omitempty,min=3"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status" binding:"omitempty,oneof=pending in-progress completed"`
}

var db *sql.DB

func main() {
	// Initialize database
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := createTable(); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		api.GET("/tasks", getTasks)
		api.GET("/tasks/:id", getTask)
		api.POST("/tasks", createTask)
		api.PUT("/tasks/:id", updateTask)
		api.DELETE("/tasks/:id", deleteTask)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDB() (*sql.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5579")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "belluk")
	dbname := getEnv("DB_NAME", "taskdb")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database")
	return db, nil
}

func createTable() error {
	query := `CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		due_date DATE NOT NULL,
		status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'in-progress', 'completed')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil

}

func getTasks(c *gin.Context) {
	query := `
		SELECT id, title, description, due_date, status, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
		`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	tasks := []Task{}

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, tasks)
}

func getTask(c *gin.Context) {
	id := c.Param("id")

	query := `SELECT id, title, description, due_date, status, created_at, updated_at 
	          FROM tasks WHERE id = $1`

	var task Task
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status, &task.CreatedAt, &task.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}

	c.JSON(http.StatusOK, task)
}
func createTask(c *gin.Context) {
    var req CreateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Println("Bind error:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    dueDate, err := time.Parse("2006-01-02", req.DueDate)
    if err != nil {
        log.Println("Parse date error:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
        return
    }

    _, err = db.Exec(
        "INSERT INTO tasks (title, description, due_date, status) VALUES ($1, $2, $3, $4)",
        req.Title, req.Description, dueDate, req.Status,
    )
    if err != nil {
        log.Printf("DB insert error: %v | title=%s desc=%s due=%v status=%s",
            err, req.Title, req.Description, dueDate, req.Status)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Task created successfully"})
}


func updateTask(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if task exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Validate date format if provided
	if req.DueDate != "" {
		if _, err := time.Parse("2006-01-02", req.DueDate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
	}

	// Build dynamic update query
	query := `UPDATE tasks SET updated_at = CURRENT_TIMESTAMP`
	args := []interface{}{}
	argCount := 1

	if req.Title != "" {
		query += fmt.Sprintf(", title = $%d", argCount)
		args = append(args, req.Title)
		argCount++
	}

	if req.Description != "" {
		query += fmt.Sprintf(", description = $%d", argCount)
		args = append(args, req.Description)
		argCount++
	}

	if req.DueDate != "" {
		query += fmt.Sprintf(", due_date = $%d", argCount)
		args = append(args, req.DueDate)
		argCount++
	}

	if req.Status != "" {
		query += fmt.Sprintf(", status = $%d", argCount)
		args = append(args, req.Status)
		argCount++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, title, description, due_date, status, created_at, updated_at", argCount)
	args = append(args, id)

	var task Task
	err = db.QueryRow(query, args...).Scan(&task.ID, &task.Title, &task.Description,
		&task.DueDate, &task.Status, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
