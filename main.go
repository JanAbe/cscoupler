package main

import (
	"database/sql"

	"github.com/janabe/cscoupler/util"

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
	dsn := util.GetDSN("./.secret.json")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return db
}
