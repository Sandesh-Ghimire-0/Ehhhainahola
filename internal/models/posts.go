// for database realated things

// sudo -u postgres psql 
// psql -U newusername -d dbname  -> conn spec db as that user
// psql -U blogadmin -d blogdb -h localhost

package models

import(
	"fmt"
	"log"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"time"
	"os"
)

type Post struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Slug      string    `json:"slug"`
    Content   string    `json:"content"`
    Published bool      `json:"published"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	ID      	int 		
	Post_id    	int 
	Author  	string
	Content 	string
}


type Post_Comment struct{
	Posts		[]Post
	Comments 	[]Comment
	Post_id 		int
}

func Connectdb() (*sql.DB, error){
	err := godotenv.Load()
	connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_SSLMODE"),
    )
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Cannot reach database:", err)
	}
	fmt.Println("Connected to PostgreSQL!")
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS posts (
		id         SERIAL PRIMARY KEY,
		title      TEXT NOT NULL,
		slug       TEXT NOT NULL UNIQUE,
		content    TEXT NOT NULL,
		published  BOOLEAN NOT NULL DEFAULT false,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comments (
		id         SERIAL PRIMARY KEY,
		post_id    INT REFERENCES posts(id) ON DELETE CASCADE,
		author     TEXT NOT NULL,
		content     TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table created!")
	return db, err
}

