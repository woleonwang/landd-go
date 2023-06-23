package mysql

import (
	"time"
)

type PartnerProfile struct {
	ID        int64     `gorm:"column:id;primarykey" json:"id"`
	UserID    int64     `gorm:"column:user_id" json:"user_id"`
	Name      string    `gorm:"column:name" json:"name"`
	Photo     string    `gorm:"column:photo" json:"photo"`
	Company   string    `gorm:"column:company" json:"company"`
	Twitter   string    `gorm:"column:twitter" json:"twitter"`
	LinkedIn  string    `gorm:"column:linkedin" json:"linkedin"`
	Website   string    `gorm:"column:website" json:"website"`
	Blog      string    `gorm:"column:blog" json:"blog"`
	Facebook  string    `gorm:"column:facebook" json:"facebook"`
	Instagram string    `gorm:"column:instagram" json:"instagram"`
	Tiktok    string    `gorm:"column:tiktok" json:"tiktok"`
	Youtube   string    `gorm:"column:youtube" json:"youtube"`
	Other     string    `gorm:"column:other" json:"other"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (*PartnerProfile) TableName() string {
	return "partner_profile"
}

func GetPartnerProfile(userID int64) (*PartnerProfile, error) {
	var profile PartnerProfile
	err := GetDB().Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func UpdatePartnerProfile(userID int64, updates *PartnerProfile) error {
	res := GetDB().Model(&PartnerProfile{}).
		Where("user_id = ?", userID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return GetDB().Create(updates).Error
	}
	return nil
}
