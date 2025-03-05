package main

import (
	"net/http"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler{
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/", http.HandlerFunc(app.home))
	mux.Handle("/snippet/view", http.HandlerFunc(app.snippetView))
	mux.Handle("/snippet/create", http.HandlerFunc(app.snippetCreate))

	standart := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standart.Then(mux)
}