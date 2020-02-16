package models

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "root"
	host     = "127.0.0.1"
	port     = "3306"
	dbName   = "go_example"
)

func InitDB() (*sql.DB, error) {
	path := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", dbName, "?charset=utf8"}, "")

	db, err := sql.Open("mysql", path)
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	if err != nil {
		fmt.Println("open database fail:", err)
		return nil, err
	}
	fmt.Println("connect database success")

	return db, nil
}
