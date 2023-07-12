package mysql

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"landd.co/landd/pkg/utils"
	"time"
)

type Job struct {
	ID            int64     `gorm:"column:id;primarykey" json:"id"`
	JobID         int64     `gorm:"column:job_id" json:"job_id"`
	Title         string    `gorm:"column:title" json:"title"`
	Company       string    `gorm:"column:company" json:"company"`
	Jd            string    `gorm:"column:jd" json:"jd"`
	AboutCompany  string    `gorm:"column:about_company" json:"about_company"`
	Comment       string    `gorm:"column:comment" json:"comment"`
	ReferralFee   int       `gorm:"column:referral_fee" json:"referral_fee"`
	LowerBoundSal int       `gorm:"column:lower_bound_sal" json:"lower_bound_sal"`
	UpperBoundSal int       `gorm:"column:upper_bound_sal" json:"upper_bound_sal"`
	PosterID      int64     `gorm:"column:poster_id" json:"poster_id"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (*Job) TableName() string {
	return "job"
}

func (u *Job) BeforeCreate(tx *gorm.DB) error {
	id, err := utils.GenerateID()
	if err != nil {
		log.Errorf("error creating job id: %v ", err)
		return err
	}
	u.JobID = id
	return nil
}

func CreateJobOpening(job *Job) (int64, error) {
	err := GetDB().Create(job).Error
	if err != nil {
		return 0, err
	}
	return job.JobID, nil
}

func GetJobOpenings(offset, limit int) ([]*Job, error) {
	var jobs []*Job
	err := GetDB().Offset(offset).Limit(limit).Order("created_at desc").Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func GetJobByID(jobID int64) (*Job, error) {
	job := Job{}
	err := GetDB().Where("job_id = ?", jobID).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func UpdateJobOpening(jobID int64, updates *Job) error {
	res := GetDB().Model(&Job{}).Where("job_id = ?", jobID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
