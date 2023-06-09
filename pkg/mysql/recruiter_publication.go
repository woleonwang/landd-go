package mysql

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type RecruiterPublication struct {
	ID        int64     `gorm:"column:id;primarykey" json:"id"`
	UserID    int64     `gorm:"column:user_id" json:"user_id"`
	Title     string    `gorm:"column:title" json:"title"`
	Link      string    `gorm:"column:link" json:"link"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (*RecruiterPublication) TableName() string {
	return "recruiter_publication"
}

func GetRecruiterPublication(userID int64) ([]*RecruiterPublication, error) {
	var pubs []*RecruiterPublication
	err := GetDB().Where("user_id = ?", userID).Find(&pubs).Error
	if err != nil {
		return nil, err
	}
	return pubs, nil
}

func SaveRecruiterPublication(userID int64, pubs []*RecruiterPublication) error {
	if err := GetDB().Where("user_id = ?", userID).Delete(&RecruiterPublication{}).Error; err != nil {
		log.Errorf("error deleting publications: %v ", err)
		return err
	}
	if len(pubs) > 0 {
		if err := GetDB().Create(pubs).Error; err != nil {
			log.Errorf("error creating publications: %v ", err)
			return err
		}
	}
	return nil
}
