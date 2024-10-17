package dto

type ReadSheetCandidateOffer struct {
	SheetUrl string `json:"sheet_url"`
}

type CreateNewSheet struct {
	SheetName string `json:"sheet_name"`
}

type CreateNewSheetResponse struct {
	SpreadsheetId  string
	SpreadsheetUrl string
}
