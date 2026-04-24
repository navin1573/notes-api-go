package main

import (
	"fmt"
	"net/http"
	"notes-api/handlers"
	"notes-api/db"
)

func main(){
	db.Init()
 http.HandleFunc("/signup",handlers.SignupHandler)	
 http.HandleFunc("/login",handlers.LoginHandler)	
 http.HandleFunc("/notes",handlers.NotesHandler)	
 http.HandleFunc("/notes/",handlers.NotesHandler)	
 fmt.Println("Server running on :3000")
 http.ListenAndServe(":3000",nil)
}
