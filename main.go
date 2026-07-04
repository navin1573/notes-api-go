package main

import (
	"fmt"
	"net/http"
	"notes-api/handlers"
	"notes-api/db"
	"os"
	"github.com/rs/cors"
)

func main(){
	db.Init()
 http.HandleFunc("/signup",handlers.SignupHandler)	
 http.HandleFunc("/login",handlers.LoginHandler)	
 http.HandleFunc("/notes",handlers.NotesHandler)	
 http.HandleFunc("/notes/",handlers.NotesHandler)	
 c:=cors.New(cors.Options{
	 AllowedOrigins: []string{"http://localhost:5173"},
	 AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	 AllowedHeaders: []string{"*"},
 })
 port := os.Getenv("PORT")
	if port == "" {
		port = "3000" 
	}
 fmt.Println("Server running on :"+port)
 http.ListenAndServe(":"+port,c.Handler(http.DefaultServeMux))
}
