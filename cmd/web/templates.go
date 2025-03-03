package main

import "github.com/konnenl/snippetbox/internal/models"

type templateData struct{
	Snippet *models.Snippet
	Snippets []*models.Snippet
}