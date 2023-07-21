package backend

type Link struct {
	ID      uint64 `gorm:"primaryKey"`
	Name    string `gorm:"notNull"`
	Href    string `gorm:"notNull"`
	GroupID uint64 `gorm:"notNull"`
}
