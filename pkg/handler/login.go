package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"landd.co/landd/pkg/middleware"
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/mysql"
	"net/http"
	"strings"
)

type LoginHandler struct{}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

type SignUpRequest struct {
	IsEmailCheck bool                         `json:"check_email"`
	Email        string                       `json:"email"`
	Password     string                       `json:"password"`
	Name         string                       `json:"name"`
	Job          string                       `json:"job"`
	Contacts     map[model.ContactType]string `json:"contacts"`
}

func (h *LoginHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.BindJSON(&req); err != nil {
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	log.Infof("signing up .. does email check only: %v, email: %v", req.IsEmailCheck, req.Email)
	_, mysqlErr := mysql.GetUserByEmail(req.Email)
	if mysqlErr == nil {
		log.Errorf("GetUserByEmail error: %v ", mysqlErr)
		h.genErrResponse(c, model.ErrCodeDuplicateEmail)
		return
	}
	if !errors.Is(mysqlErr, gorm.ErrRecordNotFound) {
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	if req.IsEmailCheck {
		c.JSON(http.StatusOK, gin.H{"message": ""})
		return
	}
	user := &mysql.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Job:      req.Job,
	}
	userID, createErr := mysql.CreateUser(user)
	if createErr != nil {
		log.Errorf("CreateUser error: %v ", createErr)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	var contacts []*mysql.UserContact
	for t, contact := range req.Contacts {
		userContact := &mysql.UserContact{
			UserID:  userID,
			Contact: contact,
			Type:    t,
		}
		contacts = append(contacts, userContact)
	}
	if err := mysql.CreateUserContact(contacts); err != nil {
		log.Errorf("CreateUserContact error: %v ", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	if err := middleware.SetUser(c, *user); err != nil {
		log.Errorf("middleware.SetUser error: %v ", err)
		h.genErrResponse(c, model.ErrCodeSessionError)
		return
	}
	log.Infof("user created with id: %v", userID)
	resp := LoginResponse{
		UserID:   userID,
		UserName: req.Name,
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"username"`
}

func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	if strings.Trim(req.Email, " ") == "" || strings.Trim(req.Password, " ") == "" {
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	user, mysqlErr := mysql.GetUserByEmail(req.Email)
	if mysqlErr != nil {
		if errors.Is(mysqlErr, gorm.ErrRecordNotFound) {
			h.genErrResponse(c, model.ErrCodeEmailNotFound)
		} else {
			h.genErrResponse(c, model.ErrCodeMysqlError)
		}
		return
	}
	if user.Password != req.Password {
		h.genErrResponse(c, model.ErrCodeIncorrectCredential)
		return
	}
	if err := middleware.SetUser(c, *user); err != nil {
		log.Errorf("middleware.SetUser error: %v ", err)
		h.genErrResponse(c, model.ErrCodeSessionError)
		return
	}
	resp := LoginResponse{
		UserID:   user.UserID,
		UserName: user.Name,
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

func (h *LoginHandler) Logout(c *gin.Context) {
	if err := middleware.ClearSession(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *LoginHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}
