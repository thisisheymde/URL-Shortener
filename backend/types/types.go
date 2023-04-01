package types

type Link struct {
	ID  string `json:"id" redis:"id"`
	URL string `json:"url" redis:"id"`
}
