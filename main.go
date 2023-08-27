package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var a App

func main() {
	godotenv.Load(".env")

	a = App{}

	a.Initilize(
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_USERNAME"))
	ensureTableExist()

	a.Run(":3001")

}

var queries []string

func ensureTableExist() {
	queries = []string{userTable, todoTable}

	if _, err := a.DB.Exec(userTable); err != nil {
		log.Fatal(err)
	}

	for _, query := range queries {

		if _, err := a.DB.Exec(query); err != nil {
			log.Fatal(err)
		}
	}

}

const userTable = `CREATE TABLE IF NOT EXISTS users
(
	id SERIAL,
	username TEXT NOT NULL,
	password TEXT NOT NULL,
	CONSTRAINT user_pkey PRIMARY KEY (id),
	CONSTRAINT user_username_unique unique (username)
)`

const todoTable = `CREATE TABLE IF NOT EXISTS todos (
	id SERIAL,
	title TEXT NOT NULL,
	status bool NOT NULL,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	author_id INT NOT NULL,
	CONSTRAINT todos_pkey PRIMARY KEY (id),
	CONSTRAINT todos_title_unique unique (title),
	CONSTRAINT fk_users
		FOREIGN KEY(author_id)
			REFERENCES users(id)
)`
