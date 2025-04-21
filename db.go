package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "comments.db")
	if err != nil {
		log.Fatal("DB接続失敗:", err)
	}

	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Fatal("スキーマ読み込み失敗:", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatal("スキーマ実行失敗:", err)
	}
}

func SaveComment(source, channel, username, message string) {
	stmt, err := db.Prepare("INSERT INTO comments (source, channel, username, message) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("INSERT 準備失敗:", err)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(source, channel, username, message); err != nil {
		log.Println("コメント保存失敗:", err)
	}
}
