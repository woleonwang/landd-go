package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"landd.co/landd/pkg/middleware"
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/mysql"
	"net/http"
	"strconv"
)

type RecruiterJobHandler struct {
}

func NewRecruiterJobHandler() *RecruiterJobHandler {
	return &RecruiterJobHandler{}
}

func (h *RecruiterJobHandler) GetJobs(c *gin.Context) {
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
	jobs, mysqlErr := mysql.GetRecruiterJobs(middleware.GetUser(c).UserID)
	if mysqlErr != nil {
		log.Errorf("GetRecruiterJobs error: %v ", mysqlErr)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": jobs})
}

type UpdateJobRequest struct {
	UserID int64  `json:"user_id"`
	Jobs   []*Job `json:"jobs"`
}

type Job struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Description string `json:"description"`
}

func (h *RecruiterJobHandler) UpdateJobs(c *gin.Context) {
	var req UpdateJobRequest
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("UpdateJobRequest bind request error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	user := middleware.GetUser(c)
	if user.UserID != req.UserID {
		log.Errorf("UpdateJobs requested user not logged in")
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	log.Infof("UpdateJobs request: %v ", req)
	var jobs []*mysql.RecruiterJob
	for _, j := range req.Jobs {
		job := &mysql.RecruiterJob{
			UserID:      user.UserID,
			Title:       j.Title,
			Company:     j.Company,
			Description: j.Description,
		}
		jobs = append(jobs, job)
	}
	if err := mysql.SaveRecruiterJobs(user.UserID, jobs); err != nil {
		log.Errorf("SaveRecruiterJobs mysql err:%v", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "jobs updated"})
}

func (h *RecruiterJobHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}
