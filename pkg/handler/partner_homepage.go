package handler

import (
	"encoding/json"
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

type PartnerHomepageHandler struct {
}

func NewPartnerHomepageHandler() *PartnerHomepageHandler {
	return &PartnerHomepageHandler{}
}

func (h *PartnerHomepageHandler) Get(c *gin.Context) {
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
	audience := c.Query("audience")
	if audience != "talent" && audience != "employer" {
		log.Errorf("invalid homepage audience: %v ", audience)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	homepage, err1 := mysql.GetPartnerHomepage(userID, audience)
	if err1 != nil {
		if errors.Is(err1, gorm.ErrRecordNotFound) {
			h.genErrResponse(c, model.ErrCodeProfileNotFound)
			return
		}
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": homepage})
}

func (h *PartnerHomepageHandler) Update(c *gin.Context) {
	var req UpdatePartnerHomepageRequest
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("UpdatePartnerHomepage bind request error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	log.Infof("UpdatePartnerHomepage request: %v ", req)
	user := middleware.GetUser(c)
	if user.UserID != req.UserID {
		log.Errorf("UpdatePartnerHomepage requested user not logged in")
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	if req.Audience != "talent" && req.Audience != "employer" {
		log.Errorf("invalid homepage audience: %v ", req.Audience)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	updates := &mysql.PartnerHomepage{
		UserID:      user.UserID,
		Audience:    req.Audience,
		DisplayName: req.DisplayName,
		Summary:     req.Summary,
		CTPSummary:  req.CTPSummary,
		Companies:   req.Companies,
		DataPolicy:  req.DataPolicy,
		Applicants:  req.Applicants,
		HowTo:       req.HowTo,
	}
	if len(req.Reasons) > 0 {
		reasons, err := json.Marshal(req.Reasons)
		if err != nil {
			log.Errorf("UpdatePartnerHomepage marshal reasons error:%v", err)
			h.genErrResponse(c, model.ErrCodeInternalServerError)
			return
		}
		updates.Reasons = string(reasons)
	}
	if err := mysql.UpdatePartnerHomepage(user.UserID, req.Audience, updates); err != nil {
		log.Errorf("UpdatePartnerHomepage err:%v", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "homepage updated"})
}

func (h *PartnerHomepageHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}
