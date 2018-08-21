package main

import (
	"net/http"
	"html/template"
)

var listingUpdateImagesTemplate = template.Must(template.ParseFiles("template/root.html", "template/listingUpdateImages.html"))
func listingUpdateImagesHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, listingUpdateImagesTemplate, nil)
}
