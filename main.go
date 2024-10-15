package main

import "github.com/homily707/ohmygo/crawler"

func main() {
	//crawler.ImportedByRepo("https://pkg.go.dev/github.com/segmentio/kafka-go?tab=importedby")
	crawler.GetStars("github.com/segmentio/kafka-go")
}
