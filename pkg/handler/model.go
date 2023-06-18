package handler

import (
	"landd.co/landd/pkg/model"
	"landd.co/landd/pkg/mysql"
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
