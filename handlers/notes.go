package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"notes-api/db"
	"notes-api/models"
)

func GetIDFromPath(path string) (int, error) {
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		return 0, fmt.Errorf("no id")
	}

	return strconv.Atoi(parts[2])
}

func NotesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {

		id, err := GetIDFromPath(r.URL.Path)

		if err != nil {
			if r.URL.Path == "/notes/" || r.URL.Path == "/notes" {

				rows, err := db.DB.Query("SELECT id, text FROM notes")
				if err != nil {
					http.Error(w, "DB error", 500)
					return
				}
				defer rows.Close()

				var notes []models.Note

				for rows.Next() {
					var n models.Note
					rows.Scan(&n.ID, &n.Text)
					notes = append(notes, n)
				}

				json.NewEncoder(w).Encode(notes)
				return
			}

			http.Error(w, "Invalid ID", 400)
			return
		}

		var n models.Note

		err = db.DB.QueryRow("SELECT id, text FROM notes WHERE id = ?", id).
			Scan(&n.ID, &n.Text)

		if err != nil {
			http.Error(w, "Not found", 404)
			return
		}

		json.NewEncoder(w).Encode(n)
		return
	}

	if r.Method == "POST" {

		var n models.Note

		err := json.NewDecoder(r.Body).Decode(&n)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		result, err := db.DB.Exec("INSERT INTO notes(text) VALUES(?)", n.Text)
		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}

		id, _ := result.LastInsertId()
		n.ID = int(id)

		json.NewEncoder(w).Encode(n)
		return
	}

	if r.Method == "DELETE" {

		id, err := GetIDFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "Invalid ID", 400)
			return
		}

		_, err = db.DB.Exec("DELETE FROM notes WHERE id = ?", id)
		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Deleted Successfully",
		})
		return
	}

	if r.Method == "PUT" {

		id, err := GetIDFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "Invalid ID", 400)
			return
		}

		var updated models.Note

		err = json.NewDecoder(r.Body).Decode(&updated)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		_, err = db.DB.Exec("UPDATE notes SET text = ? WHERE id = ?", updated.Text, id)
		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}

		updated.ID = id

		json.NewEncoder(w).Encode(updated)
		return
	}

	http.Error(w, "Method not allowed", 405)
}
