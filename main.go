package main

import (
	"fmt"
	"net/http"
	"notes-api/handlers"
	"notes-api/db"
	"os"
)

func main(){
	db.Init()
 http.HandleFunc("/signup",handlers.SignupHandler)	
 http.HandleFunc("/login",handlers.LoginHandler)	
 http.HandleFunc("/notes",handlers.NotesHandler)	
 http.HandleFunc("/notes/",handlers.NotesHandler)	
 port := os.Getenv("PORT")
	if port == "" {
		port = "3000" 
	}
 fmt.Println("Server running on :"+port)
 http.ListenAndServe(":"+port,nil)
}
