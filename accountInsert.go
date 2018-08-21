package main

import (
	"net/http"
	"wed.rentals/ross/website/pkgs/validate"
	"html/template"
	"database/sql"
	"github.com/mattn/go-sqlite3"
)

var accountInsertTemplate = template.Must(template.ParseFiles("template/root.html", "template/accountInsert.html"))

func accountInsertHandler(w http.ResponseWriter, r *http.Request) {
	var f accountInsertForm
	if r.Method == "POST" {
		f.setValues(r)
		if f.checkValidation() {
			db, err := sqlOpen()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer db.Close()
			newAccountID, err := f.exec(db)
			if err != nil {
				if sqliteErr, ok := err.(sqlite3.Error); ok {
					if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
						f.Error = "Check below and try again"
						f.EmailError = "An account with that email already exists"
						renderTemplate(w, accountInsertTemplate, f)
						return
					}
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = setSession(db, w, newAccountID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/profile/create", http.StatusFound)
			return
		}
	}
	renderTemplate(w, accountInsertTemplate, f)
}

type accountInsertForm struct {
	Email,
	EmailError,
	Password,
	PasswordError,
	Error string
}

func (f *accountInsertForm) setValues(r *http.Request) {
	f.Email = r.FormValue("email")
	f.Password = r.FormValue("password")
}

func (f *accountInsertForm) checkValidation() bool {
	ok := true

	err := validate.Email(f.Email)
	if err != nil {
		ok = false
		f.Error = "Check below and try again"
		f.EmailError = err.Error()
	}

	err = validate.Password(f.Password)
	if err != nil {
		ok = false
		f.Error = "Check below and try again"
		f.PasswordError = err.Error()
	}
	return ok
}

func (f *accountInsertForm) exec(db *sql.DB) (int64, error) {
	query := "INSERT INTO account(email, password) VALUES(?,?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(f.Email, f.Password)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}