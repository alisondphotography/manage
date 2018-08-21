package main

import (
	"net/http"
	"html/template"
)

var listingUpdateImageTemplate = template.Must(template.ParseFiles("template/root.html", "template/listingUpdateImage.html"))
func listingUpdateImageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, listingUpdateImageTemplate, nil)
}
