package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	pass   = "root"
	dbname = "FoodOrder"
)

func OpenDBConnection() *sql.DB {
	connectionString := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", host, port, user, pass, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("DB connection success.")

	return db
}
