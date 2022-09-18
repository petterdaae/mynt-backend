package types

type Category struct {
	ID int64 `json:"id"`
	DraftCategory
}

type DraftCategory struct {
	Name     string  `json:"name"`
	ParentID *int64  `json:"parentId"`
	Color    *string `json:"color"`
	Ignore   *bool   `json:"ignore"`
}
