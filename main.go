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
	tracks        = RecentlyPlayedResult{}
)

type RecentlyPlayedResult struct {
	Items []RecentlyPlayedItem `json:"items"`
}

type RecentlyPlayedItem struct {
	Track           SimpleTrack     `json:"track"`
	PlayedAt        time.Time       `json:"played_at"`
	PlaybackContext PlaybackContext `json:"context"`
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

type PlaybackContext struct {
	ExternalURLs map[string]string `json:"external_urls"`
	Endpoint     string            `json:"href"`
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
}

type SimpleAlbum struct {
	Name   string  `json:"name"`
	Images []Image `json:"images"`
}

type SimpleArtist struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	URI      string `json:"uri"`
	Endpoint string `json:"href"`
}

type Image struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

func main() {
	jsonFile, err := os.Open("new_sample.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &tracks)

	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		fmt.Println("bb")
		return
	}
	indexTemplate.Execute(w, tracks)
}
