package dto

type SheetCandidateOffer struct {
	ID               string `json:"id" sheets:"No"`
	FullName         string `json:"full_name" sheets:"Full Name"`
	Position         string `json:"position" sheets:"Position"`
	AcceptDate       string `json:"accept_date" sheets:"Accept-offer date"`
	OnboardDate      string `json:"onboard_date" sheets:"Expected-onboard date"`
	Division         string `json:"division" sheets:"Division"`
	HR               string `json:"hr" sheets:"HR"`
	Point            string `json:"point" sheets:"Point"`
	Level            string `json:"level" sheets:"Level"`
	CloseRequestDate string `json:"close_request_date" sheets:"Close-request date"`
	Source           string `json:"source" sheets:"Source"`
}
