module github.com/myl7/tgchan2tw

go 1.17

require (
	github.com/PuerkitoBio/goquery v1.7.1
	github.com/dghubble/go-twitter v0.0.0-20210609183100-2fdbf421508e
	github.com/dghubble/oauth1 v0.7.0
	github.com/joho/godotenv v1.3.0
	github.com/mattn/go-sqlite3 v1.14.8
	github.com/mmcdole/gofeed v1.1.3
	github.com/myl7/twitter-text-parse-go v1.0.0
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	gopkg.in/yaml.v2 v2.2.8
)

require (
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/cenkalti/backoff v2.1.1+incompatible // indirect
	github.com/dghubble/sling v1.3.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/mmcdole/goxpp v0.0.0-20181012175147-0068e33feabf // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/text v0.3.6 // indirect
)

replace github.com/myl7/twitter-text-parse-go => ./third_party/twitter-text-parse-go
