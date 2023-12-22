package main

import (
	_ "github.com/lib/pq"
	"assignment-2/app"
	"log"
)

func main() {
	app.InitDB()
	log.Println("Initialized the DB successfully")
}