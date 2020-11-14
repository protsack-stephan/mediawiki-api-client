package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/protsack-stephan/mediawiki-api-client"
)

func main() {
	client := mediawiki.NewClient("https://en.wikipedia.org/")
	meta, status, err := client.PageMeta(context.Background(), "Pet_door")

	if err != nil {
		log.Panic(err)
	}

	if status != http.StatusOK {
		log.Panic("bad request")
	}

	fmt.Println(meta)
}
