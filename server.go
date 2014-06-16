package main

import (
	"bytes"
	"encoding/json"
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

type Song struct {
	Beatmap_ID   string
	Score        string
	MaxCombo     string
	Count50      string
	Count100     string
	Count300     string
	CountMiss    string
	CountKatu    string
	CountGeki    string
	Perfect      string
	Enabled_Mods string
	User_ID      string
	Date         string
	Rank         string
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there, you are currently at %s. The URL is: %s.",
		r.URL.Path[1:], buildRecentURL(userID, 0))
}

func main() {
	setAPIKey()
	fmt.Printf("Server started on Port %s.\n", PORT)
	http.HandleFunc("/", mainPage)
	http.ListenAndServe(LISTEN_PORT, nil)
}

func setAPIKey() {
	var err error

	API_KEY, err = ioutil.ReadFile("./APIKEY.txt")
	checkError(err, "API Key")

	// Trims spaces and trailing newlines from the API key so that the URL
	// to retrieve songs can be built properly.
	API_KEY = bytes.TrimSpace(API_KEY)
	API_KEY = bytes.Trim(API_KEY, "\r\n")

	fmt.Println("API Key set to:", string(API_KEY))
}

func buildRecentURL(USER_ID string, GAME_TYPE int) string {
	return API_URL + API_RECENT_PLAYS + "?k=" + string(API_KEY) + "&u=" + USER_ID
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatal("ERROR (", msg, "): ", err)
	}
}

func getRecentPlays(url string) string {
	var songs []Song

	res, err := http.Get(url)
	defer res.Body.Close()
	checkError(err, "Get HTTP")

	html, err := ioutil.ReadAll(res.Body)
	checkError(err, "Read HTML")

	err = json.Unmarshal(html, &songs)
	checkError(err, "Unmarshal JSON")

	// Prints out each song's entry, used for debugging purposes.
	for i := 0; i < len(songs); i++ {
		fmt.Printf("%s: ID=%s> Score=%s, %s\n",
			songs[i].Date, songs[i].Beatmap_ID, songs[i].Score, songs[i].Rank)
	}

	return ""
}
