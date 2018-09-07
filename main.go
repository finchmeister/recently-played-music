package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"google.golang.org/appengine"
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

type templateParams struct {
	Name string
}

type RecentlyPlayedItem struct {
	// Track is the track information
	Track SimpleTrack `json:"track"`

	// PlayedAt is the time that this song was played
	PlayedAt time.Time `json:"played_at"`

	// PlaybackContext is the current playback context
	PlaybackContext PlaybackContext `json:"context"`
}

type RecentlyPlayedResult struct {
	Items []RecentlyPlayedItem `json:"items"`
}

// PlaybackContext is the playback context
type PlaybackContext struct {
	ExternalURLs map[string]string `json:"external_urls"`
	Endpoint     string            `json:"href"`
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
}

type SimpleTrack struct {
	Album      SimpleAlbum    `json:"album"`
	Artists    []SimpleArtist `json:"artists"`
	Endpoint   string         `json:"href"`
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	PreviewURL string         `json:"preview_url"`
	URI        string         `json:"uri"`
}

type SimpleArtist struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	URI      string `json:"uri"`
	Endpoint string `json:"href"`
}

type SimpleAlbum struct {
	Name   string  `json:"name"`
	Images []Image `json:"images"`
}

type Image struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

type songDetails struct {
	Name   string
	Artist string
	Img    string
}

func main() {

	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		fmt.Println("bb")
		return
	}

	//tracks := RecentlyPlayedResult{}

	jsonFile, err := os.Open("new_sample.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var tracks RecentlyPlayedResult

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &tracks)

	//tracks.Items[0].Track.Album.Images

	indexTemplate.Execute(w, tracks)
}
