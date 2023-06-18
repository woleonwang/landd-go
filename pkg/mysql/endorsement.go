package mysql

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/utils"
	"time"
)

type Endorsement struct {
	ID           int64                   `gorm:"column:id;primarykey" json:"id"`
	InviteID     int64                   `gorm:"column:invite_id" json:"invite_id"`
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

func GetEndorsements(userID int64, statuses []model.EndorsementStatus, offset, limit int) ([]*Endorsement, error) {
	var endorsements []*Endorsement
	query := GetDB().Where("user_id = ?", userID)
	if len(statuses) > 0 {
		query = query.Where("status in (?)", statuses)
	}
	err := query.Offset(offset).Limit(limit).Order("created_at desc").
		Find(&endorsements).Error
	if err != nil {
		return nil, err
	}
	return endorsements, nil
}

func GetEndorsementByID(userID, inviteID int64) (*Endorsement, error) {
	endorse := Endorsement{}
	err := GetDB().Where("user_id = ? and invite_id = ?", userID, inviteID).
		First(&endorse).Error
	if err != nil {
		return nil, err
	}
	return &endorse, nil
}

func CreateEndorsement(endorse *Endorsement) (int64, error) {
	err := GetDB().Create(endorse).Error
	if err != nil {
		return 0, err
	}
	return endorse.InviteID, nil
}

func (e *Endorsement) BeforeCreate(tx *gorm.DB) error {
	id, err := utils.GenerateID()
	if err != nil {
		log.Errorf("error creating endorsement id: %v ", err)
		return err
	}
	e.InviteID = id
	return nil
}

func UpdateEndorsement(userID, inviteID int64, updates *Endorsement) error {
	res := GetDB().Model(&Endorsement{}).Where("user_id = ? and invite_id = ?", userID, inviteID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("no such endorsement found")
	}
	return nil
}
