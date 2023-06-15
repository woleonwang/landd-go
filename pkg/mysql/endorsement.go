package mysql

import (
	"landd.co/landd/pkg/model"
	"time"
)

type Endorsement struct {
	ID           int64                   `gorm:"column:id;primarykey" json:"id"`
	InvitationID int64                   `gorm:"column:invitation_id" json:"invitation_id"`
	UserID       int64                   `gorm:"column:user_id" json:"user_id"`
	EndorserName string                  `gorm:"column:endorser_name" json:"endorser_name"`
	Title        string                  `gorm:"column:title" json:"title"`
	Company      string                  `gorm:"column:company" json:"company"`
	Identity     model.EndorserIdentity  `gorm:"column:identity" json:"identity"`
	Status       model.EndorsementStatus `gorm:"column:status" json:"status"`
	Content      string                  `gorm:"column:content" json:"content"`
	CreatedAt    time.Time               `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time               `gorm:"column:updated_at" json:"updated_at"`
}

func (*Endorsement) TableName() string {
	return "endorsement"
}

func GetEndorsements(userID int64, status int) ([]*Endorsement, error) {
	var endorsements []*Endorsement
	err := GetDB().Where("user_id = ? and status = ?", userID, status).
		Find(&endorsements).Error
	if err != nil {
		return nil, err
	}
	return endorsements, nil
}
