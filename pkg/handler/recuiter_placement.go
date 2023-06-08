package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"landd.co/landd/pkg/middleware"
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/mysql"
	"net/http"
	"strconv"
	"time"
)

type RecruiterPlacementHandler struct {
}

func NewRecruiterPlacementHandler() *RecruiterPlacementHandler {
	return &RecruiterPlacementHandler{}
}

func (h *RecruiterPlacementHandler) GetPlacements(c *gin.Context) {
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
	placements, mysqlErr := mysql.GetRecruiterPlacements(middleware.GetUser(c).UserID)
	if mysqlErr != nil {
		log.Errorf("GetRecruiterPlacements error: %v ", mysqlErr)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": placements})
}

type UpdatePlacementRequest struct {
	UserID     int64        `json:"user_id"`
	Placements []*Placement `json:"placements"`
}

type Placement struct {
	Date     int64  `json:"date"`
	Position string `json:"position"`
	Company  string `json:"company"`
	Verified bool   `json:"verified"`
}

func (h *RecruiterPlacementHandler) UpdatePlacements(c *gin.Context) {
	var req UpdatePlacementRequest
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("UpdatePlacements bind request error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	user := middleware.GetUser(c)
	if user.UserID != req.UserID {
		log.Errorf("UpdatePlacements requested user not logged in")
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	log.Infof("UpdatePlacements request: %v ", req)
	var mysqlPlacements []*mysql.RecruiterPlacement
	for _, p := range req.Placements {
		placement := &mysql.RecruiterPlacement{
			UserID:   user.UserID,
			Date:     time.Unix(p.Date, 0),
			Position: p.Position,
			Company:  p.Position,
			Verified: p.Verified,
		}
		mysqlPlacements = append(mysqlPlacements, placement)
	}
	if err := mysql.SaveRecruiterPlacements(user.UserID, mysqlPlacements); err != nil {
		log.Errorf("SaveRecruiterPlacements mysql err:%v", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "placements updated"})
}

func (h *RecruiterPlacementHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}
