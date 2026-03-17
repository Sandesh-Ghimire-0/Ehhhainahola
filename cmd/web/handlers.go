package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"waduhek/internal/models"
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/parser"
)

func GetCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("username")
	if err != nil {
    	return nil
	} else {
    	return cookie
	}
}

func (app *Application) serverError(w http.ResponseWriter, r *http.Request, err error) {
    var (
        method = r.Method
        uri    = r.URL.RequestURI()
    )
    app.logger.Error(err.Error(), "method", method, "uri: ", uri)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	cookie := GetCookie(r)

	db := app.DB
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
	}

	rows, err := db.Query(`SELECT id, author, content FROM comments WHERE post_id IS NULL OR post_id = 0`)
	if err != nil {
    	log.Fatal(err)
	}
	defer rows.Close()
	var comments [] models.Comment
	for rows.Next() {
    	var c models.Comment
    	err := rows.Scan(&c.ID, &c.Author, &c.Content)
    	if err != nil {
        	http.Error(w, `{"error" : "oops something happened when fetching all the comments"}`, http.StatusBadRequest)
			return
    	}
		if err := rows.Err(); err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)
	}
	ts, err := template.ParseFiles(files...)
	pnc := models.Post_Comment{
		Posts		: nil,
		Comments	: comments,
		Post_id		: 0,	
		Cookie 		: cookie,
	}
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri: ", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}	
	err = ts.ExecuteTemplate(w, "base", pnc)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *Application) blog(w http.ResponseWriter, r *http.Request) {
	cookie := GetCookie(r)
	// we opted for buffer method to solve the "superfluous response" error
	w.Header().Add("Server", "Go")

	db := app.DB
	
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/blogs.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri: ", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	return
	}
	rows, err := db.Query(`SELECT id, title, slug, content, published, created_at, updated_at FROM posts WHERE published = true`)
	if err := rows.Err(); err != nil {
    	http.Error(w, "something went wrong", http.StatusInternalServerError)
    	return
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
    	var p models.Post
    	err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt)
    	if err != nil {
        	http.Error(w, `{"error" : "oops something happened when fetching all the posts"}`, http.StatusBadRequest)
    	}
    posts = append(posts, p)
	}

	var comments []models.Comment
	var c models.Comment
	rows, err = db.Query(`SELECT id, author, content FROM comments WHERE post_id IS NULL OR post_id = 0`)
	if err := rows.Err(); err != nil {
    	http.Error(w, "something went wrong", http.StatusInternalServerError)
    	return
	}
	defer rows.Close()
	for rows.Next() {
    	err := rows.Scan(&c.ID, &c.Author, &c.Content)
    	if err != nil {
        	http.Error(w, `{"error" : "oops something happened when fetching all the posts"}`, http.StatusBadRequest)
    	}
    comments = append(comments, c)
	}
	pnc := models.Post_Comment{
		Posts		: posts,
		Comments	: comments,
		Post_id		: 0,	
		Cookie		: cookie,	
	}
	err = ts.ExecuteTemplate(w, "base", pnc) 
	
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
	}
}

func (app *Application) indiv_blog_Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	cookie := GetCookie(r)
	db := app.DB
	// Extract the ID from the URL e.g. /post/3 -> "3"
	slug_Str := strings.TrimPrefix(r.URL.Path, "/post/")
	if slug_Str == "" {
		http.Error(w, `{"error": "missing post id"}`, http.StatusBadRequest)
		return
	}
	var p models.Post

	// logic for posts
	row := db.QueryRow("SELECT id, title, content, published, created_at, updated_at FROM posts WHERE slug = $1", slug_Str)
	if db == nil {
    	http.Error(w, "db is nil", http.StatusInternalServerError)
    return
	}
	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, `{"error": "database error"}`, http.StatusInternalServerError)
		return
	}

	// logic for converting md to html

	md := goldmark.New(
          goldmark.WithExtensions(extension.GFM),
          goldmark.WithParserOptions(
              parser.WithAutoHeadingID(),
          ),
          goldmark.WithRendererOptions(
              html.WithHardWraps(),
              html.WithXHTML(),
          ),
      )
	var buf bytes.Buffer
	if err := md.Convert([]byte(p.Content), &buf); err != nil {
		panic(err)
	}

	htmlTemplate := buf.String()


	//logic for post comments
	var comments []models.Comment
	var c models.Comment
	rows, err := db.Query(`SELECT id, author, content FROM comments WHERE post_id = $1`, p.ID)
	if err := rows.Err(); err != nil {
    	http.Error(w, "something went wrong", http.StatusInternalServerError)
    	return
	}
	defer rows.Close()
	for rows.Next() {
    	err := rows.Scan(&c.ID, &c.Author, &c.Content)
    	if err != nil {
        	http.Error(w, `{"error" : "oops something happened when fetching all the posts"}`, http.StatusBadRequest)
    	}
    comments = append(comments, c)
	}

	pnc := models.Post_Comment{
		Posts			: []models.Post{p},
		Comments		: comments,
		Post_id			: p.ID,
		HTMLContent		: template.HTML(htmlTemplate),
		Cookie			: cookie,
	}
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/indiv-blog.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error(), "method: ", r.Method, "uri: ", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	return
	}
	err = ts.ExecuteTemplate(w, "base", pnc) 
	if err != nil {
    	app.logger.Error(err.Error(), "method: ", r.Method, "uri: ", r.URL.RequestURI())
    	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    	return
	}
}

func (app *Application) createComment(w http.ResponseWriter, r *http.Request) {
	db := app.DB
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	author := r.FormValue("author")
	content := r.FormValue("content")
	post_id, err := strconv.Atoi(r.FormValue("post_no"))

	http.SetCookie(w, &http.Cookie{
    Name:     "username",
    Value:    author,
    Path:     "/",
    MaxAge:   7 * 24 * 3600, // 1 week
    HttpOnly: true,
    SameSite: http.SameSiteLaxMode,  // prevent cross site request forgery.(CSRF)
	})

	if (author == "" || content == "" ){
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(
		`INSERT INTO comments (post_id,author, content) VALUES ($1, $2, $3)`,
		post_id, author, content, 
	)
	if err != nil {
		app.logger.Error(err.Error(), "method: ", r.Method, "uri: ", r.URL.RequestURI())
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	log.Println("path:", r.URL.Path)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}