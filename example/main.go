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

	_, status, err = client.PageHTML(context.Background(), "Pet_door", meta.Rev)

	if err != nil {
		log.Panic(err)
	}

	if status != http.StatusOK {
		log.Panic("bad request")
	}

	matrix, status, err := client.Sitematrix(context.Background())

	if err != nil {
		log.Panic(err)
	}

	if status != http.StatusOK {
		log.Panic("bad request")
	}

	for _, project := range matrix.Projects {
		fmt.Println(project)
	}

	for _, special := range matrix.Specials {
		fmt.Println(special)
	}
}
