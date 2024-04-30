package db_mgmt

import (
	"fmt"
	"database/sql"
	"models"
	"time"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "12345678"
	dbname   = "LetsNoteDB"
)

var Main_tbl_id string = "note_table_0"

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

func CreateSubnoteTable(parent_note models.Note) {
	cur_note_level := strconv.Itoa(parent_note.Note_level + 1)
	var parent_tbl_id string
	if parent_note.Parent_table_id == "none" {
		parent_tbl_id = Main_tbl_id
	} else {
		parent_tbl_id = parent_note.Parent_table_id + "_" + parent_note.Parent_note_id
	}
	table_name := parent_tbl_id + "_" + strconv.Itoa(parent_note.Id)

	createNoteTableStatement := `CREATE TABLE ` + table_name + ` (
		id SERIAL PRIMARY KEY NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_date TEXT NOT NULL,
		last_modified_date TEXT NOT NULL,
		subnote_exist BOOLEAN NOT NULL DEFAULT FALSE,
		note_level INT NOT NULL DEFAULT ` + cur_note_level + `,
		parent_table_id TEXT NOT NULL DEFAULT '` + parent_tbl_id + `',
		parent_note_id TEXT NOT NULL DEFAULT ` + strconv.Itoa(parent_note.Id) + `
	);`
	_, err := db.Exec(createNoteTableStatement)
	
	if err != nil {
		fmt.Println("Error when creating table", err)
	}
	
	fmt.Println("Created new table: " + table_name)
	return
}

func checkTableExist(parent_tbl_id string, parent_note_id string) (table_name string) {
	if parent_tbl_id == "none" {
		return Main_tbl_id
	}
	parent_note := GetNote(parent_tbl_id, parent_note_id)
	if !parent_note.Subnote_exist {
		CreateSubnoteTable(parent_note)
		setSubnoteExist(parent_tbl_id, parent_note_id, true)
	}
	table_name = parent_tbl_id + "_" + parent_note_id
	return table_name
}

func setSubnoteExist(tbl_id string, note_id string, val bool) (affected_rows int64){
	sqlStatement := `
		UPDATE ` + tbl_id + `
		SET subnote_exist = $1
		WHERE id = $2;`
	res, err := db.Exec(sqlStatement, strconv.FormatBool(val), note_id)	
	if err != nil {
		fmt.Println("Error when updating note", err)
	}

	affected_rows, _ = res.RowsAffected()
	return
}

func GetNoteListAllSubtable(table_id string) []models.Note {
	place := models.Note{} // Initialize a User struct to hold retrieved data

	// Execute a SQL query to select "username" and "email" columns from the "users" table
	rows, _ := db.Query("SELECT id, title, content, created_date, last_modified_date, subnote_exist, note_level, parent_table_id, parent_note_id FROM " + table_id + " ORDER BY last_modified_date")
	note_list := []models.Note{}

	for rows.Next() {
		err := rows.Scan(&place.Id, &place.Title, &place.Content, &place.Created_date, &place.Last_modified_date, &place.Subnote_exist, &place.Note_level, &place.Parent_table_id, &place.Parent_note_id) // Scan the current row into the "place" variable
		if err != nil {
			fmt.Println("Error when scanning rows", err)
		}
		if place.Subnote_exist {
			subnote_list := GetNoteListAllSubtable(table_id + "_" + strconv.Itoa(place.Id))
			note_list = append(subnote_list, note_list...)
		} 
		note_list = append([]models.Note{place}, note_list...) // ... unpack into note_list[0], note_list[1], ...
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		fmt.Println("Error during note iteration", err)
	}

	return note_list
}

func GetNoteListSingleTable(table_name string) []models.Note {
	place := models.Note{} // Initialize a User struct to hold retrieved data

	// Execute a SQL query to select "username" and "email" columns from the "users" table
	rows, _ := db.Query("SELECT id, title, content, created_date, last_modified_date, subnote_exist, note_level, parent_table_id, parent_note_id FROM " + table_name)
	note_list := []models.Note{}

	for rows.Next() {
		err := rows.Scan(&place.Id, &place.Title, &place.Content, &place.Created_date, &place.Last_modified_date, &place.Subnote_exist, &place.Note_level, &place.Parent_table_id, &place.Parent_note_id) // Scan the current row into the "place" variable
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

func GetNote(table_id string, note_id string) models.Note {
	note := models.Note{} // Initialize struct to hold retrieved data

	queryStat := "SELECT id, title, content, created_date, last_modified_date, subnote_exist, note_level, parent_table_id, parent_note_id FROM " + table_id + " WHERE id = " + note_id

	rows, _ := db.Query(queryStat)

	if rows == nil {
		fmt.Println("Can't retrieve note with ID: ", note_id)
	}

	rows.Next()
	err := rows.Scan(&note.Id, &note.Title, &note.Content, &note.Created_date, &note.Last_modified_date, &note.Subnote_exist, &note.Note_level, &note.Parent_table_id, &note.Parent_note_id) // Scan the current row into the "place" variable
	if err != nil {
		fmt.Println("Error when scanning rows", err)
	}

	return note
}

func InsertNote(parent_tbl_id string, parent_note_id string, title string, content string) (note_id int, table_id string) {
	currentTime := time.Now()
	created_date := currentTime.Format("2006-01-02")
	last_modified_date := currentTime.Format("2006-01-02")
	fmt.Println("Insert note:")
	fmt.Println(title, content, created_date, last_modified_date)
	fmt.Println("Parent_tbl_id: "+parent_tbl_id+", parent_note_id: "+parent_note_id)

	table_name := checkTableExist(parent_tbl_id, parent_note_id)
	fmt.Println("table_name: "+table_name)

	sqlStatement := `
		INSERT INTO ` + table_name + ` (title, content, created_date, last_modified_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	
	err = db.QueryRow(sqlStatement, title, content, created_date, last_modified_date).Scan(&note_id)
	if err != nil {
		fmt.Println("Error when inserting note", err)
	}

	table_id = table_name
	return
}

func DeleteNote(table_id string, note_id string) (affected_rows int64) {
	note := GetNote(table_id, note_id)
	fmt.Println("Delete note with ID: ", note_id)

	if note.Subnote_exist {
		DeleteSubNoteTable(table_id + "_" + note_id)
	}

	sqlStatement := `
		DELETE FROM ` + table_id + `
		WHERE id = $1;`
	res, err := db.Exec(sqlStatement, note_id)
	
	if err != nil {
		fmt.Println("Error when deleting note", err)
	}

	affected_rows, _ = res.RowsAffected()

	table_deleted := checkIfNeedDeleteTable(table_id)

	parent_tbl_id := note.Parent_table_id
	parent_note_id := note.Parent_note_id
	if table_deleted && parent_tbl_id != "none" {
		setSubnoteExist(parent_tbl_id, parent_note_id, false)
	}
	return
}

func DeleteSubNoteTable(table_id string){
	fmt.Println("Delete subnote table: ", table_id)

	note_list := GetNoteListSingleTable(table_id)

	for _, note := range note_list {
		if note.Subnote_exist {
			DeleteSubNoteTable(strconv.Itoa(note.Id))
		}
	}

	sqlStatement := `
	DROP TABLE ` + table_id + `;`
	_, err := db.Exec(sqlStatement)
	
	if err != nil {
		fmt.Println("Error when deleting table:", table_id)
		fmt.Println("Error:", err)
	}
	return
}

func checkIfNeedDeleteTable(table_id string) (table_deleted bool) {
	table_deleted = false

	if table_id == Main_tbl_id {
		return table_deleted
	}

	rows_cnt := 0
	queryStat := "SELECT COUNT(*) FROM " + table_id

	rows, _ := db.Query(queryStat)

	if rows == nil {
		fmt.Println("Can't query COUNT of Table ID: ", table_id)
	}

	rows.Next()
	err := rows.Scan(&rows_cnt)
	if err != nil {
		fmt.Println("Error when scanning rows in checkIfNeedDeleteTable", err)
	}

	if rows_cnt == 0 {
		sqlStatement := `
		DROP TABLE ` + table_id + `;`
		_, err := db.Exec(sqlStatement)
		
		if err != nil {
			fmt.Println("Error when deleting table:", table_id)
			fmt.Println("Error:", err)
		}
		table_deleted = true
	}
	return table_deleted
}

func UpdateNote(table_id string, note_id string, title string, content string) (affected_rows int64) {
	fmt.Println("Update note with table ID: " + table_id + ", note ID: " + note_id)

	currentTime := time.Now()
	last_modified_date := currentTime.Format("2006-01-02")

	sqlStatement := `
		UPDATE ` + table_id + `
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