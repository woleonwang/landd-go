package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"landd.co/landd/pkg/middleware"
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/mysql"
)

type RecruiterProfileHandler struct {
}

func NewRecruiterProfileHandler() *RecruiterProfileHandler {
	return &RecruiterProfileHandler{}
}

func (h *RecruiterProfileHandler) GetProfileInfo(c *gin.Context) {
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
	profile, err1 := mysql.GetRecruiterProfile(userID)
	if err1 != nil {
		if errors.Is(err1, gorm.ErrRecordNotFound) {
			log.Warnf("GetRecruiterProfile not found: %v ", err1)
			h.genErrResponse(c, model.ErrCodeProfileNotFound)
			return
		}
		log.Errorf("GetRecruiterProfile error: %v ", err1)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	placements, err2 := mysql.GetRecruiterPlacements(userID)
	if err2 != nil {
		log.Errorf("GetRecruiterPlacements error: %v ", err2)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	jobs, err3 := mysql.GetRecruiterJobs(userID)
	if err3 != nil {
		log.Errorf("GetRecruiterJobs error: %v ", err3)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	candidates, err4 := mysql.GetRecruiterCandidates(userID)
	if err4 != nil {
		log.Errorf("GetRecruiterCandidates error: %v ", err4)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	pubs, err5 := mysql.GetRecruiterPublication(userID)
	if err5 != nil {
		log.Errorf("GetRecruiterPublication error: %v ", err5)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	resp := GetProfileInfoResponse{
		UserID:       userID,
		Profile:      profile,
		Placements:   placements,
		Jobs:         jobs,
		Candidates:   candidates,
		Publications: pubs,
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

func (h *RecruiterProfileHandler) UpdateProfileInfo(c *gin.Context) {
	var req UpdateProfileInfoRequest
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("UpdateProfileInfo bind request error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	log.Infof("UpdateProfileInfo request: %v ", req)
	user := middleware.GetUser(c)
	if user.UserID != req.UserID {
		log.Errorf("UpdateProfileInfo requested user not logged in")
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	if req.ProfileChanges != nil {
		if err := h.updateProfile(req); err != nil {
			log.Errorf("updateProfile err:%v", err)
			h.genErrResponse(c, model.ErrCodeMysqlError)
			return
		}
	}
	if req.PlacementChanges != nil {
		if err := h.updatePlacement(req); err != nil {
			log.Errorf("updatePlacement err:%v", err)
			h.genErrResponse(c, model.ErrCodeMysqlError)
			return
		}
	}
	if req.JobChanges != nil {
		if err := h.updateJob(req); err != nil {
			log.Errorf("updateJob err:%v", err)
			h.genErrResponse(c, model.ErrCodeMysqlError)
			return
		}
	}
	if req.CandidateChanges != nil {
		if err := h.updateCandidate(req); err != nil {
			log.Errorf("updateCandidate err:%v", err)
			h.genErrResponse(c, model.ErrCodeMysqlError)
			return
		}
	}
	if req.PublicationChanges != nil {
		if err := h.updatePublication(req); err != nil {
			log.Errorf("updatePublication err:%v", err)
			h.genErrResponse(c, model.ErrCodeMysqlError)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

func (h *RecruiterProfileHandler) updateProfile(req UpdateProfileInfoRequest) error {
	updates := &mysql.RecruiterProfile{
		UserID:                req.UserID,
		Name:                  req.ProfileChanges.Name,
		Photo:                 req.ProfileChanges.Photo,
		Summary:               req.ProfileChanges.Summary,
		Company:               req.ProfileChanges.Company,
		YearsExpr:             req.ProfileChanges.YearsExpr,
		Expertise:             req.ProfileChanges.Expertise,
		TotalPlacedCandidates: req.ProfileChanges.TotalPlacedCandidates,
		TotalPlacedSalary:     req.ProfileChanges.TotalPlacedSalary,
		TotalCandidates:       req.ProfileChanges.TotalCandidates,
	}
	return mysql.UpdateRecruiterProfile(req.UserID, updates)
}

func (h *RecruiterProfileHandler) updatePlacement(req UpdateProfileInfoRequest) error {
	var mysqlPlacements []*mysql.RecruiterPlacement
	for _, p := range req.PlacementChanges {
		placement := &mysql.RecruiterPlacement{
			UserID:   req.UserID,
			Date:     time.Unix(p.Date, 0),
			Position: p.Position,
			Company:  p.Company,
			Verified: p.Verified,
		}
		mysqlPlacements = append(mysqlPlacements, placement)
	}
	return mysql.SaveRecruiterPlacements(req.UserID, mysqlPlacements)
}

func (h *RecruiterProfileHandler) updateJob(req UpdateProfileInfoRequest) error {
	var jobs []*mysql.RecruiterJob
	for _, j := range req.JobChanges {
		job := &mysql.RecruiterJob{
			UserID:      req.UserID,
			Title:       j.Title,
			Company:     j.Company,
			Description: j.Description,
		}
		jobs = append(jobs, job)
	}
	return mysql.SaveRecruiterJobs(req.UserID, jobs)
}

func (h *RecruiterProfileHandler) updateCandidate(req UpdateProfileInfoRequest) error {
	var candidates []*mysql.RecruiterCandidate
	for _, c := range req.CandidateChanges {
		candidate := &mysql.RecruiterCandidate{
			UserID:     req.UserID,
			Title:      c.Title,
			Percentage: c.Percentage,
		}
		candidates = append(candidates, candidate)
	}
	return mysql.SaveRecruiterCandidates(req.UserID, candidates)
}

func (h *RecruiterProfileHandler) updatePublication(req UpdateProfileInfoRequest) error {
	var pubs []*mysql.RecruiterPublication
	for _, p := range req.PublicationChanges {
		pub := &mysql.RecruiterPublication{
			UserID: req.UserID,
			Title:  p.Title,
			Link:   p.Link,
		}
		pubs = append(pubs, pub)
	}
	return mysql.SaveRecruiterPublication(req.UserID, pubs)
}

func (h *RecruiterProfileHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}
