package urls

type UrlEntry struct {
	Id       string `json:"id"`
	UrlLong  string `json:"urlLong"`
	UrlShort string `json:"urlShort"`
	Date     string `json:"date"`
}

var UrlEntries = []UrlEntry{
	{Id: "1", UrlLong: "www.google.com", UrlShort: "bit.ly/7thyF", Date: "2024-02-20"},
	{Id: "2", UrlLong: "www.bing.com", UrlShort: "bit.ly/h6Y4f", Date: "2024-02-20"},
}
