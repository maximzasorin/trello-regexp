package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/bbolt"

	"github.com/maximzasorin/trello-regexp/auth"
	"github.com/maximzasorin/trello-regexp/store"
)

var (
	appName      *string
	appURL       *string
	dbFile       *string
	trelloKey    *string
	trelloSecret *string
)

func main() {
	appName = flag.String(
		"name",
		"Trello Regexp",
		"Name of App",
	)
	appURL = flag.String(
		"url",
		"http://localhost:8080",
		"App url.",
	)
	dbFile = flag.String(
		"file",
		"bolt.db",
		"Storage file.",
	)
	trelloKey = flag.String(
		"key",
		"",
		"Trello key from https://trello.com/1/appKey/generate.",
	)
	trelloSecret = flag.String(
		"secret",
		"",
		"Trello secret from https://trello.com/1/appKey/generate.",
	)

	flag.Parse()

	// Create store
	db, err := bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	store := store.NewStore(db)

	// Create auth
	auth := auth.NewAuth(&auth.Config{
		Name:        *appName,
		CallbackURL: *appURL + "/auth/callback",
		Key:         *trelloKey,
		Secret:      *trelloSecret,
	}, store)

	// Routes
	http.HandleFunc("/", serveHomePage)
	http.HandleFunc("/auth", auth.GetRedirectHandler())
	http.HandleFunc("/auth/callback", auth.GetCallbackHandler())

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8" />
			<title>Trello Regex</title>
		</head>
		<body>
			<a href="/auth">Login</a>
		</body>
	</html>`)
}