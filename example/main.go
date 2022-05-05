package main

import (
	"context"
	"fmt"
	"log"

	"github.com/protsack-stephan/mediawiki-api-client"
)

func main() {
	client := mediawiki.
		NewBuilder("https://en.wikipedia.org").
		Headers(map[string]string{
			"User-Agent": "test@gmail.com",
		}).
		Build()
	ctx := context.Background()

	meta, err := client.PageMeta(ctx, "Barack_Obama")

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(meta)

	data, err := client.PageHTML(ctx, "Barack_Obama", meta.Rev)

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(data))

	revisions, err := client.PageRevisions(ctx, "Barack_Obama", 10,
		mediawiki.PageRevisionsOptions{
			Order: mediawiki.RevisionOrderingOlder,
			Props: []string{"content", "ids", "timestamp"},
		},
	)

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

	wikitext, err := client.PageWikitext(ctx, "Barack_Obama")

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(wikitext))

	psdata, err := client.PagesData(ctx, []string{"Barack_Obama", "Earth"}, mediawiki.PageDataOptions{
		RevisionsLimit: 2,
		RevisionProps:  []string{"ids", "timestamp", "content", "sha1"},
	})

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(psdata)

	pdata, err := client.PageData(ctx, "Barack_Obama")

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(pdata)

	users, err := client.Users(ctx, pdata.Revisions[0].UserID, 3333333)

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(users)

	user, err := client.User(ctx, pdata.Revisions[0].UserID)

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(user)

	// An example of working with page category
	opt := mediawiki.PageDataOptions{
		CategoriesLimit: 1,
		CategoriesProps: []string{"hidden"},
	}
	pdata, err = client.PageData(ctx, "Rideau_River", opt)

	if err != nil {
		log.Panic(err)
	}
	if len(pdata.Categories) != 1 {
		fmt.Println("Error: PageData response should have only one category.")
	}
	fmt.Println("Category namespace : ", pdata.Categories[0].Ns)
	fmt.Println("Category title : ", pdata.Categories[0].Title)
	fmt.Println("Category prop hidden : ", pdata.Categories[0].Hidden)
}
