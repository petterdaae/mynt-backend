package types

type Name struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Fields struct {
		Regex       string `json:"regex"`
		ReplaceWith string `json:"replaceWith"`
	} `json:"fields"`
}
