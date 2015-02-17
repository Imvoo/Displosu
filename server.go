package main

import (
	"encoding/json"
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
	"strconv"
	"time"
)

var (
	LISTEN_PORT string
	DATABASE    GOsu.Database
	USER_ID     string
	session     *mgo.Session
)

type Config struct {
	ApiKey     string
	DBURL      string
	DBUsername string
	DBPassword string
	Port       int
	SaveSongs  bool
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	songs := RetrieveSongs()

	t := template.New("song")
	t, _ = template.ParseFiles("template.html")
	t.Execute(w, songs)
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	confFile, err := ioutil.ReadFile(dir + "/conf.json")
	if err != nil {
		log.Fatal("Could not find configuration file (conf.json).")
	}

	var Configuration Config
	err = json.Unmarshal(confFile, &Configuration)

	LISTEN_PORT = ":" + strconv.Itoa(Configuration.Port)
	DATABASE.SetAPIKey(Configuration.ApiKey)
	USER_ID = "Imvoo"

	mongoDB := "mongodb://" + Configuration.DBUsername + ":" + Configuration.DBPassword + "@" + Configuration.DBURL

	session, err = mgo.Dial(mongoDB)
	if err != nil {
		log.Fatal("Cannot authenticate with the database, are your credentials correct in the conf.json file?")
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	fmt.Printf("INFO: Started on Port %s.\n", strconv.Itoa(Configuration.Port))

	if Configuration.SaveSongs {
		cronJob := cron.New()
		cronJob.AddFunc("0 * * * * *", func() { SaveRecentSongs() })
		cronJob.Start()
		fmt.Println("INFO: Recording new songs.")
	} else {
		fmt.Println("INFO: NOT recording new songs.")
	}

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
