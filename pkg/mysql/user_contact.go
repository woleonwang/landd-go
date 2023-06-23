package mysql

import (
	log "github.com/sirupsen/logrus"
	"landd.co/landd/pkg/model"
	"time"
)

type UserContact struct {
	ID        int64             `gorm:"column:id;primarykey" json:"id"`
	UserID    int64             `gorm:"column:user_id" json:"user_id"`
	Contact   string            `gorm:"column:contact" json:"contact"`
	Type      model.ContactType `gorm:"column:type" json:"type"`
	CreatedAt time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time         `gorm:"column:updated_at" json:"updated_at"`
}

func (*UserContact) TableName() string {
	return "user_contact"
}

func CreateUserContact(contacts []*UserContact) error {
	err := GetDB().Create(contacts).Error
	if err != nil {
		log.Errorf("error creating user contacts: %v ", err)
		return err
	}
	return nil
}

func GetUserContacts(userID int64) ([]*UserContact, error) {
	var contacts []*UserContact
	err := GetDB().Where("user_id = ?", userID).Find(&contacts).Error
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func SaveUserContacts(userID int64, contacts []*UserContact) error {
	if err := GetDB().Where("user_id = ?", userID).Delete(&UserContact{}).Error; err != nil {
		log.Errorf("error deleting contacts: %v ", err)
		return err
	}
	if len(contacts) > 0 {
		if err := GetDB().Create(contacts).Error; err != nil {
			log.Errorf("error creating contacts: %v ", err)
			return err
		}
	}
	return nil
}
