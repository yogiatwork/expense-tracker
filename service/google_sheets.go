package service

import (
	"log/slog"
	"time"

	"github.com/yogiatwork/expense-tracker/config"
	"github.com/yogiatwork/expense-tracker/model"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// sheetService is a struct that holds google sheet service
type sheetService struct {
	GSheetSrv *sheets.Service
}

// SheetService is a global variable that holds the sheet service
var SheetService *sheetService

// InitSheetService initializes the sheet service
func InitSheetService() {
	// create a new sheet service
	slog.Info("creating sheet service", slog.String("credentials_file", config.AppConfig.GoogleSheets.CredentialsFile))
	srv, err := sheets.NewService(nil, option.WithCredentialsFile(config.AppConfig.GoogleSheets.CredentialsFile))
	if err != nil {
		slog.Error("failed to create sheet service", slog.String("error", err.Error()))
	}
	slog.Info("sheet service created successfully", slog.Any("service", srv))

	SheetService = &sheetService{
		GSheetSrv: srv,
	}
}

// GetSheetMetadata retrieves the metadata of the sheet
func (s *sheetService) GetSheetMetadata() error {
	sheet, err := s.GSheetSrv.Spreadsheets.Get(config.AppConfig.GoogleSheets.SpreadsheetId).Do()
	if err != nil {
		slog.Error("failed to get sheet metadata", slog.String("error", err.Error()))
		return err
	}
	slog.Info("sheet metadata retrieved successfully", slog.Any("metadata", sheet))
	return nil
}

func (s *sheetService) AppendSheet(txData model.Transaction) error {
	if s.GSheetSrv == nil {
		slog.Error("sheet service is not initialized")
		return nil
	}

	appendValueRange := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values: [][]interface{}{
			txData.ToRow(), // Call our mapping method here
		},
	}

	// Define the range and values to append
	// TODO: Make this dynamic based on the current date or other logic
	// Use date time to determine the correct sheet name and if missing create a new sheet for the month
	sheetName := getSheetName()
	slog.Info("sheet name determined", slog.String("sheet_name", sheetName))
	if sheetsExist, err := s.checkIfSheetExists(sheetName); err != nil {
		slog.Error("failed to check if sheet exists", slog.String("error", err.Error()))
		return err
	} else if !sheetsExist {
		if err := s.createNewSheet(sheetName); err != nil {
			slog.Error("failed to create new sheet", slog.String("error", err.Error()))
			return err
		}
	}

	rangeToAppend := sheetName + "!A1" // Change this to your desired range

	// Append the values to the sheet
	_, err := s.GSheetSrv.Spreadsheets.Values.Append(config.AppConfig.GoogleSheets.SpreadsheetId, rangeToAppend, appendValueRange).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		slog.Error("failed to append values to sheet", slog.String("error", err.Error()))
		return err
	}

	slog.Info("values appended successfully to sheet")
	return nil
}

// getSheetName returns the sheet name based on the currentDate
func getSheetName() string {
	// Get the current date
	currentDate := time.Now()

	// Format the date to "Jan2006" format (e.g., "Sep2023")
	sheetName := currentDate.Format("Jan2006")

	return sheetName
}

// checkIfSheetExists checks if a sheet with the given name exists in the spreadsheet
func (s *sheetService) checkIfSheetExists(sheetName string) (bool, error) {
	slog.Info("spreadsheet id", slog.String("spreadsheet_id", config.AppConfig.GoogleSheets.SpreadsheetId))
	spreadsheet, err := s.GSheetSrv.Spreadsheets.Get(config.AppConfig.GoogleSheets.SpreadsheetId).Do()
	if err != nil {
		slog.Error("failed to get spreadsheet", slog.String("error", err.Error()))
		return false, err
	}

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == sheetName {
			return true, nil
		}
	}

	return false, nil
}

// createNewSheet creates a new sheet with the given name in the spreadsheet
func (s *sheetService) createNewSheet(sheetName string) error {
	request := &sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: sheetName,
			},
		},
	}

	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{request},
	}

	_, err := s.GSheetSrv.Spreadsheets.BatchUpdate(config.AppConfig.GoogleSheets.SpreadsheetId, batchUpdateRequest).Do()
	if err != nil {
		slog.Error("failed to create new sheet", slog.String("error", err.Error()))
		return err
	}

	slog.Info("new sheet created successfully", slog.String("sheet_name", sheetName))
	return nil
}
