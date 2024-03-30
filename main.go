package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"
	"time"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345678"
	dbname   = "LetsNoteDB"
)

var db *sql.DB
var err error

// CREATE TABLE letsnote_note (
// 	id SERIAL PRIMARY KEY NOT NULL,
// 	title TEXT NOT NULL,
// 	content TEXT NOT NULL,
// 	created_date TEXT NOT NULL,
// 	last_modified_date TEXT NOT NULL
//   );

type Note struct {
	Id int `db:"id"`
	Title string `db:"title"`
	Content  string `db:"content"`
	Created_date  string `db:"created_date"`
	Last_modified_date string `db:"last_modified_date"`
}

type ModifyNoteInfo struct {
	Note_list []Note
	Note_id int
}

func listNote(w http.ResponseWriter, r *http.Request) {
	var fileName = "note_list.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("Error when parsing file", err)
		return
	}

	place := Note{} // Initialize a User struct to hold retrieved data

	// Execute a SQL query to select "username" and "email" columns from the "users" table
	rows, _ := db.Query("SELECT id, title, content, created_date, last_modified_date FROM letsnote_note")
	note_list := []Note{}

	for rows.Next() {
		err := rows.Scan(&place.Id, &place.Title, &place.Content, &place.Created_date, &place.Last_modified_date) // Scan the current row into the "place" variable
		if err != nil {
			fmt.Println("Error when scanning rows", err)
		}
		fmt.Printf("%#v\n", place) // Log the content of the "place" struct for each row
		note_list = append([]Note{place}, note_list...) // ... unpack into note_list[0], note_list[1], ...
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		fmt.Println("Error during note iteration", err)
	}

	err = t.ExecuteTemplate(w, fileName, note_list)
	if err != nil {
		fmt.Println("Error when executing template", err)
		return
	}
}

func insertNote(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")

	currentTime := time.Now()
	created_date := currentTime.Format("2006-01-02")
	last_modified_date := currentTime.Format("2006-01-02")
	fmt.Println("Insert note:")
	fmt.Println(title, content, created_date, last_modified_date)

	sqlStatement := `
		INSERT INTO letsnote_note (title, content, created_date, last_modified_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	
	id := 0
	err = db.QueryRow(sqlStatement, title, content, created_date, last_modified_date).Scan(&id)
	if err != nil {
		fmt.Println("Error when inserting note", err)
	}
	fmt.Println("New record ID is:", id)

	http.Redirect(w, r, "/list_note", http.StatusSeeOther)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	note_id := r.FormValue("note_id")
	fmt.Println("Delete note with ID: ", note_id)

	sqlStatement := `
		DELETE FROM letsnote_note
		WHERE id = $1;`
		_, err := db.Exec(sqlStatement, note_id)
	if err != nil {
		fmt.Println("Error when deleting note", err)
	}
	http.Redirect(w, r, "/list_note", http.StatusSeeOther)
}

func modifyNote(w http.ResponseWriter, r *http.Request) {
	note_id, _ := strconv.Atoi(r.FormValue("note_id"))
	fmt.Println("Modify note with ID: ", note_id)

	var fileName = "note_modify.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("Error when parsing file", err)
		return
	}

	place := Note{} // Initialize a User struct to hold retrieved data

	// Execute a SQL query to select "username" and "email" columns from the "users" table
	rows, _ := db.Query("SELECT id, title, content, created_date, last_modified_date FROM letsnote_note")
	note_list := []Note{}

	for rows.Next() {
		err := rows.Scan(&place.Id, &place.Title, &place.Content, &place.Created_date, &place.Last_modified_date) // Scan the current row into the "place" variable
		if err != nil {
			fmt.Println("Error when scanning rows", err)
		}
		fmt.Printf("%#v\n", place) // Log the content of the "place" struct for each row
		note_list = append([]Note{place}, note_list...) // ... unpack into note_list[0], note_list[1], ...
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		fmt.Println("Error during note iteration", err)
	}

	// modify_note_info := ModifyNoteInfo{Note_list: note_list, Note_id: note_id} // Initialize a User struct to hold retrieved data
	modify_note_info := map[string]interface{}{"Note_id": note_id, "Note_list": note_list}

	err = t.ExecuteTemplate(w, fileName, modify_note_info)
	if err != nil {
		fmt.Println("Error when executing template", err)
		return
	}
}

func modifyNoteAction(w http.ResponseWriter, r *http.Request) {
	
	action := r.FormValue("action")
	if action != "Update" {
		http.Redirect(w, r, "/list_note", http.StatusSeeOther)
	}

	note_id, _ := strconv.Atoi(r.FormValue("note_id"))
	fmt.Println("Confirm modification on note with ID: ", note_id)

	title := r.FormValue("title")
	content := r.FormValue("content")
	currentTime := time.Now()
	last_modified_date := currentTime.Format("2006-01-02")

	sqlStatement := `
		UPDATE letsnote_note
		SET title = $2, content = $3, last_modified_date = $4
		WHERE id = $1
		RETURNING title, content;`
	// _, err = db.Exec(sqlStatement, id, "Updated title 5", "Updated content 5")
	err = db.QueryRow(sqlStatement, note_id, title, content, last_modified_date).Scan(&title, &content)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/list_note", http.StatusSeeOther)
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case ("/list_note"):
		listNote(w, r)
	case ("/insert_note_item"):
		insertNote(w, r)
	case ("/delete_note_item"):
		deleteNote(w, r)
	case ("/modify_note_item"):
		modifyNote(w, r)
	case ("/modify_note_item_action"):
		modifyNoteAction(w, r)
	default:
		fmt.Fprintf(w, "Visit /list_note")
	}
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error when connecting to database", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error when pinging database", err)
	}

	fmt.Println("Successfully connected!")

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8000", nil)
	// http.ListenAndServeTLS(":8000", "cert.pem", "key.pem", nil)
	// go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=localhost
}
