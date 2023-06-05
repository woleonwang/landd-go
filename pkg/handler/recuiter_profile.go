package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"landd.co/landd/pkg/model"
	"net/http"
	"strings"
)

type RecruiterProfileHandler struct {
}

func NewRecruiterProfileHandler() *RecruiterProfileHandler {
	return &RecruiterProfileHandler{}
}

func (h *RecruiterProfileHandler) GetProfileInfo(c *gin.Context) {

}

func (h *RecruiterProfileHandler) UpdateProfileInfo(c *gin.Context) {

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
