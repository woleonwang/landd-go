package mysql

import (
	"landd.co/landd/pkg/model"
	"time"
)

type EndorsementDraft struct {
	ID           int64                  `gorm:"column:id;primarykey" json:"id"`
	UserID       int64                  `gorm:"column:user_id" json:"user_id"`
	EndorserName string                 `gorm:"column:endorser_name" json:"endorser_name"`
	Title        string                 `gorm:"column:title" json:"title"`
	Company      string                 `gorm:"column:company" json:"company"`
	Identity     model.EndorserIdentity `gorm:"column:identity" json:"identity"`
	Content      string                 `gorm:"column:content" json:"content"`
	CreatedAt    time.Time              `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time              `gorm:"column:updated_at" json:"updated_at"`
}

func (*EndorsementDraft) TableName() string {
	return "endorsement_draft"
}

func GetEndorsementDraft(userID int64) (*EndorsementDraft, error) {
	var draft EndorsementDraft
	err := GetDB().Where("user_id = ?", userID).First(&draft).Error
	if err != nil {
		return nil, err
	}
	return &draft, nil
}

func UpdateEndorsementDraft(userID int64, updates *EndorsementDraft) error {
	res := GetDB().Model(&EndorsementDraft{}).Where("user_id = ?", userID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return GetDB().Create(updates).Error
	}
	return nil
}
