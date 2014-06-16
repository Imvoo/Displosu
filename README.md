# Displosu!
###### Website displaying recently played songs for Osu!

### Overview
This is a remake of my original OsuScoreRetriever project, the main difference being the usage of the Osu! API and it being written in GoLang. The main aim of this is to display a user's history of plays alongside many fun stats, kind of like a database of song plays.

### Requirements
If you want to build and run this yourself, you need the following:

- Go (v1.2.2)
- Osu! API Key (https://osu.ppy.sh/p/api)


Inside the directory (e.g. Displosu/), run in your terminal:

    $    go build

You must then create a file called APIKEY.txt inside the directory and paste in your API key.

After that, if running on Linux run:

    $    ./Displosu.exe

or if you're using Windows,

    $    Displosu

and your server will start on Port 8080.

Alternatively, you can also run the Displosu.exe file from a File Browser (e.g. Windows Explorer).

### Acknowledgements
- Peppy (https://github.com/peppy), for including the Recently Played aspect of the Osu! Api on request.
- The makers of GoLang.
- Heroku, for website hosting.
