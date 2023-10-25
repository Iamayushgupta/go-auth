package config

import (
	"database/sql"
	"fmt"
	"os"
)

var DB *sql.DB

func ConnectToDB() {
	var err error
	DB, err = sql.Open("mysql", os.Getenv("MYSQL_DB_URL")) 
	fmt.Println(os.Getenv("MYSQL_DB_URL"))
	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to the database!")
	}
	fmt.Println("Database connected!")
}