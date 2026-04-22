package handlers

import (
	"strings"
	"fmt"
	"encoding/json"
	"net/http"
	"notes-api/models"
	"strconv"
)

var Notes = []models.Note{}
var ID = 1

func GetIDFromPath(path string) (int,error){
	parts:= strings.Split(path,"/")
	if len(parts)< 3 {
		return 0,fmt.Errorf("no id")
	}
	return strconv.Atoi(parts[2])
}

func NotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type","application/json")
		id, err:=GetIDFromPath(r.URL.Path)
		if err!=nil{
			if r.URL.Path=="/notes/" || r.URL.Path == "/notes" {
				json.NewEncoder(w).Encode(Notes)
				return
			}
			http.Error(w,"Invalid ID",http.StatusBadRequest)
			return
		}
		for _,note:=range Notes {
			if note.ID == id {
				json.NewEncoder(w).Encode(note)
				return
			}
		}
		http.Error(w,"Note not Found",http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		var n models.Note

		err := json.NewDecoder(r.Body).Decode(&n)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		n.ID = ID
		ID++

		Notes = append(Notes, n)

		json.NewEncoder(w).Encode(n)
		return
	}

	if r.Method == "DELETE" {
		w.Header().Set("Conetnt-Type","application/json")
		id,err:=GetIDFromPath(r.URL.Path)
		if err!=nil {
			http.Error(w,"Invalid ID",http.StatusBadRequest)
			return
		}
		for i,note:=range Notes {
			if note.ID == id {
				Notes = append(Notes[:i],Notes[i+1:]...)
				json.NewEncoder(w).Encode(map[string]string{
					"message":"Deleted Succesfully",
				})
				return
			}
		}
		http.Error(w,"No Note Found",http.StatusNotFound)
		return
}
if r.Method == "PUT" {

	w.Header().Set("Content-Type", "application/json")

	id, err := GetIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updated models.Note

	err = json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for i, note := range Notes {
		if note.ID == id {

			Notes[i].Text = updated.Text

			json.NewEncoder(w).Encode(Notes[i])
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
	return
}
}
