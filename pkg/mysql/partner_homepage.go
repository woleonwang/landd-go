package mysql

import (
	"time"
)

type PartnerHomepage struct {
	ID          int64     `gorm:"column:id;primarykey" json:"id"`
	UserID      int64     `gorm:"column:user_id" json:"user_id"`
	Audience    string    `gorm:"column:audience" json:"audience"`
	DisplayName string    `gorm:"column:display_name" json:"display_name"`
	Summary     string    `gorm:"column:summary" json:"summary"`
	CTPSummary  string    `gorm:"column:ctp_summary" json:"ctp_summary"`
	Reasons     string    `gorm:"column:reasons" json:"reasons"`
	Companies   string    `gorm:"column:companies" json:"companies"`
	DataPolicy  string    `gorm:"column:data_policy" json:"data_policy"`
	Applicants  string    `gorm:"column:applicants" json:"applicants"`
	HowTo       string    `gorm:"column:howto" json:"howto"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (*PartnerHomepage) TableName() string {
	return "partner_homepage"
}

func GetPartnerHomepage(userID int64, audience string) (*PartnerHomepage, error) {
	var page PartnerHomepage
	err := GetDB().Where("user_id = ? and audience = ?", userID, audience).
		First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func UpdatePartnerHomepage(userID int64, audience string, updates *PartnerHomepage) error {
	res := GetDB().Model(&PartnerHomepage{}).
		Where("user_id = ? and audience = ?", userID, audience).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return GetDB().Create(updates).Error
	}
	return nil
}
