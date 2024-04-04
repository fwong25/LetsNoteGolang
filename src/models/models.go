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
	Subnote_exist bool `db:"subnote_exist"`
	Note_level int `db:"note_level"`
	Parent_table_id string `db:"parent_table_id"`
	Parent_note_id string `db:"parent_note_id"`
}

func PrintNote(note Note) {
	fmt.Println("Note ID : ", note.Id)
	fmt.Print(", Title : ", note.Title)
	fmt.Print(", Content : ", note.Content)
	fmt.Print(", Created date : ", note.Created_date)
	fmt.Print(", Last modified date : ", note.Last_modified_date)
	fmt.Print(", Subnote exist : ", note.Subnote_exist)
	fmt.Print(", Note level : ", note.Note_level)
	fmt.Print(", Parent table id : ", note.Parent_table_id)
	fmt.Println(", Parent table id : ", note.Parent_note_id)
}

func PrintNoteList(note_list []Note) {
	for _, note := range note_list {
		PrintNote(note)
	}
}