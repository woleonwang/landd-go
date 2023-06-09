package mysql

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type RecruiterCandidate struct {
	ID         int64     `gorm:"column:id;primarykey" json:"id"`
	UserID     int64     `gorm:"column:user_id" json:"user_id"`
	Title      string    `gorm:"column:title" json:"title"`
	Percentage float32   `gorm:"column:percentage" json:"percentage"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (*RecruiterCandidate) TableName() string {
	return "recruiter_candidate"
}

func GetRecruiterCandidates(userID int64) ([]*RecruiterCandidate, error) {
	var candidates []*RecruiterCandidate
	err := GetDB().Where("user_id = ?", userID).Find(&candidates).Error
	if err != nil {
		return nil, err
	}
	return candidates, nil
}

func SaveRecruiterCandidates(userID int64, candidates []*RecruiterCandidate) error {
	if err := GetDB().Where("user_id = ?", userID).Delete(&RecruiterCandidate{}).Error; err != nil {
		log.Errorf("error deleting candidates: %v ", err)
		return err
	}
	if len(candidates) > 0 {
		if err := GetDB().Create(candidates).Error; err != nil {
			log.Errorf("error creating candidates: %v ", err)
			return err
		}
	}
	return nil
}
