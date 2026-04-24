package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"notes-api/db"
	"notes-api/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"os"
)
var jwtKey = []byte(os.Getenv("JWT_SECRET"))
func GetIDFromPath(path string) (int, error) {
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		return 0, fmt.Errorf("no id")
	}

	return strconv.Atoi(parts[2])
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Missing fields", 400)
		return
	}

	var id int
	var dbPassword string

	err = db.DB.QueryRow(
		"SELECT id, password FROM users WHERE username = ?",
		user.Username,
	).Scan(&id, &dbPassword)

	if err != nil {
		http.Error(w, "User not found", 401)
		return
	}

	e:=bcrypt.CompareHashAndPassword([]byte(dbPassword),[]byte(user.Password))
	
		if e!=nil {
		http.Error(w, "Invalid password", 401)
		return
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"user_id":id,
		"exp":time.Now().Add(time.Hour*24).Unix(),
	})
	tokenString,err:=token.SignedString(jwtKey)
	if err!=nil {
		http.Error(w,"Token error",500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"token":tokenString,
	})
	
}
func NotesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	userID, err := GetUserFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", 401)
			return
		}

	if r.Method == "GET" {

		id, err := GetIDFromPath(r.URL.Path)

		if err != nil {
			if r.URL.Path == "/notes/" || r.URL.Path == "/notes" {
				rows, err := db.DB.Query("SELECT id, text FROM notes WHERE user_id =?",userID)
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
		err = db.DB.QueryRow("SELECT id, text FROM notes WHERE id = ? AND user_id=?", id,userID).
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
		result, err := db.DB.Exec("INSERT INTO notes(text,user_id) VALUES(?,?)", n.Text,userID,)
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
		_, err = db.DB.Exec("DELETE FROM notes WHERE id = ? AND user_id=?", id,userID)
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
		_, err = db.DB.Exec("UPDATE notes SET text = ? WHERE id = ? AND user_id=?", updated.Text, id,userID)
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


func SignupHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Missing fields", 400)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		http.Error(w, "Error hashing password", 500)
		return
	}

	_, err = db.DB.Exec(
		"INSERT INTO users(username,password) VALUES(?,?)",
		user.Username,
		string(hash),
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			http.Error(w, "Username already exists", 400)
			return
		}
		http.Error(w, "DB error", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created",
	})
}

func GetUserFromToken(r *http.Request)(int,error) {
	authHeader:=r.Header.Get("Authorization")
	if authHeader == "" {
		return 0,fmt.Errorf("missing token")
	}
	parts:=strings.Split(authHeader," ")
	if(len(parts)!=2){
		return 0,fmt.Errorf("invalid Header")
	}
	tokenStr:=parts[1]
	token,err:=jwt.Parse(tokenStr,func(token *jwt.Token) (interface{}, error){
		return jwtKey,nil
	})
	if err !=nil || !token.Valid {
		return 0,fmt.Errorf("invalid token")
	}
	claims,ok:=token.Claims.(jwt.MapClaims)
	if !ok {
		return 0,fmt.Errorf("Invalid claims")
	}
	userID:= int(claims["user_id"].(float64))
	return userID,nil
}