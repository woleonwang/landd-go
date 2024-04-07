package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/mysql"
	"net/http"
	"strconv"
)

type JobHandler struct {
}

func NewJobHandler() *JobHandler {
	return &JobHandler{}
}

func (h *JobHandler) Get(c *gin.Context) {
	jobIDStr := c.Query("job_id")
	if jobIDStr != "" {
		jobID, err1 := strconv.ParseInt(jobIDStr, 10, 64)
		if err1 != nil {
			log.Errorf("convert request jobID error: %v ", err1)
			h.genErrResponse(c, model.ErrCodeInvalidRequest)
			return
		}
		job, err := mysql.GetJobByID(jobID)
		if err != nil {
			log.Errorf("GetJobByID err:%v", err)
			h.genErrResponse(c, model.ErrCodeMysqlError)
			return
		}
		resp := GetJobsResponse{
			Jobs: []*mysql.Job{job},
		}
		c.JSON(http.StatusOK, gin.H{"message": resp})
		return
	}
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	if size > 100 || size <= 0 || page <= 0 {
		log.Errorf("page size invalid")
		h.genErrResponse(c, model.ErrCodePageSizeInvalid)
		return
	}
	jobs, err2 := mysql.GetJobOpenings((page-1)*size, size)
	if err2 != nil {
		log.Errorf("GetJobOpenings err:%v", err2)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	resp := GetJobsResponse{
		Jobs: jobs,
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

func (h *JobHandler) Create(c *gin.Context) {
	var req CreateJobRequest
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("CreateJobRequest bind request error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	log.Infof("CreateJobRequest request: %v ", req)
	job := &mysql.Job{
		Title:         req.Title,
		Company:       req.Company,
		Jd:            req.Jd,
		AboutCompany:  req.AboutCompany,
		Comment:       req.Comment,
		ReferralFee:   req.ReferralFee,
		LowerBoundSal: req.LowerBoundSal,
		UpperBoundSal: req.UpperBoundSal,
		PosterID:      req.PosterID,
	}
	jobID, err := mysql.CreateJobOpening(job)
	if err != nil {
		log.Errorf("CreateJobOpening err:%v", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	resp := CreateJobResponse{
		JobID: strconv.FormatInt(jobID, 10),
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})
}

func (h *JobHandler) Update(c *gin.Context) {
	var req UpdateJobRequest
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("UpdateJobRequest bind request error: %v ", err)
		h.genErrResponse(c, model.ErrCodeInvalidRequest)
		return
	}
	log.Infof("UpdateJobRequest request: %v ", req)
	updates := &mysql.Job{
		JobID:         req.JobID,
		Title:         req.Title,
		Company:       req.Company,
		Jd:            req.Jd,
		AboutCompany:  req.AboutCompany,
		Comment:       req.Comment,
		ReferralFee:   req.ReferralFee,
		LowerBoundSal: req.LowerBoundSal,
		UpperBoundSal: req.UpperBoundSal,
	}
	if err := mysql.UpdateJobOpening(req.JobID, updates); err != nil {
		log.Errorf("UpdateJobOpening err:%v", err)
		h.genErrResponse(c, model.ErrCodeMysqlError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "job opening updated"})
}

func (h *JobHandler) genErrResponse(c *gin.Context, err model.CustomError) {
	resp := BaseResponse{ErrCode: err.Code, Message: err.Message}
	c.JSON(err.HttpCode, gin.H{"message": resp})
}
