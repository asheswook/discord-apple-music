package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/asheswook/discord-apple-music/song"
	"github.com/hugolgst/rich-go/client"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	err = client.Login(os.Getenv("CLIENT_ID"))
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected to Discord")

	currentSong := song.Song{}
	for {
		song, err := song.GetNowPlaying()
		if err != nil {
			panic(err)
		}

		if currentSong["name"] != song["name"] {
			currentSong = song
			log.Println("Song changed to", song["name"], "by", song["artist"])

			duration, _ := strconv.ParseFloat(song["duration"], 64)
			durationInt := int(duration)

			nowTime := time.Now()
			endTIme := nowTime.Add(time.Second * time.Duration(durationInt))
			// endtime nowtime 음수될때 처리 필요함. 

			searchUrl := "https://music.apple.com/us/search?term=" + strings.ReplaceAll(song["name"], " ", "%20")
			fmt.Println(searchUrl)

			err = client.SetActivity(client.Activity{
				State:      song["artist"] + " — " + song["album"],
				Details:    song["name"],
				Timestamps: &client.Timestamps{
					Start: &nowTime,
					End: &endTIme,
				},
				Buttons: []*client.Button{
					&client.Button{
						Label: "Listen on Apple Music",
						Url: searchUrl,
					},
					&client.Button{
						Label: "Listen on Spotify",
						Url: "https://open.spotify.com/search/" + strings.ReplaceAll(song["name"], " ", "%20"),
					},
				},
			})
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Second * 5)
	}
}