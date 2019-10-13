package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/janabe/cscoupler/server"
)

func main() {
	db := ConnectToDB()
	defer db.Close()
	server := server.NewServer(db)
	server.Run()
}

// ConnectToDB connects to the database
// and returns
func ConnectToDB() *sql.DB {
	dsn := "user=postgres password=[password] dbname=cscoupler sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return db
}
