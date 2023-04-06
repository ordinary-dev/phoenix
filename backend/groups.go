package backend

import (
    "gorm.io/gorm"
)

type Group struct {
	ID    uint64   `gorm:"primaryKey"`
	Name  string `gorm:"unique,notNull"`
	Links []Link
}

func GetGroups(db *gorm.DB) ([]Group, error) {
    var groups []Group
    result := db.Model(&Group{}).Preload("Links").Find(&groups)
    if result.Error != nil {
        return nil, result.Error
    }
    return groups, nil
}

func CreateGroup(db *gorm.DB, groupName string) (Group, error) {
    group := Group{
        Name: groupName,
    }
    result := db.Create(&group)
    if result.Error != nil {
        return Group{}, result.Error
    }

    return group, nil
}
