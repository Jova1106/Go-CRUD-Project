package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const DB_INFO_FILE = "database/dbinfo.json"

var db *sql.DB

// InitializeDatabase
//
// 1. Read the database info env file
// 2. Open the database
func InitializeDatabase() {
	// 1. Read the database info env file
	dbInfoMap, err := godotenv.Read("database/dbinfo.env")

	if err != nil {
		fmt.Println("Error reading dbinfo.env:", err)
		return
	}

	USERNAME := dbInfoMap["USERNAME"]
	PASSWORD := dbInfoMap["PASSWORD"]

	dbOpenStr := fmt.Sprintf("%s:%s@/dbname", USERNAME, PASSWORD)

	// 2. Open the database
	db, err = sql.Open("mysql", dbOpenStr)

	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	fmt.Printf("Database initialized successfully!\nLogged in as user '%s.'\n", USERNAME)
}
