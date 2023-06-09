package mysql

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type RecruiterPlacement struct {
	ID        int64     `gorm:"column:id;primarykey" json:"id"`
	UserID    int64     `gorm:"column:user_id" json:"user_id"`
	Date      time.Time `gorm:"column:date" json:"date"`
	Position  string    `gorm:"column:position" json:"position"`
	Company   string    `gorm:"column:company" json:"company"`
	Verified  bool      `gorm:"column:verified" json:"verified"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (*RecruiterPlacement) TableName() string {
	return "recruiter_placement"
}

func GetRecruiterPlacements(userID int64) ([]*RecruiterPlacement, error) {
	var placements []*RecruiterPlacement
	err := GetDB().Where("user_id = ?", userID).Order("date desc").Find(&placements).Error
	if err != nil {
		return nil, err
	}
	return placements, nil
}

func SaveRecruiterPlacements(userID int64, placements []*RecruiterPlacement) error {
	if err := GetDB().Where("user_id = ?", userID).Delete(&RecruiterPlacement{}).Error; err != nil {
		log.Errorf("error deleting placements: %v ", err)
		return err
	}
	if len(placements) > 0 {
		if err := GetDB().Create(placements).Error; err != nil {
			log.Errorf("error creating user placements: %v ", err)
			return err
		}
	}
	return nil
}
