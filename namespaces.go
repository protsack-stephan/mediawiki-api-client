package mediawiki

const namespacesURL = "/w/api.php?action=query&format=json&meta=siteinfo&siprop=namespaces&formatversion=2"

// Namespace single namespace
type Namespace struct {
	ID            int    `json:"id"`
	Case          string `json:"case"`
	Name          string `json:"name"`
	Subpages      bool   `json:"subpages"`
	Canonical     string `json:"canonical"`
	Content       bool   `json:"content"`
	Nonincludable bool   `json:"nonincludable"`
}

type namespacesResponse struct {
	Query struct {
		Namespaces map[int]Namespace
	}
}
