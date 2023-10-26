package config

import (
	"database/sql"
	"fmt"
	"os"
)

// An instance of DB
var DB *sql.DB

func ConnectToDB() {
	var err error
	DB, err = sql.Open("mysql", os.Getenv("MYSQL_DB_URL")) 
	
	if err != nil {
		fmt.Println(err)
		//Panic prints the error and terminates the program
		panic("Failed to connect to the database!")
	}
	fmt.Println("Database connected!")
}