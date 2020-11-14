package mediawiki

const sitematrixURL = "/w/api.php?action=sitematrix&format=json&formatversion=2"

// Site wikimedia site
type Site struct {
	URL      string `json:"url"`
	DBName   string `json:"dbname"`
	Code     string `json:"code"`
	Sitename string `json:"sitename"`
	Closed   bool   `json:"closed,omitempty"`
}

// Project project and his sites
type Project struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	Site      []Site `json:"site"`
	Dir       string `json:"dir"`
	Localname string `json:"localname"`
}

// Special project with some special attributes
type Special struct {
	URL       string `json:"url"`
	DBName    string `json:"dbname"`
	Code      string `json:"code"`
	Lang      string `json:"lang"`
	Sitename  string `json:"sitename"`
	Private   bool   `json:"private,omitempty"`
	Closed    bool   `json:"closed,omitempty"`
	Fishbowl  bool   `json:"fishbowl,omitempty"`
	Nonglobal bool   `json:"nonglobal,omitempty"`
}

// Sitematrix all wikimedia projects
type Sitematrix struct {
	Count    int       `json:"count"`
	Specials []Special `json:"specials"`
	Projects []Project `json:"projects"`
}

type siteMatrixMainResponse struct {
	Sitematrix map[string]Project `json:"sitematrix"`
}

type siteMatrixSpecialResponce struct {
	Sitematrix struct {
		Count    int       `json:"count"`
		Specials []Special `json:"specials"`
	} `json:"sitematrix"`
}
