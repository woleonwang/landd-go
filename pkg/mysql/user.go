package mysql

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/utils"
	"time"
)

type User struct {
	ID        int64          `gorm:"column:id;primarykey" json:"id"`
	UserID    int64          `gorm:"column:user_id" json:"user_id"`
	Name      string         `gorm:"column:name" json:"name"`
	Email     string         `gorm:"column:email" json:"email"`
	Password  string         `gorm:"column:password" json:"password"`
	Job       string         `gorm:"column:job" json:"job"`
	Role      model.UserRole `gorm:"column:role" json:"role"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (*User) TableName() string {
	return "user"
}

func GetUserByEmail(email string) (*User, error) {
	user := User{}
	err := GetDB().Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *User) (int64, error) {
	err := GetDB().Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.UserID, nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	id, err := utils.GenerateID()
	if err != nil {
		log.Errorf("error creating user id: %v ", err)
		return err
	}
	u.UserID = id
	return nil
}
