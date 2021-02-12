# Mediawiki API client for Go

This is golang client for accessing Mediawiki API endpoints, such as [REST API](https://en.wikipedia.org/api/rest_v1/#/) and [Actions API](https://www.mediawiki.org/wiki/API:Main_page). Purpose of this package is to provide some kind of SDK like abstraction for developers. Intention here take away headache of parsing JSON responses and keeping track of API endpoints to provide clear and easy to use way for access to the Mediawiki data.

Small example of getting page meta data:
```go
client := mediawiki.NewClient("https://en.wikipedia.org/") // creating the client

ctx := context.Background() // getting context instance

meta, err := client.PageMeta(ctx, "Pet_door") // accessing "Pet_door" page meta data

if err != nil {
  log.Panic(err)
}

fmt.Println(meta)
```

### *Note that we are far from supporting all API endpoints here and will be adding more support based on our needs. Feel free to open a PR and add new endpoints to support.