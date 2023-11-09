package model

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
