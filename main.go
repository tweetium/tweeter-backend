package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"                     // required driver
	_ "github.com/mattes/migrate/source/file" // required driver

	"tweeter/db"
	"tweeter/handlers"
	"tweeter/util"
)

func main() {
	port := util.MustGetEnvUInt32("PORT")
	dbURL := util.MustGetEnv("DATABASE_URL")

	err := db.Init(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %s", err)
	}

	runWebserver(port)
}

func runWebserver(port uint32) {
	http.HandleFunc("/api/v1/users", logMiddleware(handlers.UsersHandler))

	log.Printf("Listening at :%d!..", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func logMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s - %s", r.Method, r.URL)
		handler(w, r)
	}
}
