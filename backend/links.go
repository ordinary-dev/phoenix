package backend

import (
    "gorm.io/gorm"
)

type Link struct {
	ID      uint64   `gorm:"primaryKey"`
	Name    string `gorm:"notNull"`
	Href    string `gorm:"notNull"`
	GroupID uint64   `gorm:"notNull"`
}

func CreateLink(db *gorm.DB, linkName string, href string, groupID uint64) (Link, error) {
    link := Link{
        Name: linkName,
        Href: href,
        GroupID: groupID,
    }
    result := db.Create(&link)
    if result.Error != nil {
        return Link{}, result.Error
    }

    return link, nil
}

func UpdateLink(db *gorm.DB, id uint64, linkName string, href string) (Link, error) {
    var link Link
    db.First(&link, id)

    link.Name = linkName
    link.Href = href
    result := db.Save(&link)
    if result.Error != nil {
        return Link{}, result.Error
    }

    return link, nil
}

func DeleteLink(db *gorm.DB, id uint64) error {
    result := db.Delete(&Link{}, id)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
