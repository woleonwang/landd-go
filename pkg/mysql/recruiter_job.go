package mysql

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type RecruiterJob struct {
	ID          int64     `gorm:"column:id;primarykey" json:"id"`
	UserID      int64     `gorm:"column:user_id" json:"user_id"`
	Title       string    `gorm:"column:title" json:"title"`
	Company     string    `gorm:"column:company" json:"company"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (*RecruiterJob) TableName() string {
	return "recruiter_job"
}

func GetRecruiterJobs(userID int64) ([]*RecruiterJob, error) {
	var jobs []*RecruiterJob
	err := GetDB().Where("user_id = ?", userID).Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func SaveRecruiterJobs(userID int64, jobs []*RecruiterJob) error {
	if err := GetDB().Where("user_id = ?", userID).Delete(&RecruiterJob{}).Error; err != nil {
		log.Errorf("error deleting jobs: %v ", err)
		return err
	}
	if err := GetDB().Create(jobs).Error; err != nil {
		log.Errorf("error creating user jobs: %v ", err)
		return err
	}
	return nil
}
