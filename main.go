package main

import (
	"cloud.google.com/go/datastore"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	indexTemplate  = template.Must(template.ParseFiles("index.html"))
	tracks         = RecentlyPlayedResult{}
	refreshedToken = RefreshedToken{}
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

type Settings struct {
	ClientId     string
	ClientSecret string
	RefreshToken string
}

func main() {
	//jsonFile, err := os.Open("new_sample.json")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer jsonFile.Close()
	//byteValue, _ := ioutil.ReadAll(jsonFile)
	//json.Unmarshal(byteValue, &tracks)

	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	ctx := appengine.NewContext(r)
	client, _ := datastore.NewClient(ctx, "recently-played-music")
	settingsKey := datastore.NameKey("settings", "spotify_secrets", nil)
	settings := new(Settings)
	client.Get(ctx, settingsKey, settings)

	if refreshedToken.AccessToken == "" {
		//fmt.Println(fmt.Sprintf("%s - No Access token - requesting new token", time.Now()))
		requestNewAccessToken(w, r, settings)
	}

	if refreshedToken.Expires.Before(time.Now()) {
		//fmt.Println(fmt.Sprintf("%s - Token expired - requesting new token", time.Now()))
		requestNewAccessToken(w, r, settings)
	}

	//fmt.Println(refreshedToken.AccessToken)
	getRecentlyPlayed(refreshedToken.AccessToken, w, r)

	indexTemplate.Execute(w, tracks)
}

func getRecentlyPlayed(accessToken string, w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/recently-played?limit=50", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error")

	}
	defer resp.Body.Close()

	byteValue, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(byteValue, &tracks)
}

func requestNewAccessToken(w http.ResponseWriter, r *http.Request, settings *Settings) {

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	body := strings.NewReader(`grant_type=refresh_token&refresh_token=` + settings.RefreshToken)
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", body)
	if err != nil {
		// handle err
	}

	sEnc := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", settings.ClientId, settings.ClientSecret)))

	req.Header.Set("Authorization", "Basic "+sEnc)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	byteValue, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(byteValue, &refreshedToken)
	refreshedToken.Expires = time.Now().Add(time.Second * 3600)
}

type RefreshedToken struct {
	AccessToken string `json:"access_token"`
	Expires     time.Time
}
