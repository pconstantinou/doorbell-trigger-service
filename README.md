
To get started, create an account at pusher.com and get an API key

Copy pusherConfig.go.default to pusherConfig.go and set the key, secret and API id
```
cp pusherConfig.go.default pusherConfig.go
```

Test the code locally:
```
go run *.go
```

Once everything looks good, push it to Google cloud:
```
gcloud app deploy
```

