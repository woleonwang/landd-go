package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"landd.co/landd/pkg/middleware"
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/mysql"
	"net/http"
	"strconv"
	"strings"
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
	profile, mysqlErr := mysql.GetRecruiterProfile(userID)
	if mysqlErr != nil {
		if errors.Is(mysqlErr, gorm.ErrRecordNotFound) {
			log.Warnf("GetRecruiterProfile not found: %v ", err)
			h.genErrResponse(c, model.ErrCodeProfileNotFound)
			return
		}
		log.Errorf("GetRecruiterProfile error: %v ", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": profile})
}

type UpdateProfileInfoRequest struct {
	UserID                int64  `json:"user_id"`
	Name                  string `json:"name"`
	Photo                 string `json:"photo"`
	Summary               string `json:"summary"`
	Company               string `json:"company"`
	YearsExpr             int    `json:"years_of_expr"`
	Expertise             string `json:"expertise"`
	TotalPlacedCandidates int    `json:"total_placed_candidates"`
	TotalPlacedSalary     int64  `json:"total_placed_salary"`
	TotalCandidates       int    `json:"total_candidates"`
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
	updates := &mysql.RecruiterProfile{
		UserID:                user.UserID,
		Name:                  req.Name,
		Photo:                 req.Photo,
		Summary:               req.Summary,
		Company:               req.Company,
		YearsExpr:             req.YearsExpr,
		Expertise:             req.Expertise,
		TotalPlacedCandidates: req.TotalPlacedCandidates,
		TotalPlacedSalary:     req.TotalPlacedSalary,
		TotalCandidates:       req.TotalCandidates,
	}
	if err := mysql.UpdateRecruiterProfile(user.UserID, updates); err != nil {
		log.Errorf("UpdateRecruiterProfile mysql err:%v", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

type UploadPhotoResponse struct {
	PhotoID string `json:"photo_id"`
}

func (h *RecruiterProfileHandler) UploadPhoto(c *gin.Context) {
	// TODO save images to cloud storage
	file, err := c.FormFile("image")
	if err != nil {
		log.Errorf("image upload error: %v ", err)
		h.genErrResponse(c, model.ErrCodeImageInvalid)
		return
	}
	filename := uuid.New().String()
	filenameParts := strings.Split(file.Filename, ".")
	if len(filenameParts) < 2 {
		log.Errorf("file name invalid")
		h.genErrResponse(c, model.ErrCodeImageInvalid)
		return
	}
	image := fmt.Sprintf("%s.%s", filename, filenameParts[len(filenameParts)-1])

	if err = c.SaveUploadedFile(file, fmt.Sprintf("./images/%s", image)); err != nil {
		log.Errorf("image save error: %v ", err)
		h.genErrResponse(c, model.ErrCodeImageUploadError)
		return
	}
	resp := UploadPhotoResponse{
		PhotoID: image,
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

func (h *RecruiterProfileHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}
