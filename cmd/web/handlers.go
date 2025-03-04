package main

import (
	"net/http"
	"fmt"
	"strconv"
	"errors"
	"github.com/konnenl/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil{
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1{
		app.notFound(w)
		return 
	}
	snippet, err := app.snippets.Get(id)
	if err != nil{
		if errors.Is(err, models.ErrNoRecord){
			app.notFound(w)
		}else{
			app.serverError(w, err)
		}
		return 
	}
	
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "home.html", data)
}


func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "Title"
	content := "This is\ncontent"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil{
		app.serverError(w, err)
		return 
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
