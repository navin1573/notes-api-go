package main

import (
	"fmt"
	"net/http"
	"notes-api/handlers"
)

func main(){
 http.HandleFunc("/notes",handlers.NotesHandler)	
 fmt.Println("Server running on :3000")
 http.ListenAndServe(":3000",nil)
}