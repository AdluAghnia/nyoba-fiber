package connection

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitiliazedDB() (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "gofiber",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}
	log.Println("Connected to the database")
	return db, err
}
