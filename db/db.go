package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)


func OpenDB() *sql.DB{
	connStr := "user=postgres dbname=postgres password=example host=localhost port=5432 sslmode=disable"
	
	conn, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to DATABSE ESTABLISHED✅")
	
	return conn
}

func PingDB(conn *sql.DB) error{
	err := conn.Ping()

	if err != nil{
		return err
	}

	return nil
}
