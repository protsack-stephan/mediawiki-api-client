package main

import (
	"context"
	"fmt"
	"log"

	"github.com/protsack-stephan/mediawiki-api-client"
)

func main() {
	client := mediawiki.NewClient("https://en.wikipedia.org/")
	ctx := context.Background()

	meta, err := client.PageMeta(ctx, "Pet_door")

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(meta)

	data, err := client.PageHTML(ctx, "Pet_door", meta.Rev)

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(data))

	revisions, err := client.PageRevisions(ctx, "Pet_door", 10)

	if err != nil {
		log.Panic(err)
	}

	for _, rev := range revisions {
		fmt.Println(rev)
	}

	matrix, err := client.Sitematrix(ctx)

	if err != nil {
		log.Panic(err)
	}

	for _, project := range matrix.Projects {
		fmt.Println(project)
	}

	for _, special := range matrix.Specials {
		fmt.Println(special)
	}

	namespaces, err := client.Namespaces(ctx)

	if err != nil {
		log.Panic(err)
	}

	for _, ns := range namespaces {
		fmt.Println(ns)
	}

	wikitext, err := client.PageWikitext(ctx, "Main")

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(wikitext))
}
