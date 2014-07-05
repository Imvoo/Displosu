package main

import (
	"fmt"
	"github.com/Imvoo/GOsu"
	"log"
	"net/http"
)

var (
	PORT        string = "8080"
	LISTEN_PORT string = ":" + PORT
	DATABASE    GOsu.Database
	USER_ID     string
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there, you are currently at %s. The URL is: %s.",
		r.URL.Path[1:], DATABASE.BuildRecentURL(USER_ID, 0))
}

func main() {
	DATABASE.SetAPIKey()
	USER_ID = "Imvoo"
	fmt.Printf("Server started on Port %s.\n", PORT)
	songs, err := DATABASE.GetRecentPlays(USER_ID, 0)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(songs)

	http.HandleFunc("/", mainPage)
	http.ListenAndServe(LISTEN_PORT, nil)
}
