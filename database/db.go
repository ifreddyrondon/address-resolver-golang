package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	tableCreationQuery = `CREATE TABLE IF NOT EXISTS addresses (
				id SERIAL,
				address TEXT NOT NULL,
				lat NUMERIC(10,2) NOT NULL,
				lng NUMERIC(10,2) NOT NULL,
				CONSTRAINT addresses_pkey PRIMARY KEY (id))`
	deleteDataQuery        = "DELETE FROM addresses"
	restartIdSequenceQuery = "ALTER SEQUENCE addresses_id_seq RESTART WITH 1"
)

var db *sql.DB

func CreateConnection(connectionUrl string) {
	var err error
	db, err = sql.Open("postgres", connectionUrl)
	if err != nil {
		log.Fatal(err)
	}

	pingConnection()
	ensureTableExists()
}

func pingConnection() {
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func GetDB() *sql.DB {
	return db
}

func ensureTableExists() {
	if _, err := db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func ClearTable() {
	db.Exec(deleteDataQuery)
	db.Exec(restartIdSequenceQuery)
}
