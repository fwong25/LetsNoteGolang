package models

import (
	"fmt"
)

type Note struct {
	Id int `db:"id"`
	Title string `db:"title"`
	Content  string `db:"content"`
	Created_date  string `db:"created_date"`
	Last_modified_date string `db:"last_modified_date"`
}

func PrintNote(note Note) {
	fmt.Println("Note ID : ", note.Id)
	fmt.Print(", Title : ", note.Title)
	fmt.Print(", Content : ", note.Content)
	fmt.Print(", Created date : ", note.Created_date)
	fmt.Println(", Last modified date : ", note.Last_modified_date)
}

func PrintNoteList(note_list []Note) {
	for _, note := range note_list {
		PrintNote(note)
	}
}