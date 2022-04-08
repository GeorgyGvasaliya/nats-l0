package psql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1q2w3e"
	dbname   = "wb-3"
)

func NewPostgresDB() *sqlx.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sqlx.Open("postgres", psqlconn)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to db")
	return db
}
