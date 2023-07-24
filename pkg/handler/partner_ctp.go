package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"landd.co/landd/pkg/middleware"
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/mysql"
	"net/http"
	"strconv"
)

type PartnerCTPHandler struct {
}

func NewPartnerCTPHandler() *PartnerCTPHandler {
	return &PartnerCTPHandler{}
}

func (h *PartnerCTPHandler) Get(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Errorf("convert request param error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	user := middleware.GetUser(c)
	if user.Role != model.Admin && user.UserID != userID {
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	candidateIDStr := c.Query("candidate_id")
	if candidateIDStr != "" {
		candidateID, err1 := strconv.ParseInt(candidateIDStr, 10, 64)
		if err1 != nil {
			log.Errorf("convert request candidate_id error: %v ", err1)
			h.genErrResponse(c, model.ErrCodeInvalidRequest)
			return
		}
		candidate, err := mysql.GetCandidateByID(userID, candidateID)
		if err != nil {
			log.Errorf("GetCandidateByID err:%v", err)
			h.genErrResponse(c, model.ErrCodeMysqlError)
			return
		}
		candidates, err2 := h.convertCandidates([]*mysql.PartnerCandidate{candidate})
		if err2 != nil {
			log.Errorf("convertCandidates err:%v", err2)
			h.genErrResponse(c, model.ErrCodeInternalServerError)
			return
		}
		resp := GetCTPCandidateResponse{
			UserID:     userID,
			Candidates: candidates,
		}
		c.JSON(http.StatusOK, gin.H{"message": resp})
		return
	}
	statusList := c.QueryArray("status")
	var statuses []model.CTPStatus
	for _, s := range statusList {
		status, err := strconv.Atoi(s)
		if err != nil {
			log.Errorf("convert request status error: %v ", err)
			h.genErrResponse(c, model.ErrCodeInvalidRequest)
			return
		}
		statuses = append(statuses, model.CTPStatus(status))
	}
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	if size > 100 || size <= 0 || page <= 0 {
		log.Errorf("page size invalid")
		h.genErrResponse(c, model.ErrCodePageSizeInvalid)
		return
	}
	candidates, err2 := mysql.GetCTPCandidates(userID, statuses, (page-1)*size, size)
	if err2 != nil {
		log.Errorf("GetCTPCandidates err:%v", err2)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	converted, err3 := h.convertCandidates(candidates)
	if err3 != nil {
		log.Errorf("convertCandidates err:%v", err3)
		h.genErrResponse(c, model.ErrCodeInternalServerError)
		return
	}
	resp := GetCTPCandidateResponse{
		UserID:     userID,
		Candidates: converted,
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

func (h *PartnerCTPHandler) Create(c *gin.Context) {
	var req CreateCTPCandidateRequest
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("CreateCTPCandidateRequest bind request error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	log.Infof("CreateCTPCandidateRequest request: %v ", req)
	user := middleware.GetUser(c)
	if user.UserID != req.UserID {
		log.Errorf("CreateCTPCandidateRequest requested user not logged in")
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	candidate := &mysql.PartnerCandidate{
		UserID:    req.UserID,
		Status:    model.CTPStatusPending,
		Vet:       req.Vet,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Mobile:    req.Mobile,
		Email:     req.Email,
		Expr:      req.Expr,
		LinkedIn:  req.LinkedIn,
		Skill:     req.Skill,
		Tag:       req.Tag,
		Comment:   req.Comment,
		Note:      req.Note,
	}
	resume, errJson := json.Marshal(req.Resume)
	if errJson != nil {
		log.Errorf("json marshal resume err:%v", errJson)
		h.genErrResponse(c, model.ErrCodeInternalServerError)
		return
	}
	workExpr, errJson := json.Marshal(req.WorkExpr)
	if errJson != nil {
		log.Errorf("json marshal workexpr err:%v", errJson)
		h.genErrResponse(c, model.ErrCodeInternalServerError)
		return
	}
	edu, errJson := json.Marshal(req.Education)
	if errJson != nil {
		log.Errorf("json marshal education err:%v", errJson)
		h.genErrResponse(c, model.ErrCodeInternalServerError)
		return
	}
	candidate.Resume = string(resume)
	candidate.WorkExpr = string(workExpr)
	candidate.Education = string(edu)
	candidateID, err := mysql.CreateCTPCandidate(candidate)
	if err != nil {
		log.Errorf("CreateCTPCandidate err:%v", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	resp := CreateCTPCandidateResponse{
		UserID:      user.UserID,
		CandidateID: candidateID,
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

func (h *PartnerCTPHandler) Update(c *gin.Context) {
	var req UpdateCTPCandidateRequest
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("UpdateCTPCandidateRequest bind request error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	log.Infof("UpdateCTPCandidateRequest request: %v ", req)
	user := middleware.GetUser(c)
	if user.UserID != req.UserID {
		log.Errorf("UpdateCTPCandidateRequest requested user not logged in")
		h.genErrResponse(c, model.ErrCodeRequestUserNotLogin)
		return
	}
	updates := &mysql.PartnerCandidate{
		UserID:      req.UserID,
		CandidateID: req.CandidateID,
		Status:      req.Status,
		Vet:         req.Vet,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Mobile:      req.Mobile,
		Email:       req.Email,
		Expr:        req.Expr,
		LinkedIn:    req.LinkedIn,
		Skill:       req.Skill,
		Tag:         req.Tag,
		Comment:     req.Comment,
		Note:        req.Note,
	}
	resume, errJson := json.Marshal(req.Resume)
	if errJson != nil {
		log.Errorf("json marshal resume err:%v", errJson)
		h.genErrResponse(c, model.ErrCodeInternalServerError)
		return
	}
	workExpr, errJson := json.Marshal(req.WorkExpr)
	if errJson != nil {
		log.Errorf("json marshal workexpr err:%v", errJson)
		h.genErrResponse(c, model.ErrCodeInternalServerError)
		return
	}
	edu, errJson := json.Marshal(req.Education)
	if errJson != nil {
		log.Errorf("json marshal education err:%v", errJson)
		h.genErrResponse(c, model.ErrCodeInternalServerError)
		return
	}
	if len(req.Resume) > 0 {
		updates.Resume = string(resume)
	}
	if len(req.WorkExpr) > 0 {
		updates.WorkExpr = string(workExpr)
	}
	if len(req.Education) > 0 {
		updates.Education = string(edu)
	}
	if err := mysql.UpdateCTPCandidate(req.UserID, req.CandidateID, updates); err != nil {
		log.Errorf("UpdateCTPCandidate err:%v", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "candidate updated"})
}

func (h *PartnerCTPHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}

func (h *PartnerCTPHandler) convertCandidates(candidates []*mysql.PartnerCandidate) ([]*CTPCandidate, error) {
	var res []*CTPCandidate
	for _, c := range candidates {
		can := &CTPCandidate{
			UserID:      c.UserID,
			CandidateID: c.CandidateID,
			Status:      c.Status,
			Vet:         c.Vet,
			FirstName:   c.FirstName,
			LastName:    c.LastName,
			Mobile:      c.Mobile,
			Email:       c.Email,
			Expr:        c.Expr,
			LinkedIn:    c.LinkedIn,
			Skill:       c.Skill,
			Tag:         c.Tag,
			Comment:     c.Comment,
			Note:        c.Note,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		}
		var resume []*CTPResume
		if err := json.Unmarshal([]byte(c.Resume), &resume); err != nil {
			return nil, err
		}
		var workExpr []*CTPWorkExpr
		if err := json.Unmarshal([]byte(c.WorkExpr), &workExpr); err != nil {
			return nil, err
		}
		var edu []*CTPEducation
		if err := json.Unmarshal([]byte(c.Education), &edu); err != nil {
			return nil, err
		}
		can.Resume = resume
		can.WorkExpr = workExpr
		can.Education = edu
		res = append(res, can)
	}
	return res, nil
}
