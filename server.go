package main

import (
	"fmt"
	"github.com/Imvoo/GOsu"
	"github.com/robfig/cron"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	PORT           string = "8080"
	LISTEN_PORT    string = ":" + PORT
	DATABASE       GOsu.Database
	USER_ID        string
	session        *mgo.Session
	collectionName string = "Imvoo"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there, you are currently at %s.", r.URL.Path[1:])
}

func main() {
	DATABASE.SetAPIKey()
	USER_ID = "Imvoo"

	// Setup the database for incoming connections.
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	userName, err := ioutil.ReadFile(dir + "\\dbUser.txt")
	if err != nil {
		log.Fatal(err)
	}

	password, err := ioutil.ReadFile(dir + "\\dbPass.txt")
	if err != nil {
		log.Fatal(err)
	}

	mongoDB := "mongodb://" + string(userName) + ":" + string(password) + "@ds053439.mongolab.com:53439/displosu"

	session, err = mgo.Dial(mongoDB)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	fmt.Printf("Server started on Port %s.\n", PORT)

	c := cron.New()
	c.AddFunc("0 * * * * *", func() { RecentSongs() })
	c.Start()

	http.HandleFunc("/", mainPage)
	http.ListenAndServe(LISTEN_PORT, nil)
}

func RecentSongs() {
	recentSongs, err := DATABASE.GetRecentPlays(USER_ID, GOsu.OSU)
	if err != nil {
		log.Fatal(err)
	}

	c := session.DB("displosu").C(USER_ID)

	// Grabs the latest song for tracking which songs to record.
	emptyDatabase := false
	latestSong := GOsu.Song{}
	err = c.Find(nil).Sort("-date").One(&latestSong)
	if err != nil {
		fmt.Println("WARN: no songs found in database!")
		emptyDatabase = true
	}

	var latestTime time.Time
	if !emptyDatabase {
		latestTime, _ = time.Parse("2006-01-02 15:04:05", latestSong.Date)
	}

	for _, result := range recentSongs {
		if !emptyDatabase {
			resultTime, err := time.Parse("2006-01-02 15:04:05", result.Date)
			if err != nil {
				log.Fatal(err)
			}

			if resultTime.After(latestTime) {
				fmt.Printf("Inserted song @ %s scoring %s.\n", result.Date, result.Score)
				err = c.Insert(result)

				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			fmt.Printf("Inserted song @ %s scoring %s.\n", result.Date, result.Score)
			err = c.Insert(result)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
