package main

import (
	"net/http"
	"db_mgmt"
	"interfaces"
)
func main() {
	db_mgmt.ConnectDB()

	http.HandleFunc("/", interfaces.Handler)
	http.ListenAndServe(":8000", nil)

	db_mgmt.CloseDB()
}
