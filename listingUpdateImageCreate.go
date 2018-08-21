package main

import (
	"net/http"
	"html/template"
)

var listingUpdateImageCreateTemplate = template.Must(template.ParseFiles("template/root.html", "template/listingUpdateImageCreate.html"))
func listingUpdateImageCreateHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, listingUpdateImageCreateTemplate, nil)
}
