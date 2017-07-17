package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS addresses
(
id SERIAL,
address TEXT NOT NULL,
lat NUMERIC(10,2) NOT NULL,
lng NUMERIC(10,2) NOT NULL,
CONSTRAINT addresses_pkey PRIMARY KEY (id)
)
`

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
	db.Exec("DELETE FROM addresses")
	db.Exec("ALTER SEQUENCE addresses_id_seq RESTART WITH 1")
}
