package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	PORT             string = "8080"
	LISTEN_PORT      string = ":" + PORT
	API_KEY          []byte
	API_URL          string = "https://osu.ppy.sh/api/"
	API_RECENT_PLAYS string = "get_user_recent"
	userID           string = "Imvoo"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there, you are currently at %s. The URL is: %s. Response: %s",
		r.URL.Path[1:], buildRecentURL("Imvoo", 0), getRecentPlays(buildRecentURL("Imvoo", 0)))
}

func main() {
	setAPIKey()
	fmt.Printf("Server started on Port %s.", PORT)
	http.HandleFunc("/", mainPage)
	http.ListenAndServe(LISTEN_PORT, nil)
}

func setAPIKey() {
	var err error
	API_KEY, err = ioutil.ReadFile("./APIKEY.txt")

	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}

func buildRecentURL(USER_ID string, GAME_TYPE int) string {
	return API_URL + API_RECENT_PLAYS + "?k=" + string(API_KEY) + "&u=" + USER_ID
}

func getRecentPlays(url string) string {
	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	html, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	return string(html)
}