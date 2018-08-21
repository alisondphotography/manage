package main

import (
	"net/http"
	"html/template"
	"wed.rentals/ross/website/pkgs/validate"
	"database/sql"
)

var profileInsertTemplate = template.Must(template.ParseFiles("template/root.html", "template/profileInsert.html"))

func profileInsertHandler(w http.ResponseWriter, r *http.Request) {
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

	var f profileInsertForm
	if r.Method == "POST" {
		f.setValues(r)
		if f.isValid() {
			err = f.Exec(db, accountID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/listing/create", http.StatusFound)
			return
		}
	}
	renderTemplate(w, profileInsertTemplate, f)
}

type profileInsertForm struct {
	FirstName,
	FirstNameError,
	LastName,
	LastNameError,
	BusinessName,
	BusinessNameError,
	Phone,
	PhoneError,
	Website,
	WebsiteError,
	Error string
}

func (f *profileInsertForm) setValues(r *http.Request) {
	f.FirstName = r.FormValue("first_name")
	f.LastName = r.FormValue("last_name")
	f.BusinessName = r.FormValue("business_name")
	f.Phone = r.FormValue("phone")
	f.Website = r.FormValue("website")
}

func (f *profileInsertForm) isValid() bool {
	isValid := true

	err := validate.LengthBetween(f.FirstName, 1, 20)
	if err != nil {
		isValid = false
		f.FirstNameError = err.Error()
	}

	err = validate.LengthBetween(f.LastName, 1, 20)
	if err != nil {
		isValid = false
		f.LastNameError = err.Error()
	}

	err = validate.LengthBetween(f.BusinessName, 1, 20)
	if err != nil {
		isValid = false
		f.BusinessNameError = err.Error()
	}

	// todo: validate phone
	err = validate.LengthBetween(f.Phone, 1, 20)
	if err != nil {
		isValid = false
		f.PhoneError = err.Error()
	}

	// todo: validate url
	err = validate.LengthBetween(f.Website, 1, 20)
	if err != nil {
		isValid = false
		f.WebsiteError = err.Error()
	}

	if !isValid {
		f.Error = "check below and try again"
	}

	return isValid
}

func (f profileInsertForm) Exec(db *sql.DB, accountID int) error {
	query := "INSERT INTO profile(account_id, first_name, last_name, business_name, phone, website)" +
		" VALUES(?,?,?,?,?,?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(accountID, f.FirstName, f.LastName, f.BusinessName, f.Phone, f.Website)
	return err
}