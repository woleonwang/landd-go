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

type FileHandler struct {
}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

type UploadFileResponse struct {
	FileID string `json:"file_id"`
}

func (h *FileHandler) Upload(c *gin.Context) {
	// TODO save files to cloud storage
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("file upload error: %v ", err)
		h.genErrResponse(c, model.ErrCodeFileInvalid)
		return
	}
	filename := uuid.New().String()
	filenameParts := strings.Split(file.Filename, ".")
	if len(filenameParts) < 2 {
		log.Errorf("file name invalid")
		h.genErrResponse(c, model.ErrCodeFileInvalid)
		return
	}
	fileID := fmt.Sprintf("%s.%s", filename, filenameParts[len(filenameParts)-1])

	if err = c.SaveUploadedFile(file, fmt.Sprintf("./files/%s", fileID)); err != nil {
		log.Errorf("file save error: %v ", err)
		h.genErrResponse(c, model.ErrCodeImageUploadError)
		return
	}
	resp := UploadFileResponse{
		FileID: fileID,
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

func (h *FileHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}
