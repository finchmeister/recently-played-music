# Recently Played Music

Web application written in Go showing the last 50 tracks I've listened to. 
Deployed to Google App Engine on the free tier.

## Commands

Install
```
export GO112MODULE=on 
go get
```

Run locally at <http://localhost:8080>:
```
go run
```

Deploy:
```
gcloud app deploy
```

View:

<https://recently-played-music.appspot.com>

Switch project:

```
gcloud config set project recently-played-music
```