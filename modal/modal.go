package modal

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"munchkin/system"
)

type Database struct {
	DB *sql.DB
}

func (d *Database) connect() {
	var err error

	d.DB, err = sql.Open("postgres", system.GetConfig().ConnectPostgresString())
	if err != nil {
		log.Printf("[DB] Open sql %v", err)
	}

	if err := d.DB.Ping(); err != nil {
		log.Printf("[DB] Ping sql %v", err)
	}
}

func New() *Database {
	db := new(Database)
	db.connect()
	return db
}
