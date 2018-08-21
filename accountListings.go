package main

import (
	"net/http"
	"html/template"
	"database/sql"
)

var accountListingsTemplate = template.Must(template.ParseFiles("template/root.html", "template/accountListings.html"))
func accountListingsHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := sessionIDCookie(r)
	if err != nil {
		// todo: change to login
		http.Redirect(w, r, "/account/create", http.StatusFound)
		return
	}

	db, err := sqlOpen()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	accountID, err := accountIDBySessionID(db, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			// todo: change to login
			http.Redirect(w, r, "/account/create", http.StatusFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profileID, err := profileIDByAccountID(db, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/profile/create", http.StatusFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query := "SELECT id, title FROM listing WHERE profile_id = ?;"
	stmt, err := db.Prepare(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(profileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var listings []Listing
	for rows.Next() {
		var l Listing
		err = rows.Scan(&l.ID, &l.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		listings = append(listings, l)
	}
	renderTemplate(w, accountListingsTemplate, listings)
}

type Listing struct {
	ID int
	Title string
}