package handler

import (
	"github.com/gin-gonic/gin"
	"landd.co/landd/pkg/model"
)

type PartnerProfileHandler struct {
}

func NewPartnerProfileHandler() *PartnerProfileHandler {
	return &PartnerProfileHandler{}
}

func (h *PartnerProfileHandler) Get(c *gin.Context) {

}

func (h *PartnerProfileHandler) Update(c *gin.Context) {

}

func (h *PartnerProfileHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}
