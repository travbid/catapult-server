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
	dbConfig := flag.String("db_config", "", "Path to a JSON config to generate a database from")
	flag.Parse()

	fail := false
	if *port == 0 {
		log.Println("port not set")
		fail = true
	}
	if *database == "" && *dbConfig == "" {
		log.Println("Neither database or db_config set")
		fail = true
	}
	if *database != "" && *dbConfig != "" {
		log.Println("Both database and db_config set")
		fail = true
	}
	if fail {
		flag.PrintDefaults()
		os.Exit(1)
	}

	dbPath := *database
	var err error
	if *dbConfig != "" {
		log.Println("Generating database...")
		dbPath, err = lib.GenerateDatabase(*dbConfig)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println("Initialising database...")
	db, err := lib.InitDatabase(dbPath)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Creating routes...")
	router := lib.NewRouter(db)

	log.Println("Listening for connections...")
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(*port), router))
}
