package interfaces_v2

import (
	"fmt"
	"net/http"
	"html/template"
	"strconv"
	"path/filepath"
	"db_mgmt"
)

const templatesDirPath = "templates"

var funcMap = map[string]interface {} {
	"Iterate": func(count int) []uint {
		var i uint
		var Items []uint
		for i = 0; i < uint(count); i++ {
			Items = append(Items, i)
		}
		return Items
	},
	"IsSubnote": func(note_level int) bool {
		if note_level == 0 {
			return false
		}
		return true
	},
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case ("/add_new_note"):
		addNewNote(w, r)
	case ("/list_note"):
		listNote(w, r)
	case ("/view_note"):
		viewNote(w, r)
	case ("/insert_note_item_action"):
		insertNoteAction(w, r)
	case ("/delete_note_item"):
		deleteNote(w, r)
	case ("/modify_note_item"):
		modifyNote(w, r)
	case ("/modify_note_item_action"):
		modifyNoteAction(w, r)
	default:
		// fmt.Fprintf(w, "Visit /list_note")
		fmt.Println("Redirect ", r.URL.Path, " to /list_note")
		http.Redirect(w, r, "/list_note", http.StatusSeeOther)
	}
}

func addNewNote(w http.ResponseWriter, r *http.Request) {
	parent_tbl_id := r.FormValue("Parent_tbl_id")
	parent_note_id := r.FormValue("Parent_note_id")

	var fileName = "note_add_new.html"
	t, err := template.New(fileName).Funcs(funcMap).ParseFiles(filepath.Join(templatesDirPath, fileName))
	if err != nil {
		fmt.Println("Error when parsing file", err)
		return
	}

	note_list := db_mgmt.GetNoteListAllSubtable(db_mgmt.Main_tbl_id)

	if parent_tbl_id == "none_none" {
		parent_tbl_id = db_mgmt.Main_tbl_id
	}

	add_note_info := map[string]interface{}{"Parent_tbl_id": parent_tbl_id, "Parent_note_id": parent_note_id, "Note_list": note_list}
	fmt.Println(add_note_info)
	err = t.ExecuteTemplate(w, fileName, add_note_info)
	if err != nil {
		fmt.Println("Error when executing template", err)
		return
	}
}

func listNote(w http.ResponseWriter, r *http.Request) {
	var fileName = "note_list.html"
	t, err := template.New(fileName).Funcs(funcMap).ParseFiles(filepath.Join(templatesDirPath, fileName))
	if err != nil {
		fmt.Println("Error when parsing file", err)
		return
	}

	note_list := db_mgmt.GetNoteListAllSubtable(db_mgmt.Main_tbl_id)

	err = t.ExecuteTemplate(w, fileName, note_list)
	if err != nil {
		fmt.Println("Error when executing template", err)
		return
	}
}

func viewNote(w http.ResponseWriter, r *http.Request) {
	note_id := r.FormValue("note_id")
	tbl_id := r.FormValue("tbl_id")

	note_list := db_mgmt.GetNoteListAllSubtable(db_mgmt.Main_tbl_id)
	if tbl_id == "none_none" {
		tbl_id = db_mgmt.Main_tbl_id
	}
	fmt.Println("Table ID: ", tbl_id)
	selected_note := db_mgmt.GetNote(tbl_id, note_id)

	var fileName = "note_view.html"
	
	t, err := template.New(fileName).Funcs(funcMap).ParseFiles(filepath.Join(templatesDirPath, fileName))
	if err != nil {
		fmt.Println("Error when parsing file", err)
		return
	}

	note_id_int, _ := strconv.Atoi(note_id)

	modify_note_info := map[string]interface{}{"Note_id": note_id_int, "Note_list": note_list, "Selected_note": selected_note}

	err = t.ExecuteTemplate(w, fileName, modify_note_info)
	if err != nil {
		fmt.Println("Error when executing template", err)
		return
	}
}

func insertNoteAction(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("action")
	if action != "Insert" {
		parent_tbl_id := r.FormValue("Parent_tbl_id")
		parent_note_id := r.FormValue("Parent_note_id")

		if parent_tbl_id == "none" {
			http.Redirect(w, r, "/list_note", http.StatusSeeOther)
		} else {
			http.Redirect(w, r,  "/view_note?tbl_id=" + parent_tbl_id + "&note_id=" + parent_note_id, http.StatusSeeOther)
		}
	} else {
		insertNote(w, r)
	}
}

func insertNote(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	parent_tbl_id := r.FormValue("Parent_tbl_id")
	parent_note_id := r.FormValue("Parent_note_id")

	note_id, table_id := db_mgmt.InsertNote(parent_tbl_id, parent_note_id, title, content)
	fmt.Println("New record ID is: ", note_id)
	fmt.Println("Table ID: ", table_id)

	http.Redirect(w, r,  "/view_note?tbl_id=" + table_id + "&note_id=" + strconv.Itoa(note_id), http.StatusSeeOther)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	tbl_id := r.FormValue("tbl_id")
	note_id := r.FormValue("note_id")

	if tbl_id == "none_none" {
		tbl_id = db_mgmt.Main_tbl_id
	}
	note_to_delete := db_mgmt.GetNote(tbl_id, note_id)
	db_mgmt.DeleteNote(tbl_id, note_id)

	if tbl_id == db_mgmt.Main_tbl_id {
		http.Redirect(w, r, "/list_note", http.StatusSeeOther)
	} else {
		parent_tbl_id := note_to_delete.Parent_table_id
		parent_note_id := note_to_delete.Parent_note_id
		http.Redirect(w, r,  "/view_note?tbl_id=" + parent_tbl_id + "&note_id=" + parent_note_id, http.StatusSeeOther)

	}
}

func modifyNote(w http.ResponseWriter, r *http.Request) {
	note_id := r.FormValue("note_id")
	table_id := r.FormValue("tbl_id")
	note_list := db_mgmt.GetNoteListAllSubtable(db_mgmt.Main_tbl_id)

	if table_id == "none_none" {
		table_id = db_mgmt.Main_tbl_id
	}
	note_to_edit := db_mgmt.GetNote(table_id, note_id)

	var fileName = "note_modify.html"
	t, err := template.New(fileName).Funcs(funcMap).ParseFiles(filepath.Join(templatesDirPath, fileName))
	if err != nil {
		fmt.Println("Error when parsing file", err)
		return
	}

	note_id_int, _ := strconv.Atoi(note_id)

	modify_note_info := map[string]interface{}{"Note_id": note_id_int, "Note_list": note_list, "Note_to_edit": note_to_edit}

	err = t.ExecuteTemplate(w, fileName, modify_note_info)
	if err != nil {
		fmt.Println("Error when executing template", err)
		return
	}
}

func modifyNoteConfirm(w http.ResponseWriter, r *http.Request) {
	note_id := r.FormValue("note_id")
	tbl_id := r.FormValue("tbl_id")
	fmt.Println("Confirm modification on note with ID: ", note_id)
	fmt.Println("Table ID: ", tbl_id)
	title := r.FormValue("title")
	content := r.FormValue("content")

	if tbl_id == "none_none" {
		tbl_id = db_mgmt.Main_tbl_id
	}
	db_mgmt.UpdateNote(tbl_id, note_id, title, content)

	http.Redirect(w, r,  "/view_note?tbl_id=" + tbl_id + "&note_id=" + note_id, http.StatusSeeOther)
}

func modifyNoteAction(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("action")
	note_id := r.FormValue("note_id")
	tbl_id := r.FormValue("tbl_id")
	if action != "Update" {
		http.Redirect(w, r, "/view_note?tbl_id=" + tbl_id + "&note_id=" + note_id, http.StatusSeeOther)
	} else {
		modifyNoteConfirm(w, r)
	}
}