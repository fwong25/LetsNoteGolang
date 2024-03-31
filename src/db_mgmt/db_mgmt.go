package db_mgmt

import (
	"fmt"
	"database/sql"
	"models"
	"time"

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

func ConnectDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error when connecting to database", err)
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error when pinging database", err)
	}

	fmt.Println("Successfully connected to database!")
}

func CloseDB() {
	db.Close()
}

func GetNoteList() []models.Note {
	place := models.Note{} // Initialize a User struct to hold retrieved data

	// Execute a SQL query to select "username" and "email" columns from the "users" table
	rows, _ := db.Query("SELECT id, title, content, created_date, last_modified_date FROM letsnote_note")
	note_list := []models.Note{}

	for rows.Next() {
		err := rows.Scan(&place.Id, &place.Title, &place.Content, &place.Created_date, &place.Last_modified_date) // Scan the current row into the "place" variable
		if err != nil {
			fmt.Println("Error when scanning rows", err)
		}
		// fmt.Printf("%#v\n", place) // Log the content of the "place" struct for each row
		note_list = append([]models.Note{place}, note_list...) // ... unpack into note_list[0], note_list[1], ...
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		fmt.Println("Error during note iteration", err)
	}

	return note_list
}

func InsertNote(title string, content string) (note_id int) {
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
	return id
}

func DeleteNote(note_id string) (affected_rows int64) {
	fmt.Println("Delete note with ID: ", note_id)

	sqlStatement := `
		DELETE FROM letsnote_note
		WHERE id = $1;`
		res, err := db.Exec(sqlStatement, note_id)
	
	if err != nil {
		fmt.Println("Error when deleting note", err)
	}

	affected_rows, _ = res.RowsAffected()
	return
}

func UpdateNote(note_id string, title string, content string) (affected_rows int64) {
	fmt.Println("Update note with ID: ", note_id)

	currentTime := time.Now()
	last_modified_date := currentTime.Format("2006-01-02")

	sqlStatement := `
		UPDATE letsnote_note
		SET title = $2, content = $3, last_modified_date = $4
		WHERE id = $1
		RETURNING title, content;`
	res, err := db.Exec(sqlStatement, note_id, title, content, last_modified_date)
	// err = db.QueryRow(sqlStatement, note_id, title, content, last_modified_date).Scan(&title, &content)
	
	if err != nil {
		fmt.Println("Error when updating note", err)
	}

	affected_rows, _ = res.RowsAffected()
	return
}