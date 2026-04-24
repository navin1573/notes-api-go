package db
import(
	"database/sql"
	"log"
  _ "github.com/mattn/go-sqlite3"
)
var DB *sql.DB
func Init(){
var err error
DB,err = sql.Open("sqlite3","./notes.db")
if err!=nil{
	log.Fatal(err)
}
createTable()
}

func createTable() {

	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT
	);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS notes(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT,
		user_id INTEGER
	);`)
	if err != nil {
		log.Fatal(err)
	}
}