package main

import (
	"net/http"
	"html/template"
	"log"
	"database/sql"
	"math/rand"
	"os"
	"errors"
	"strconv"
)

func main() {
	/*
	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/account/create", accountInsertHandler) // update
	http.HandleFunc("/account", accountHandler)

	http.HandleFunc("/profile/create", profileInsertHandler) // update
	http.HandleFunc("/profile", profileHandler)
	*/

	// http.HandleFunc("/listing/create", listingInsertHandler) // update // i want create and update page to be same if possible, draft?
	http.HandleFunc("/", rerouter)
	http.HandleFunc("/listing/update", listingUpdateHandler)
	http.HandleFunc("/listing/update/content", listingUpdateContentHandler)
	http.HandleFunc("/listing/update/images", listingUpdateImagesHandler)
	http.HandleFunc("/listing/update/image/add", listingUpdateImageCreateHandler)
	http.HandleFunc("/listing/update/image", listingUpdateImageHandler)

	// http.HandleFunc("/listing", listingHandler)

	http.HandleFunc("/account/listings", accountListingsHandler)
	// http.HandleFunc("/listings", indexHandler)
	// http.HandleFunc("/faves", indexHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func rerouter(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/listing/update", http.StatusFound)
}

func sqlOpen() (*sql.DB, error) {
	return sql.Open("sqlite3", "wed.db?_fk=1") // todo: i don't always need fk set
}

func setSession(db *sql.DB, w http.ResponseWriter, accountID int64) error {
	// todo: more random session id
	// todo: retry after unique error on session id
	// todo: remember me
	query := "INSERT INTO session(id, account_id) VALUES(?,?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	newSessionID := rand.Intn(10000000000000000)
	_, err = stmt.Exec(newSessionID, accountID)

	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: strconv.Itoa(newSessionID), Path: "/"})
	return err
}

func sessionIDCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", errors.New("empty")
	}
	return cookie.Value, nil
}

func accountIDBySessionID(db *sql.DB, sessionID string) (int, error) {
	query := "SELECT account_id FROM session WHERE id = ?;"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(sessionID)
	var accountID int
	err = row.Scan(&accountID)
	return accountID, err
}

func profileIDByAccountID(db *sql.DB, accountID int) (int, error) {
	query := "SELECT id FROM profile WHERE account_id = ?;"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(accountID)
	var profileID int
	err = row.Scan(&profileID)
	return profileID, err
}

var indexTemplate = template.Must(template.ParseFiles("template/root.html", "template/index.html"))
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, indexTemplate, nil)
}

var accountTemplate = template.Must(template.ParseFiles("template/root.html", "template/account.html"))
func accountHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, accountTemplate, nil)
}

var profileTemplate = template.Must(template.ParseFiles("template/root.html", "template/profile.html"))
func profileHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, profileTemplate, nil)
}

var listingTemplate = template.Must(template.ParseFiles("template/root.html", "template/listing.html"))
func listingHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, listingTemplate, nil)
}

func renderTemplate(w http.ResponseWriter, t *template.Template, data interface{}) {
	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	os.Remove("wed.db")
	queries := []string{
		"CREATE TABLE account(" +
			"id INTEGER NOT NULL PRIMARY KEY," +
			"email TEXT NOT NULL UNIQUE," +
			"password TEXT NOT NULL" +
		");",
		"CREATE TABLE session(" +
			"id INTEGER NOT NULL UNIQUE," +
			"account_id INTEGER NOT NULL REFERENCES account(id)" +
		");",
		"CREATE TABLE profile(" +
			"id INTEGER NOT NULL PRIMARY KEY," +
			"account_id INTEGER NOT NULL UNIQUE REFERENCES account(id)," +
			"first_name TEXT NOT NULL," +
			"last_name TEXT NOT NULL," +
			"business_name TEXT NOT NULL," +
			"phone INTEGER NOT NULL," +
			"website TEXT NOT NULL" +
		");",
		"CREATE TABLE listing(" +
			"id INTEGER NOT NULL PRIMARY KEY," +
			"profile_id INTEGER NOT NULL REFERENCES profile(id)," +
			"location_id INTEGER NOT NULL," +
			"category_id INTEGER NOT NULL," +
			"style_id INTEGER NOT NULL," +
			"color_id INTEGER NOT NULL," +
			"title TEXT NOT NULL," +
			"description TEXT NOT NULL," +
			"length INTEGER NOT NULL," +
			"width INTEGER NOT NULL," +
			"height INTEGER NOT NULL," +
			"price_per_day REAL NOT NULL" +
		");",
	}
	db, err := sqlOpen()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err, query)
		}
	}
}