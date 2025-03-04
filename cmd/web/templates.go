package main

import (
	"github.com/konnenl/snippetbox/internal/models"
	"html/template"
	"path/filepath"
	"fmt"
)

type templateData struct{
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
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

		ts, err := template.ParseFiles("./ui/html/base.html")
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