package handler

import (
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/mysql"
	"time"
)

type BaseResponse struct {
	ErrCode int    `json:"err_code"` // error code, 0 if none
	Message string `json:"message"`
}

type GetProfileInfoResponse struct {
	UserID       int64                         `json:"user_id"`
	Profile      *mysql.RecruiterProfile       `json:"profile"`
	Placements   []*mysql.RecruiterPlacement   `json:"placements"`
	Jobs         []*mysql.RecruiterJob         `json:"jobs"`
	Candidates   []*mysql.RecruiterCandidate   `json:"candidates"`
	Publications []*mysql.RecruiterPublication `json:"publications"`
}

type UpdateProfileInfoRequest struct {
	UserID             int64                 `json:"user_id"`
	ProfileChanges     *ProfileChanges       `json:"profile"`
	PlacementChanges   []*PlacementChanges   `json:"placements"`
	JobChanges         []*JobChanges         `json:"jobs"`
	CandidateChanges   []*CandidateChanges   `json:"candidates"`
	PublicationChanges []*PublicationChanges `json:"publications"`
}

type ProfileChanges struct {
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

type PlacementChanges struct {
	Date     int64  `json:"date"`
	Position string `json:"position"`
	Company  string `json:"company"`
	Verified bool   `json:"verified"`
}

type JobChanges struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Description string `json:"description"`
}

type CandidateChanges struct {
	Title      string  `json:"title"`
	Percentage float32 `json:"percentage"`
}

type PublicationChanges struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type InviteEndorseResponse struct {
	UserID   int64 `json:"user_id"`
	InviteID int64 `json:"invite_id"`
}

type CreateEndorsementRequest struct {
	UserID int64 `json:"user_id"`
}

type UpdateEndorsementRequest struct {
	UserID   int64                   `json:"user_id"`
	InviteID int64                   `json:"invite_id"`
	Endorser string                  `json:"endorser"`
	Title    string                  `json:"title"`
	Company  string                  `json:"company"`
	Identity model.EndorserIdentity  `json:"identity"`
	Content  string                  `json:"content"`
	Status   model.EndorsementStatus `json:"status"`
}

type GetEndorsementResponse struct {
	UserID       int64                `json:"user_id"`
	Endorsements []*mysql.Endorsement `json:"endorsements"`
}

type UpdateEndorseDraftRequest struct {
	UserID   int64                  `json:"user_id"`
	Endorser string                 `json:"endorser"`
	Title    string                 `json:"title"`
	Company  string                 `json:"company"`
	Identity model.EndorserIdentity `json:"identity"`
	Content  string                 `json:"content"`
}

type UpdatePartnerHomepageRequest struct {
	UserID      int64                   `json:"user_id"`
	Audience    string                  `json:"audience"`
	DisplayName string                  `json:"display_name"`
	Summary     string                  `json:"summary"`
	CTPSummary  string                  `json:"ctp_summary"`
	Reasons     []PartnerHomepageReason `json:"reasons"`
	Companies   string                  `json:"companies"`
	DataPolicy  string                  `json:"data_policy"`
	Applicants  string                  `json:"applicants"`
	HowTo       string                  `json:"howto"`
}

type PartnerHomepageReason struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type GetPartnerProfileResponse struct {
	UserID   int64                 `json:"user_id"`
	Profile  *mysql.PartnerProfile `json:"profile"`
	Email    string                `json:"email"`
	Contacts []*mysql.UserContact  `json:"contacts"`
}

type UpdatePartnerProfileRequest struct {
	UserID    int64                        `json:"user_id"`
	Name      string                       `json:"name"`
	Photo     string                       `json:"photo"`
	Company   string                       `json:"company"`
	Contacts  map[model.ContactType]string `json:"contacts"`
	Twitter   string                       `json:"twitter"`
	LinkedIn  string                       `json:"linkedin"`
	Website   string                       `json:"website"`
	Blog      string                       `json:"blog"`
	Facebook  string                       `json:"facebook"`
	Instagram string                       `json:"instagram"`
	Tiktok    string                       `json:"tiktok"`
	Youtube   string                       `json:"youtube"`
	Other     string                       `json:"other"`
}

type GetCTPCandidateResponse struct {
	UserID     int64           `json:"user_id"`
	Candidates []*CTPCandidate `json:"candidates"`
}

type CTPCandidate struct {
	UserID      int64           `json:"user_id"`
	CandidateID int64           `json:"candidate_id"`
	Status      model.CTPStatus `json:"status"`
	Vet         model.Vet       `json:"vet"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	Mobile      string          `json:"mobile"`
	Email       string          `json:"email"`
	Expr        int             `json:"expr"`
	LinkedIn    string          `json:"linkedin"`
	Resume      []*CTPResume    `json:"resume"`
	WorkExpr    []*CTPWorkExpr  `json:"work_expr"`
	Education   []*CTPEducation `json:"education"`
	Skill       string          `json:"skill"`
	Tag         string          `json:"tag"`
	Comment     string          `json:"comment"`
	Note        string          `json:"note"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CTPResume struct {
	FileName string `json:"file_name"`
	FileID   string `json:"file_id"`
}

type CTPWorkExpr struct {
	Company  string `json:"company"`
	Title    string `json:"title"`
	Function string `json:"function"`
}

type CTPEducation struct {
	School string `json:"school"`
	Degree string `json:"degree"`
	Major  string `json:"major"`
}

type CreateCTPCandidateRequest struct {
	UserID    int64           `json:"user_id"`
	Vet       model.Vet       `json:"vet"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Mobile    string          `json:"mobile"`
	Email     string          `json:"email"`
	Expr      int             `json:"expr"`
	LinkedIn  string          `json:"linkedin"`
	Resume    []*CTPResume    `json:"resume"`
	WorkExpr  []*CTPWorkExpr  `json:"work_expr"`
	Education []*CTPEducation `json:"education"`
	Skill     string          `json:"skill"`
	Tag       string          `json:"tag"`
	Comment   string          `json:"comment"`
	Note      string          `json:"note"`
}

type CreateCTPCandidateResponse struct {
	UserID      int64 `json:"user_id"`
	CandidateID int64 `json:"candidate_id"`
}

type UpdateCTPCandidateRequest struct {
	UserID      int64           `json:"user_id"`
	CandidateID int64           `json:"candidate_id"`
	Status      model.CTPStatus `json:"status"`
	Vet         model.Vet       `json:"vet"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	Mobile      string          `json:"mobile"`
	Email       string          `json:"email"`
	Expr        int             `json:"expr"`
	LinkedIn    string          `json:"linkedin"`
	Resume      []*CTPResume    `json:"resume"`
	WorkExpr    []*CTPWorkExpr  `json:"work_expr"`
	Education   []*CTPEducation `json:"education"`
	Skill       string          `json:"skill"`
	Tag         string          `json:"tag"`
	Comment     string          `json:"comment"`
	Note        string          `json:"note"`
}

type GetJobsResponse struct {
	Jobs []*mysql.Job `json:"jobs"`
}

type CreateJobRequest struct {
	Title         string `json:"title"`
	Company       string `json:"company"`
	Jd            string `json:"jd"`
	AboutCompany  string `json:"about_company"`
	Comment       string `json:"comment"`
	ReferralFee   int    `json:"referral_fee"`
	LowerBoundSal int    `json:"lower_bound_sal"`
	UpperBoundSal int    `json:"upper_bound_sal"`
	PosterID      int64  `json:"poster_id"`
}

type CreateJobResponse struct {
	JobID int64 `json:"job_id"`
}

type UpdateJobRequest struct {
	JobID         int64  `json:"job_id"`
	Title         string `json:"title"`
	Company       string `json:"company"`
	Jd            string `json:"jd"`
	AboutCompany  string `json:"about_company"`
	Comment       string `json:"comment"`
	ReferralFee   int    `json:"referral_fee"`
	LowerBoundSal int    `json:"lower_bound_sal"`
	UpperBoundSal int    `json:"upper_bound_sal"`
}
