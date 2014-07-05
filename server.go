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
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there, you are currently at %s. The URL is: %s.",
		r.URL.Path[1:], DATABASE.BuildRecentURL(GOsu.USER_ID, 0))
}

func main() {
	DATABASE.SetAPIKey()
	GOsu.SetUserID("Imvoo")
	fmt.Printf("Server started on Port %s.\n", PORT)
	_, err := GOsu.GetRecentPlays(DATABASE.BuildRecentURL(GOsu.USER_ID, 0))

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", mainPage)
	http.ListenAndServe(LISTEN_PORT, nil)
}
