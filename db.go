package main

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const createNewDb = `CREATE TABLE "users" (
	"id"	INTEGER NOT NULL UNIQUE,
	"tg_id"	INTEGER NOT NULL,
	"lang"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);`

func GetDataBasePtr(dbName string) (*sql.DB, error) {
	file, err := os.Open(dbName)
	file.Close()

	var db *sql.DB
	if err != nil {
		db, err = sql.Open("sqlite3", dbName)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(createNewDb)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return sql.Open("sqlite3", dbName)
}

func CreateNewUser(db *sql.DB, tgId int) error {
	_, err := db.Exec("insert into users (tg_id, lang) values ($1, $2)", tgId, FlagGb)
	return err
}

func SetUserLang(db *sql.DB, tgId int, lang string) error {
	_, err := db.Exec("update users set lang = $1 where tg_id = $2", lang, tgId)
	return err
}

func GetUserLang(db *sql.DB, tgId int) (lang string, err error) {
	if err = db.QueryRow("select lang from users where tg_id = $1", tgId).Scan(&lang); err != nil {
		return "", err
	}
	return lang, nil
}

func CheckUser(db *sql.DB, tgId int) bool {
	if err := db.QueryRow("select id from users where tg_id = $1", tgId).Scan(new(int)); err != nil {
		return false
	}
	return true
}
