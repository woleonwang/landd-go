package model

type ContactType int

const (
	ContactMobile   ContactType = 1
	ContactWhatsapp ContactType = 2
	ContactTelegram ContactType = 3
	ContactWechat   ContactType = 4
)

type EndorsementStatus int

const (
	EndorsementStatusInvited   EndorsementStatus = 1
	EndorsementStatusSubmitted EndorsementStatus = 2
	EndorsementStatusApproved  EndorsementStatus = 3
	EndorsementStatusDeclined  EndorsementStatus = 4
)

type EndorserIdentity int

const (
	EndorserCandidate EndorserIdentity = 1
	EndorserClient    EndorserIdentity = 2
)

type UserRole int

const (
	Recruiter UserRole = 1
	Partner   UserRole = 2
	Admin     UserRole = 3
)

type CTPStatus int

const (
	CTPStatusPending  CTPStatus = 1
	CTPStatusApproved CTPStatus = 2
	CTPStatusDeclined CTPStatus = 3
)

type Vet string

const (
	VetExColleague Vet = "colleague"
	VetFriend      Vet = "friend"
	VetSpoken      Vet = "spoken"
	VetRecommended Vet = "recommended"
	VetNone        Vet = "none"
	VetOther       Vet = "other"
)
