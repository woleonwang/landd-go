package mysql

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/utils"
	"time"
)

type PartnerCandidate struct {
	ID          int64           `gorm:"column:id;primarykey" json:"id"`
	UserID      int64           `gorm:"column:user_id" json:"user_id"`
	CandidateID int64           `gorm:"column:candidate_id" json:"candidate_id"`
	Status      model.CTPStatus `gorm:"column:status" json:"status"`
	Vet         model.Vet       `gorm:"column:vet" json:"vet"`
	FirstName   string          `gorm:"column:first_name" json:"first_name"`
	LastName    string          `gorm:"column:last_name" json:"last_name"`
	Mobile      string          `gorm:"column:mobile" json:"mobile"`
	Email       string          `gorm:"column:email" json:"email"`
	Expr        int             `gorm:"column:expr" json:"expr"`
	LinkedIn    string          `gorm:"column:linkedin" json:"linkedin"`
	Resume      string          `gorm:"column:resume" json:"resume"`
	WorkExpr    string          `gorm:"column:work_expr" json:"work_expr"`
	Education   string          `gorm:"column:education" json:"education"`
	Skill       string          `gorm:"column:skill" json:"skill"`
	Tag         string          `gorm:"column:tag" json:"tag"`
	Comment     string          `gorm:"column:comment" json:"comment"`
	Note        string          `gorm:"column:note" json:"note"`
	CreatedAt   time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at" json:"updated_at"`
}

func (*PartnerCandidate) TableName() string {
	return "partner_candidate"
}

func CreateCTPCandidate(candidate *PartnerCandidate) (int64, error) {
	err := GetDB().Create(candidate).Error
	if err != nil {
		return 0, err
	}
	return candidate.CandidateID, nil
}

func GetCTPCandidates(userID int64, statuses []model.CTPStatus, offset, limit int) (
	[]*PartnerCandidate, error) {
	var candidates []*PartnerCandidate
	query := GetDB().Where("user_id = ?", userID)
	if len(statuses) > 0 {
		query = query.Where("status IN (?)", statuses)
	}
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&candidates).Error
	if err != nil {
		return nil, err
	}
	return candidates, nil
}

func GetCandidateByID(userID, candidateID int64) (*PartnerCandidate, error) {
	candidate := PartnerCandidate{}
	err := GetDB().Where("user_id = ? and candidate_id = ?", userID, candidateID).
		First(&candidate).Error
	if err != nil {
		return nil, err
	}
	return &candidate, nil
}

func UpdateCTPCandidate(userID, candidateID int64, updates *PartnerCandidate) error {
	res := GetDB().Model(&PartnerCandidate{}).
		Where("user_id = ? and candidate_id = ?", userID, candidateID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (u *PartnerCandidate) BeforeCreate(tx *gorm.DB) error {
	id, err := utils.GenerateID()
	if err != nil {
		log.Errorf("error creating candidate id: %v ", err)
		return err
	}
	u.CandidateID = id
	return nil
}
