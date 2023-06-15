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

type EndorseHandler struct {
}

func NewEndorseHandler() *EndorseHandler {
	return &EndorseHandler{}
}

func (h *EndorseHandler) GetEndorsement(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Errorf("convert request user_id error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	user := middleware.GetUser(c)
	if user.UserID != userID {
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	statusStr := c.Query("status")
	endorseStatus, err1 := strconv.Atoi(statusStr)
	if err1 != nil {
		log.Errorf("convert requested status error: %v ", err1)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	endorsements, err2 := mysql.GetEndorsements(userID, endorseStatus)
	if err2 != nil {
		log.Errorf("GetEndorsements error: %v ", err2)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	resp := GetEndorsementResponse{
		UserID:       userID,
		Endorsements: endorsements,
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

func (h *EndorseHandler) UpdateEndorsement(c *gin.Context) {

}

func (h *EndorseHandler) GetDraft(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Errorf("convert request user_id error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	user := middleware.GetUser(c)
	if user.UserID != userID {
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	draft, err1 := mysql.GetEndorsementDraft(userID)
	if err1 != nil && !errors.Is(err1, gorm.ErrRecordNotFound) {
		log.Errorf("GetEndorsementDraft error: %v ", err1)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": draft})
}

func (h *EndorseHandler) UpdateDraft(c *gin.Context) {
	var req UpdateEndorseDraftRequest
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("UpdateDraft bind request error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	log.Infof("UpdateDraft request: %v ", req)
	user := middleware.GetUser(c)
	if user.UserID != req.UserID {
		log.Errorf("UpdateDraft requested user not logged in")
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	updates := &mysql.EndorsementDraft{
		UserID:       req.UserID,
		EndorserName: req.Endorser,
		Title:        req.Title,
		Company:      req.Company,
		Identity:     req.Identity,
		Content:      req.Content,
	}
	if err := mysql.UpdateEndorsementDraft(req.UserID, updates); err != nil {
		log.Errorf("UpdateEndorsementDraft err:%v", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "endorsement draft updated"})
}

func (h *EndorseHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}
