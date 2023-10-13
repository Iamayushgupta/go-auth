package config

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

func ConnectToDB() {
	var err error
	DB, err = sql.Open("mysql", "root:ayushsql@tcp(127.0.0.1:3306)/authDB")
	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to database!")
	}
	fmt.Println("Database connected!")
}
