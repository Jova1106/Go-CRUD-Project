package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

const DB_INFO_FILE = "database/dbinfo.json"

type DBInfo struct {
	Id       string `"json:id"`
	Password string `"json:password"`
}

var db *sql.DB
var err error

/*
InitializeDatabase

1. Read the database info file
2. Parse and store the data in dbInfoMap
3. Verify that the 'user' key exists
4. Open the database
*/
func InitializeDatabase() {
	// 1. Read the database info file
	dbInfoContent, _ := os.ReadFile(DB_INFO_FILE)
	var dbInfoMap map[string]DBInfo

	// 2. Parse and store the data in dbInfoMap
	if dbInfo := json.Unmarshal(dbInfoContent, &dbInfoMap); dbInfo != nil {
		fmt.Printf("%v\n", dbInfo)
	}

	// 3. Verify that the 'user' key exists
	dbInfo, dbInfoExists := dbInfoMap["user"]

	if !dbInfoExists {
		fmt.Println("Error: No user found in dbInfoMap")
		return
	}

	// 4. Open the database
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/dbname", dbInfo.Id, dbInfo.Password))

	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	fmt.Println("Database initialized successfully!")
}
