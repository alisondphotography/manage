package main

import (
	"net/http"
	"html/template"
)

var listingUpdateContentTemplate = template.Must(template.ParseFiles("template/root.html", "template/listingUpdateContent.html"))
func listingUpdateContentHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, listingUpdateContentTemplate, nil)
}
