package interfaces_v2

import (
	"fmt"
	"net/http"
	"html/template"
	"strconv"
	"path/filepath"
	"db_mgmt"
)

const templatesDirPath = "templates_v2"

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
	var fileName = "note_add_new.html"
	t, err := template.ParseFiles(filepath.Join(templatesDirPath, fileName))
	if err != nil {
		fmt.Println("Error when parsing file", err)
		return
	}

	note_list := db_mgmt.GetNoteList()

	err = t.ExecuteTemplate(w, fileName, note_list)
	if err != nil {
		fmt.Println("Error when executing template", err)
		return
	}
}

func listNote(w http.ResponseWriter, r *http.Request) {
	var fileName = "note_list.html"
	t, err := template.ParseFiles(filepath.Join(templatesDirPath, fileName))
	if err != nil {
		fmt.Println("Error when parsing file", err)
		return
	}

	note_list := db_mgmt.GetNoteList()

	err = t.ExecuteTemplate(w, fileName, note_list)
	if err != nil {
		fmt.Println("Error when executing template", err)
		return
	}
}

func viewNote(w http.ResponseWriter, r *http.Request) {
	note_id := r.FormValue("note_id")

	note_list := db_mgmt.GetNoteList()
	selected_note := db_mgmt.GetNote(note_id)

	var fileName = "note_view.html"
	t, err := template.ParseFiles(filepath.Join(templatesDirPath, fileName))
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
		http.Redirect(w, r, "/list_note", http.StatusSeeOther)
	} else {
		insertNote(w, r)
	}
}

func insertNote(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")

	note_id := db_mgmt.InsertNote(title, content)
	fmt.Println("New record ID is:", note_id)

	http.Redirect(w, r, "/view_note?note_id=" + strconv.Itoa(note_id), http.StatusSeeOther)

}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	note_id := r.FormValue("note_id")
	db_mgmt.DeleteNote(note_id)
	http.Redirect(w, r, "/list_note", http.StatusSeeOther)
}

func modifyNote(w http.ResponseWriter, r *http.Request) {
	note_id := r.FormValue("note_id")
	note_list := db_mgmt.GetNoteList()
	note_to_edit := db_mgmt.GetNote(note_id)

	var fileName = "note_modify.html"
	t, err := template.ParseFiles(filepath.Join(templatesDirPath, fileName))
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
	fmt.Println("Confirm modification on note with ID: ", note_id)
	title := r.FormValue("title")
	content := r.FormValue("content")

	db_mgmt.UpdateNote(note_id, title, content)

	http.Redirect(w, r, "/view_note?note_id=" + note_id, http.StatusSeeOther)
}

func modifyNoteAction(w http.ResponseWriter, r *http.Request) {
	
	action := r.FormValue("action")
	if action != "Update" {
		http.Redirect(w, r, "/view_note?note_id=" + r.FormValue("note_id"), http.StatusSeeOther)
	} else {
		modifyNoteConfirm(w, r)
	}
}