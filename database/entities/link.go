package entities

type Link struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Href    string  `json:"href"`
	GroupID int     `json:"-"`
	Icon    *string `json:"icon,omitempty"`
}
