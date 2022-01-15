package types

type Budget struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Color *string `json:"color"`
}
