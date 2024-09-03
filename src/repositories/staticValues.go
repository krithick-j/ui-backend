package repositories

import (
	"gorm.io/gorm"
)

func GetAllSocialConn(tx *gorm.DB) ([]string, error) {
	var relationship []string
	err := tx.Table("static_values").Select("social_conn").Find(&relationship).Error
	return relationship, err
}
