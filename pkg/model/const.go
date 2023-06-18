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
