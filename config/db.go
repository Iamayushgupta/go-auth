package config

import (
	"database/sql"
	"fmt"
	"os"
)

var DB *sql.DB

func ConnectToDB() {
	var err error
	DB, err = sql.Open("mysql", os.Getenv("LOCAL_MYSQL_URL"))
	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to database!")
	}
	fmt.Println("Database connected!")
}
