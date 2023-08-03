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
	"strconv"
)

type PartnerProfileHandler struct {
}

func NewPartnerProfileHandler() *PartnerProfileHandler {
	return &PartnerProfileHandler{}
}

func (h *PartnerProfileHandler) Get(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Errorf("convert request param error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	user := middleware.GetUser(c)
	if user.UserID != userID {
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	profile, err1 := mysql.GetPartnerProfile(userID)
	if err1 != nil {
		if errors.Is(err1, gorm.ErrRecordNotFound) {
			h.genErrResponse(c, model.ErrCodeProfileNotFound)
			return
		}
		log.Errorf("GetPartnerProfile err: %v ", err1)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	contacts, err2 := mysql.GetUserContacts(userID)
	if err2 != nil {
		log.Errorf("GetUserContacts err: %v ", err2)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	resp := GetPartnerProfileResponse{
		UserID:   strconv.FormatInt(userID, 10),
		Profile:  profile,
		Email:    user.Email,
		Contacts: contacts,
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

func (h *PartnerProfileHandler) Update(c *gin.Context) {
	var req UpdatePartnerProfileRequest
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("UpdatePartnerProfile bind request error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	log.Infof("UpdatePartnerProfile request: %v ", req)
	user := middleware.GetUser(c)
	if user.UserID != req.UserID {
		log.Errorf("UpdatePartnerProfile requested user not logged in")
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	profileChanges := &mysql.PartnerProfile{
		UserID:    req.UserID,
		Name:      req.Name,
		Photo:     req.Photo,
		Company:   req.Company,
		Twitter:   req.Twitter,
		LinkedIn:  req.LinkedIn,
		Website:   req.Website,
		Blog:      req.Blog,
		Facebook:  req.Facebook,
		Instagram: req.Instagram,
		Tiktok:    req.Tiktok,
		Youtube:   req.Youtube,
		Other:     req.Other,
	}
	if err := mysql.UpdatePartnerProfile(user.UserID, profileChanges); err != nil {
		log.Errorf("UpdatePartnerProfile err:%v", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	var contactChanges []*mysql.UserContact
	for t, contact := range req.Contacts {
		userContact := &mysql.UserContact{
			UserID:  user.UserID,
			Contact: contact,
			Type:    t,
		}
		contactChanges = append(contactChanges, userContact)
	}
	if err := mysql.SaveUserContacts(user.UserID, contactChanges); err != nil {
		log.Errorf("SaveUserContacts error: %v ", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "update successful"})
}

func (h *PartnerProfileHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}
