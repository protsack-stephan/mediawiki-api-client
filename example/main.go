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
	ctx := context.Background()
	meta, status, err := client.PageMeta(ctx, "Pet_door")

	if err != nil {
		log.Panic(err)
	}

	if status != http.StatusOK {
		log.Panic("bad request")
	}

	fmt.Println(meta)

	_, status, err = client.PageHTML(ctx, "Pet_door", meta.Rev)

	if err != nil {
		log.Panic(err)
	}

	if status != http.StatusOK {
		log.Panic("bad request")
	}

	revisions, status, err := client.PageRevisions(context.Background(), "Pet_door", 10)

	if err != nil {
		log.Panic(err)
	}

	if status != http.StatusOK {
		log.Panic("bad request")
	}

	for _, rev := range revisions {
		fmt.Println(rev)
	}

	matrix, status, err := client.Sitematrix(ctx)

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

	namespaces, status, err := client.Namespaces(ctx)

	if err != nil {
		log.Panic(err)
	}

	if status != http.StatusOK {
		log.Panic("bad request")
	}

	for _, ns := range namespaces {
		fmt.Println(ns)
	}
}
