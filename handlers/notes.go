package handlers

import (
	"encoding/json"
	"net/http"
	"notes-api/models"
	"strconv"
)

var Notes = []models.Note{}
var ID = 1

func NotesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		idstr := r.URL.Query().Get("id")

		if idstr != "" {
			id, err := strconv.Atoi(idstr)
			if err != nil {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}

			for _, note := range Notes {
				if note.ID == id {
					json.NewEncoder(w).Encode(note)
					return
				}
			}

			http.Error(w, "Note not Found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(Notes)
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
		idstr := r.URL.Query().Get("id")

		if idstr != "" {
			id, err := strconv.Atoi(idstr)
			if err != nil {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}

			for i, note := range Notes {
				if note.ID == id {

					Notes = append(Notes[:i], Notes[i+1:]...)

					json.NewEncoder(w).Encode(map[string]string{
						"message": "Deleted Successfully",
					})
					return
				}
			}

			http.Error(w, "Note Not Found", http.StatusNotFound)
			return
		}

		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
}