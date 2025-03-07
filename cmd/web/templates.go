package main

import (
	"github.com/konnenl/snippetbox/internal/models"
	"html/template"
	"path/filepath"
	"fmt"
	"time"
)

type templateData struct{
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
}

func humanDate(t time.Time) string{
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error){
	сache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	fmt.Println(pages)
	if err != nil{
		return nil, err
	}

	for _, page := range pages{
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil{
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil{
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil{
			return nil, err
		}
		сache[name] = ts
	}
	return сache, nil
}