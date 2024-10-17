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

func NewGSheetService(service *sheets.Service, driveService *drive.Service) *GSheetService {
	return &GSheetService{
		SheetService: service,
		DriveService: driveService,
	}
}

func (s *GSheetService) ReadCandidateOffer(spreadsheetId string) ([]dto.SheetCandidateOffer, error) {
	spreadsheetID, err := google_internal.ExtractSheetIdFromUrl(spreadsheetId)
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
	return list, nil
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
