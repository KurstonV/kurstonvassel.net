package main

import (
	"net/url"

	"kurstonvassel.net/quotebox/pkg/models"
)

type templateData struct {
	Quotes         []*models.Quote
	Quote          *models.Quote
	ErrorsFromForm map[string]string
	Flash          string
	FormData       url.Values
}
