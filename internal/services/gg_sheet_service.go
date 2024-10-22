package services

import (
	"fmt"
	"log"
	"time"

	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/google_internal"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/util"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/sheets/v4"
)

const (
	SharedDriveFolderId = "1izJFBRgsYtpzOL3ksrhFDs3uh1E1kjAx"
)

type GSheetService struct {
	SheetService *sheets.Service
	DriveService *drive.Service
}

type IGSheetService interface {
	ReadCandidateOffer(spreadsheetUrl string) ([]dto.SheetCandidateOffer, error)
	CreateNewSheet(sheetName string) (*dto.CreateNewSheetResponse, error)
	CreateNewSheetInSharedDrive(sheetName string, sharedDriveFolderId string) (*dto.CreateNewSheetResponse, error)
	InsertDataToSheet(spreadsheetID string, sheetName string, data []dto.SheetCandidateOffer) error
	HandleFileCandidateOffer(sheetUrl string) (*dto.CreateNewSheetResponse, error)
}

func NewGSheetService(service *sheets.Service, driveService *drive.Service) *GSheetService {
	return &GSheetService{
		SheetService: service,
		DriveService: driveService,
	}
}

func (s *GSheetService) ReadCandidateOffer(spreadsheetUrl string) ([]dto.SheetCandidateOffer, error) {
	spreadsheetID, err := google_internal.ExtractSheetIdFromUrl(spreadsheetUrl)
	if err != nil {
		return nil, err
	}
	list, err := util.ParseSheetIntoStructSlice[dto.SheetCandidateOffer](util.Options{
		Service:       s.SheetService,
		SpreadsheetID: spreadsheetID,
		SheetName:     "A1:K20",
	}.Build())
	if err != nil {
		return nil, err
	}

	filteredList := make([]dto.SheetCandidateOffer, 0, len(list))
	for _, offer := range list {
		if offer.FullName != "" {
			filteredList = append(filteredList, offer)
		}
	}

	return filteredList, nil

	// var list []dto.SheetCandidateOffer
	// can := util.ParseCandidateOffer[dto.SheetCandidateOffer](resp)
	// list = append(list, can)
	// return list, nil
}

func (s *GSheetService) CreateNewSheet(sheetName string) (*dto.CreateNewSheetResponse, error) {
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: sheetName,
		},
	}
	resp, err := s.SheetService.Spreadsheets.Create(spreadsheet).Do()
	if err != nil {
		log.Fatalf("Unable to create spreadsheet: %v", err)
	}
	fmt.Println("resp", resp)
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", resp.SpreadsheetId)
	return &dto.CreateNewSheetResponse{
		SpreadsheetId:  resp.SpreadsheetId,
		SpreadsheetUrl: url,
	}, nil
}

func (s *GSheetService) CreateNewSheetInSharedDrive(sheetName string, sharedDriveFolderId string) (*dto.CreateNewSheetResponse, error) {
	// Create a new spreadsheet
	spreadsheet, err := s.CreateNewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("unable to create new spreadsheet: %v", err)
	}

	// Move the spreadsheet to the shared drive folder
	file, err := s.DriveService.Files.Get(spreadsheet.SpreadsheetId).Fields("parents").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get file: %v", err)
	}

	_, err = s.DriveService.Files.Update(spreadsheet.SpreadsheetId, nil).
		AddParents(sharedDriveFolderId).
		RemoveParents(file.Parents[0].Id).
		Fields("id, parents").
		SupportsAllDrives(true).
		Do()
	if err != nil {
		return nil, fmt.Errorf("unable to move file to shared drive: %v", err)
	}

	return spreadsheet, nil
}

func (s *GSheetService) InsertDataToSheet(spreadsheetID string, sheetName string, data []dto.SheetCandidateOffer) error {
	writeRange := "A1"

	employees := make([]dto.NewEmployeesSkills, 0, len(data))
	for _, candidate := range data {
		employees = append(employees, *util.FromCandidateToEmployee(&candidate))
	}

	var vr sheets.ValueRange

	rows := util.ParseStructToSheetTable(employees)
	vr.Values = append(vr.Values, rows...)

	_, err := s.SheetService.Spreadsheets.Values.Update(spreadsheetID, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}

	return nil
}

func (s *GSheetService) HandleFileCandidateOffer(sheetUrl string) (*dto.CreateNewSheetResponse, error) {
	data, err := s.ReadCandidateOffer(sheetUrl)
	fmt.Println("ReadCandidateOffer", data, err, sheetUrl)
	if err != nil {
		return nil, err
	}

	newSheetName := fmt.Sprintf("New Employee Skill - %s", time.Now().Format("2006-01-02"))
	newEmployeeSkillFile, err := s.CreateNewSheetInSharedDrive(newSheetName, SharedDriveFolderId)
	fmt.Println("newEmployeeSkillFile", newEmployeeSkillFile)
	if err != nil {
		return nil, err
	}

	err = s.InsertDataToSheet(newEmployeeSkillFile.SpreadsheetId, "A1:K20", data)
	fmt.Println("InsertDataToSheet", err)
	if err != nil {
		return nil, err
	}

	return newEmployeeSkillFile, nil
}

// func (s *GSheetService) ReadSheetData(sheetUrl string) ([][]interface{}, []string, error) {
// 	sheetID, err := google_internal.ExtractSheetIdFromUrl(sheetUrl)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	resp, err := s.SheetService.Spreadsheets.Values.Get(sheetID, "A1:ZZ").Do()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	if len(resp.Values) == 0 {
// 		return nil, nil, fmt.Errorf("no data found")
// 	}

// 	headers := make([]string, len(resp.Values[0]))
// 	for i, v := range resp.Values[0] {
// 		headers[i] = fmt.Sprint(v)
// 	}

// 	return resp.Values[1:], headers, nil
// }

// func (s *GSheetService) applyMappings(candidateData [][]interface{}, mappings map[string]string) [][]interface{} {
// 	var skillData [][]interface{}

// 	// Create header row
// 	headerRow := make([]interface{}, len(mappings))
// 	for skillCol, _ := range mappings {
// 		headerRow = append(headerRow, skillCol)
// 	}
// 	skillData = append(skillData, headerRow)

// 	// Map data
// 	for _, row := range candidateData {
// 		newRow := make([]interface{}, len(mappings))
// 		for skillCol, candidateCol := range mappings {
// 			index := s.findIndex(candidateData[0], candidateCol)
// 			if index != -1 && index < len(row) {
// 				newRow = append(newRow, row[index])
// 			} else {
// 				newRow = append(newRow, "")
// 			}
// 		}
// 		skillData = append(skillData, newRow)
// 	}

// 	return skillData
// }

// func (s *GSheetService) findIndex(slice []interface{}, value string) int {
// 	for i, v := range slice {
// 		if fmt.Sprint(v) == value {
// 			return i
// 		}
// 	}
// 	return -1
// }

// func (s *GSheetService) AIAssistedFillSkillFile(candidateSheetUrl, skillSheetUrl string) error {
// 	// Read candidate data
// 	candidateData, candidateHeaders, err := s.ReadSheetData(candidateSheetUrl)
// 	if err != nil {
// 		return fmt.Errorf("error reading candidate data: %v", err)
// 	}

// 	// Read skill sheet headers
// 	_, skillHeaders, err := s.ReadSheetData(skillSheetUrl)
// 	if err != nil {
// 		return fmt.Errorf("error reading skill sheet headers: %v", err)
// 	}

// 	// Use AI to suggest mappings
// 	mappings, err := ai_service.SuggestColumnMappings(candidateHeaders, skillHeaders)
// 	if err != nil {
// 		return fmt.Errorf("error suggesting column mappings: %v", err)
// 	}

// 	// Apply mappings to create skill data
// 	skillData := s.applyMappings(candidateData, mappings)

// 	// Insert data into skill sheet
// 	skillSheetID, err := google_internal.ExtractSheetIdFromUrl(skillSheetUrl)
// 	if err != nil {
// 		return fmt.Errorf("error extracting skill sheet ID: %v", err)
// 	}

// 	err = s.InsertDataToSheet(skillSheetID, "A1", skillData)
// 	if err != nil {
// 		return fmt.Errorf("error inserting data into skill sheet: %v", err)
// 	}

// 	return nil
// }
