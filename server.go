package main

import (
	"fmt"
	"github.com/Imvoo/GOsu"
	"github.com/robfig/cron"
	"gopkg.in/mgo.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	PORT           string = os.Getenv("PORT")
	LISTEN_PORT    string
	DATABASE       GOsu.Database
	USER_ID        string
	session        *mgo.Session
	collectionName string = "Imvoo"
	dbUser         string
	dbPass         string
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	songs := RetrieveSongs()

	t := template.New("song")
	t, _ = template.ParseFiles("template.html")
	t.Execute(w, songs)

	// fmt.Fprintf(w, "%s", songs)
}

func extractText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "\r\n")
	return text
}

func main() {
	// For use with local server and not Heroku.
	if PORT == "" {
		PORT = "8080"
	}
	LISTEN_PORT = ":" + PORT

	DATABASE.SetAPIKey()
	USER_ID = "Imvoo"

	// Setup the database for incoming connections.
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Start to load in the Username and Password from one of two locations;
	// From the environment variables dbUser and dbPass or in the current folder
	// under dbUser.txt and dbPass.txt.
	tempUserName, err := ioutil.ReadFile(dir + "/dbUser.txt")
	if err != nil {
		dbUser = os.Getenv("dbUser")

		if len(dbUser) == 0 {
			log.Fatal("Unable to find a username for the MongoDB database in local file (dbUser.txt) or in the environment variables under dbUser.")
		} else {
			err = nil
		}
	} else {
		dbUser = string(tempUserName)
	}

	tempPassword, err := ioutil.ReadFile(dir + "/dbPass.txt")
	if err != nil {
		dbPass = os.Getenv("dbPass")

		if len(dbPass) == 0 {
			log.Fatal("Unable to find a password for the MongoDB database in local file (dbPass.txt) or in the environment variables under dbPass.")
		} else {
			err = nil
		}
	} else {
		dbPass = string(tempPassword)
	}

	// Removes EOL, spaces etc. that may disturb the Mongo URL.
	dbUser = extractText(dbUser)
	dbPass = extractText(dbPass)

	mongoDB := "mongodb://" + string(dbUser) + ":" + string(dbPass) + "@ds053439.mongolab.com:53439/displosu"

	session, err = mgo.Dial(mongoDB)
	if err != nil {
		log.Fatal("Cannot authenticate with the database, are your credentials correct in the local files or env variables (dbUser, dbPass)?")
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	fmt.Printf("Server started on Port %s.\n", PORT)

	cronJob := cron.New()
	cronJob.AddFunc("0 * * * * *", func() { SaveRecentSongs() })
	cronJob.Start()

	http.HandleFunc("/", mainPage)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.ListenAndServe(LISTEN_PORT, nil)
}

func RetrieveSongs() []GOsu.Song {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("displosu").C(USER_ID)

	songs := []GOsu.Song{}

	err := c.Find(nil).Sort("-date").Limit(50).All(&songs)

	if err != nil {
		fmt.Println("WARN: Couldn't retrieve songs from the database.")
		return songs
	}

	return songs
}

func SaveRecentSongs() {
	recentSongs, err := DATABASE.GetRecentPlays(USER_ID, GOsu.OSU)
	if err != nil {
		fmt.Println("WARN: ", err)
	} else {
		sessionCopy := session.Copy()
		defer sessionCopy.Close()

		c := sessionCopy.DB("displosu").C(USER_ID)

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
}
