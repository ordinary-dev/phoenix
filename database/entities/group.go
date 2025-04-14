package entities

type Group struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Username *string `json:"-"`
	Links    []Link  `json:"links"`
}
