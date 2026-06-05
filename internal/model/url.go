package model

type Url struct {
	ID           uint    `json:"id"`
	OriginalUrl string `json:"original_url"`
	UrlHash     string `json:"url_hash"`
	ShortUrl    string `json:"short_url"`
}
