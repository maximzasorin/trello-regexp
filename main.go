package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/bbolt"
	"github.com/gorilla/mux"

	"github.com/maximzasorin/trello-regexp/auth"
	"github.com/maximzasorin/trello-regexp/rest"
	"github.com/maximzasorin/trello-regexp/store"
)

func main() {
	var (
		appName      = flag.String("name", "Trello Regexp", "Name of App")
		appURL       = flag.String("url", "http://localhost:8080", "App url.")
		dbFile       = flag.String("file", "store.db", "Storage file.")
		trelloKey    = flag.String("trelloKey", "", "Trello key from https://trello.com/1/appKey/generate.")
		trelloSecret = flag.String("trelloSecret", "", "Trello secret from https://trello.com/1/appKey/generate.")
		secret       = flag.String("secret", "", "Secret for generate JWT.")
		cookieName   = flag.String("cookieName", "trello_regexp", "Name of JWT cookie.")
	)
	flag.Parse()

	// Create store
	db, err := bolt.Open(*dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	st := store.NewStore(db)

	// Create jwt
	jwt := auth.NewJwt(*secret, *cookieName)

	// Create auth
	auth := auth.NewAuth(st, jwt, &auth.Config{
		Name:         *appName,
		CallbackURL:  *appURL + "/auth/callback",
		TrelloKey:    *trelloKey,
		TrelloSecret: *trelloSecret,
	})

	// Routes
	r := mux.NewRouter()
	r.HandleFunc("/", serveHomePage)
	r.HandleFunc("/auth", auth.GetRedirectHandler())
	r.HandleFunc("/auth/callback", auth.GetCallbackHandler())

	// API
	rest := rest.NewRest(jwt, st)
	s := r.PathPrefix("/api").Subrouter()
	s.Use(rest.GetAuthMiddleware())
	rest.Expose(s)

	http.Handle("/", r)

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
