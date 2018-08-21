package main

import (
	"net/http"
	"html/template"
	"database/sql"
	"wed.rentals/ross/website/pkgs/validate"
)

var listingInsertTemplate = template.Must(template.ParseFiles("template/root.html", "template/listingInsert.html"))
func listingInsertHandler(w http.ResponseWriter, r *http.Request) {
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

	var f listingInsertForm
	if r.Method == "POST" {
		f.SetValues(r)
		if f.IsValid() {
			err = f.Exec(db, profileID)
			if err != nil {
				// todo: any other handling here?
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("success"))
			return
		}
	}
	renderTemplate(w, listingInsertTemplate, f)
}

type listingInsertForm struct {
	LocationID,
	CategoryID,
	StyleID,
	ColorID,
	Title,
	TitleError,
	Description,
	DescriptionError,
	Length,
	LengthError,
	Width,
	WidthError,
	Height,
	HeightError,
	PricePerDay,
	PricePerDayError,
	Error string
}

func (f *listingInsertForm) SetValues(r *http.Request) {
	f.LocationID = r.FormValue("location_id")
	f.CategoryID = r.FormValue("category_id")
	f.StyleID = r.FormValue("style_id")
	f.ColorID = r.FormValue("color_id")
	f.Title = r.FormValue("title")
	f.Description = r.FormValue("description")
	f.Length = r.FormValue("length")
	f.Width = r.FormValue("width")
	f.Height = r.FormValue("height")
	f.PricePerDay = r.FormValue("price_per_day")
}

func (f *listingInsertForm) IsValid() bool {
	// todo: change number fields to validate numbers
	isValid := true
	err := validate.LengthBetween(f.Title, 1, 50)
	if err != nil {
		isValid = false
		f.TitleError = err.Error()
	}

	err = validate.LengthBetween(f.Description, 1, 50)
	if err != nil {
		isValid = false
		f.DescriptionError = err.Error()
	}

	err = validate.LengthBetween(f.Length, 1, 5)
	if err != nil {
		isValid = false
		f.LengthError = err.Error()
	}

	err = validate.LengthBetween(f.Width, 1, 5)
	if err != nil {
		isValid = false
		f.WidthError = err.Error()
	}

	err = validate.LengthBetween(f.Height, 1, 5)
	if err != nil {
		isValid = false
		f.HeightError = err.Error()
	}

	err = validate.LengthBetween(f.PricePerDay, 1, 5)
	if err != nil {
		isValid = false
		f.PricePerDayError = err.Error()
	}

	if !isValid {
		f.Error = "check below and try again"
	}

	return isValid
}

func (f listingInsertForm) Exec(db *sql.DB, profileID int) error {
	query := "INSERT INTO listing(profile_id, location_id, category_id, style_id," +
		"color_id, title, description, length, width, height, price_per_day) " +
		"VALUES(?,?,?,?,?,?,?,?,?,?,?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(profileID, f.LocationID, f.CategoryID, f.StyleID, f.ColorID,
		f.Title, f.Description, f.Length, f.Width, f.Height, f.PricePerDay)
	return err
}