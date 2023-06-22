package mysql

import (
	"time"
)

type PartnerProfile struct {
	ID        int64     `gorm:"column:id;primarykey" json:"id"`
	UserID    int64     `gorm:"column:user_id" json:"user_id"`
	Name      string    `gorm:"column:name" json:"name"`
	Photo     string    `gorm:"column:photo" json:"photo"`
	Company   string    `gorm:"column:company" json:"company"`
	Channels  string    `gorm:"column:channels" json:"channels"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Channel struct {
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

func (*PartnerProfile) TableName() string {
	return "partner_profile"
}
