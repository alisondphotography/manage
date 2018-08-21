package main

import (
	"net/http"
	"html/template"
)

var listingUpdateTemplate = template.Must(template.ParseFiles("template/root.html", "template/listingUpdate.html"))
func listingUpdateHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, listingUpdateTemplate, nil)
}
