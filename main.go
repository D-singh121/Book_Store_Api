package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// book model
type Book struct {
	Id          int       `json:"id"`
	TitleName   string    `json:"title_name" binding:"required" unique:"true"`
	AuthorName  string    `json:"author_name" binding:"required"`
	PublishedAt time.Time `json:"published_at"`
}

var db *sql.DB // global db variable

// conneecting to the database
func init() {
	var err error // for global error

	godotenv.Load() // loading env file

	var (
		db_host     = os.Getenv("DB_HOST")
		db_user     = os.Getenv("DB_USER")
		db_password = os.Getenv("DB_PASSWORD")
		db_name     = os.Getenv("DB_NAME")
		db_port     = os.Getenv("DB_PORT")
	)

	connStr := "host=" + db_host + " user=" + db_user + " password=" + db_password + " dbname=" + db_name + " port=" + db_port + " sslmode=disable"

	db, err = sql.Open("postgres", connStr) // assigning to the db global variable
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database not reachable: ", err)
	}
	log.Println("Database connected successfully!")

	CreateTableIfNotExists() // calling the create table function intially after database connection done.
}

// Table creation if not exists after database connection successfull
func CreateTableIfNotExists() {
	query := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title_name VARCHAR(255) NOT NULL UNIQUE,
		author_name VARCHAR(255) NOT NULL,
		published_at DATE NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table: ", err)
	}

	log.Println("Books table checked/created successfully!")
}

// create book
func CreateBook(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input!"})
		return
	}

	// Step 1: Pehle check karenge title already exist karta hai ya nahi
	var exists bool
	checkQuery := "SELECT EXISTS (SELECT 1 FROM books WHERE title_name = $1)"
	err := db.QueryRowContext(ctx, checkQuery, book.TitleName).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check title duplication!"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Book with the same title already exists!"})
		return
	}

	// Step 2: Agar title unique hai, to insert karo
	query := `INSERT INTO books (title_name, author_name, published_at) VALUES ($1, $2, $3) RETURNING id`
	err = db.QueryRowContext(ctx, query, book.TitleName, book.AuthorName, book.PublishedAt).Scan(&book.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book!"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"book": book})
}

// controller for get book
func GetBookById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	id := c.Param("id")
	var book Book

	query := "SELECT * FROM books WHERE id = $1;"
	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(&book.Id, &book.TitleName, &book.AuthorName, &book.PublishedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found!"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get book!"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"book": book}) // return the book on response

}

// controller for get All books
func GetBooks(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	query := "SELECT * FROM books;"          // query banaya
	rows, err := db.QueryContext(ctx, query) // returns rows of data from db call
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Books!"})
		return
	}
	defer rows.Close()

	var books []Book // slice of books

	// scanning the rows
	for rows.Next() {
		var book Book // book object or single book
		err := rows.Scan(&book.Id, &book.TitleName, &book.AuthorName, &book.PublishedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan book!"})
		}
		books = append(books, book)
	}

	// error handling for row iteration
	if rows.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Row iteration error!"})
	}

	c.JSON(http.StatusOK, gin.H{"books": books}) // return the books on response

}

// controller for update book
func UpdateBook(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	id := c.Param("id") // jisko update karna hai uska id

	var book Book // updated book struct
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input!"})
		return
	}

	// Step 1: Pehle check karo kya naya title kisi aur book ka to nahi hai
	var exists bool
	checkQuery := "SELECT EXISTS (SELECT 1 FROM books WHERE title_name = $1 AND id != $2)"
	err := db.QueryRowContext(ctx, checkQuery, book.TitleName, id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check title duplication!"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Another book with the same title already exists!"})
		return
	}

	// Step 2: Agar title unique hai, tab update karo
	query := "UPDATE books SET title_name = $1, author_name = $2, published_at = $3 WHERE id = $4 RETURNING id"
	var updatedID int
	err = db.QueryRowContext(ctx, query, book.TitleName, book.AuthorName, book.PublishedAt, id).Scan(&updatedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully", "updated_id": updatedID})
}

// controller for delete book
func DeleteBook(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	id := c.Param("id") // jisko delete karna hai uska id

	// Step 1: Pehle check karo ki book exist karti hai ya nahi
	var exists bool
	checkQuery := "SELECT EXISTS (SELECT 1 FROM books WHERE id = $1)"
	err := db.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check book existence!"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found!"})
		return
	}

	// Step 2: Agar exist karti hai, to delete karo
	deleteQuery := "DELETE FROM books WHERE id = $1"
	_, err = db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// main func
func main() {
	router := gin.Default()
	godotenv.Load()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hello, world"})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"Health Message": "i am ok "})
	})

	// routes for books with controller define
	router.POST("/api/v1/book", CreateBook)
	router.GET("/api/v1/book/:id", GetBookById)
	router.GET("/api/v1/books", GetBooks)
	router.PUT("/api/v1/book/:id", UpdateBook)
	router.DELETE("api/v1/book/:id", DeleteBook)

	port := os.Getenv("BACKEND_PORT")
	router.Run(":" + port)
}
