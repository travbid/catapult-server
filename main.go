package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	lib "github.com/travbid/catapult-server/lib"
)

func main() {
	port := flag.Uint("port", 0, "The port number to serve on")
	database := flag.String("database", "", "Path to the sqlite database")
	flag.Parse()

	fail := false
	if *port == 0 {
		log.Println("port not set")
		fail = true
	}
	if *database == "" {
		log.Println("database not set")
		fail = true
	}
	if fail {
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Println("Initialising database...")
	db, err := lib.InitDatabase(*database)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Creating routes...")
	router := lib.NewRouter(db)

	log.Println("Listening for connections...")
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(*port), router))
}
