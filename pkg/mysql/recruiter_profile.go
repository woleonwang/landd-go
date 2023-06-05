package mysql

import (
	"time"
)

type RecruiterProfile struct {
	ID                    int64     `gorm:"column:id;primarykey" json:"id"`
	UserID                int64     `gorm:"column:user_id" json:"user_id"`
	Name                  string    `gorm:"column:name" json:"name"`
	Photo                 string    `gorm:"column:photo" json:"photo"`
	Summary               string    `gorm:"column:summary" json:"summary"`
	Company               string    `gorm:"column:company" json:"company"`
	YearsExpr             int       `gorm:"column:years_of_expr" json:"years_of_expr"`
	Expertise             string    `gorm:"column:expertise" json:"expertise"`
	TotalPlacedCandidates int       `gorm:"column:total_placed_candidates" json:"total_placed_candidates"`
	TotalPlacedSalary     int64     `gorm:"column:total_placed_salary" json:"total_placed_salary"`
	TotalCandidates       int       `gorm:"column:total_candidates" json:"total_candidates"`
	CreatedAt             time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (*RecruiterProfile) TableName() string {
	return "recruiter_profile"
}

func GetRecruiterProfile(userID int64) (*RecruiterProfile, error) {
	profile := RecruiterProfile{}
	err := GetDB().Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func SaveRecruiterProfile(profile *RecruiterProfile) error {
	err := GetDB().Save(profile).Error
	if err != nil {
		return err
	}
	return nil
}
