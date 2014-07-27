# GOsu
###### Osu! API Wrapper written in GoLang.

### Overview
This wrapper allows you to use the Osu! API in GoLang provided you have an API key to use.

#### Requirements

- Go (v1.2.2 or higher)
- Osu! API Key (https://osu.ppy.sh/p/api)

#### Installing

Assuming your GoPath environment variable is setup properly, in your terminal type in:

	go get github.com/Imvoo/GOsu

This will install GOsu to your GoLang path and then you can call it in your Go files with the import path github.com/Imvoo/GOsu.

#### Quick Start

Say we want to grab the latest plays for a specific user.

1. We first want to grab our API key and paste it in a file in the project directory under the name APIKEY.txt.

2. We then start to setup the database to handle the connections to the Osu! Api.

		package main

		import "github.com/Imvoo/GOsu"

		var (
			DATABASE GOsu.Database
			USER_ID	 string
		)

		func main() {
			DATABASE.SetAPIKey()
			USER_ID = "Imvoo" // Set this to your ID.
		}

3. After setting the database up, we can request the recent plays of the specified user by adding to the main function.

		func main() {
			DATABASE.SetAPIKey()
			USER_ID = "Imvoo" // Set this to your ID.
			songs, err := DATABASE.GetRecentPlays(USER_ID, GOsu.OSU) // GOsu.OSU represents the Osu! gametype.

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(songs)
		}

4. If successful, you'll see either see:
	
		$	[] 

	get outputted which means you haven't played any songs yet, or

		$	[{338547 115512 66 1 4 73 10 3 15 0 0 949789 2014-07-05 13:30:47 F}]

	as an example, which means it's working properly. It may look like a bunch of random numbers, but it's the recently played song's attributes in the following order: 
	Beatmap_ID, Score, MaxCombo, Count50, Count100, Count300, CountMiss, CountKatu, CountGeki, Perfect, Enabled_Mods, User_ID, Date, Rank.

	This can later be accessed using:

		songs[0].Beatmap_ID
		songs[0].Score
		...

	etc, where 0 is the index of the song.

#### Timeline

- [x] Get User Recent
- [x] Get Beatmaps
- [x] Get User
- [x] Get Scores
- [x] Get User Best
- [x] Get Match
- [x] Mods


