package main

import (
	"net/http"
)

func (app *Application) Routes() *http.ServeMux{

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /blogs/{$}", app.blog)
	mux.HandleFunc("GET /post/", app.indiv_blog_Handler)
	mux.HandleFunc("POST /comment", app.createComment)
	return mux
}