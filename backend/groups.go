package backend

import (
	"gorm.io/gorm"
)

type Group struct {
	ID    uint64 `gorm:"primaryKey"`
	Name  string `gorm:"unique,notNull"`
	Links []Link `gorm:"constraint:OnDelete:CASCADE;"`
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

func UpdateGroup(db *gorm.DB, id uint64, groupName string) (Group, error) {
	var group Group
	db.First(&group, id)

	group.Name = groupName
	result := db.Save(&group)
	if result.Error != nil {
		return Group{}, result.Error
	}

	return group, nil
}

func DeleteGroup(db *gorm.DB, id uint64) error {
	result := db.Delete(&Group{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
