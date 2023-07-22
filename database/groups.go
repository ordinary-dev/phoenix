package database

type Group struct {
	ID    uint64 `gorm:"primaryKey"`
	Name  string `gorm:"unique,notNull"`
	Links []Link `gorm:"constraint:OnDelete:CASCADE;"`
}
