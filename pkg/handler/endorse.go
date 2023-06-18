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

func (h *EndorseHandler) Get(c *gin.Context) {
	user, err := h.validateUser(c)
	if err != nil {
		h.genErrResponse(c, *err)
		return
	}
	inviteIDStr := c.Query("invite_id")
	if inviteIDStr != "" {
		inviteID, err1 := strconv.ParseInt(inviteIDStr, 10, 64)
		if err1 != nil {
			log.Errorf("convert request invite_id error: %v ", err1)
			h.genErrResponse(c, model.ErrCodeInvalidRequest)
			return
		}
		h.GetByInviteID(c, user.UserID, inviteID)
		return
	}
	statusList := c.QueryArray("status")
	var statuses []model.EndorsementStatus
	for _, s := range statusList {
		status, err := strconv.Atoi(s)
		if err != nil {
			log.Errorf("convert request status error: %v ", err)
			h.genErrResponse(c, model.ErrCodeInvalidRequest)
			return
		}
		statuses = append(statuses, model.EndorsementStatus(status))
	}
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	if size > 100 || size <= 0 || page <= 0 {
		log.Errorf("page size invalid")
		h.genErrResponse(c, model.ErrCodePageSizeInvalid)
		return
	}
	endorsements, err2 := mysql.GetEndorsements(user.UserID, statuses, (page-1)*size, size)
	if err2 != nil {
		log.Errorf("GetEndorsements err:%v", err2)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	resp := GetEndorsementResponse{
		UserID:       user.UserID,
		Endorsements: endorsements,
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

func (h *EndorseHandler) GetByInviteID(c *gin.Context, userID, inviteID int64) {
	endorsement, err := mysql.GetEndorsementByID(userID, inviteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("endorsement not found")
			h.genErrResponse(c, model.ErrCodeInvalidRequest)
			return
		}
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	resp := GetEndorsementResponse{
		UserID:       userID,
		Endorsements: []*mysql.Endorsement{endorsement},
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

// Invite takes current draft and creates a new endorsement with unique invite_id
func (h *EndorseHandler) Invite(c *gin.Context) {
	user, err := h.validateUser(c)
	if err != nil {
		h.genErrResponse(c, *err)
		return
	}
	draft, err1 := mysql.GetEndorsementDraft(user.UserID)
	if err1 != nil {
		log.Errorf("GetEndorsementDraft error: %v ", err1)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	endorsement := &mysql.Endorsement{
		UserID:       draft.UserID,
		EndorserName: draft.EndorserName,
		Title:        draft.Title,
		Company:      draft.Company,
		Identity:     draft.Identity,
		Status:       model.EndorsementStatusInvited,
		Content:      draft.Content,
	}
	inviteID, err2 := mysql.CreateEndorsement(endorsement)
	if err2 != nil {
		log.Errorf("CreateEndorsement error: %v ", err2)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	resp := InviteEndorseResponse{
		UserID:   user.UserID,
		InviteID: inviteID,
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

func (h *EndorseHandler) Update(c *gin.Context) {
	var req UpdateEndorsementRequest
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("Update bind request error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	log.Infof("Update request: %v ", req)
	user := middleware.GetUser(c)
	if user.UserID != req.UserID {
		log.Errorf("Operate requested user not logged in")
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	updates := &mysql.Endorsement{
		InviteID:     req.InviteID,
		UserID:       req.UserID,
		EndorserName: req.Endorser,
		Title:        req.Title,
		Company:      req.Company,
		Identity:     req.Identity,
		Status:       req.Status,
		Content:      req.Content,
	}
	if err := mysql.UpdateEndorsement(req.UserID, req.InviteID, updates); err != nil {
		log.Errorf("UpdateEndorsement error: %v ", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "update successful"})
}

func (h *EndorseHandler) GetDraft(c *gin.Context) {
	user, err := h.validateUser(c)
	if err != nil {
		h.genErrResponse(c, *err)
		return
	}
	draft, err1 := mysql.GetEndorsementDraft(user.UserID)
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

func (h *EndorseHandler) validateUser(c *gin.Context) (*mysql.User, *model.CustomError) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Errorf("convert request user_id error: %v ", err)
		return nil, &model.ErrCodeInvalidRequest
	}
	user := middleware.GetUser(c)
	if user.UserID != userID {
		return nil, &model.ErrCodeRequestUserNotLogin
	}
	return &user, nil
}

func (h *EndorseHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}
