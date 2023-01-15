package repository

import (
	"fmt"

	"google.golang.org/api/sheets/v4"
)

type GSheetRepository struct {
	gsheetCli *sheets.Service
	GSheetID  string
}

func NewGSheetRepository(gsheetCli *sheets.Service, gsheetID string) *GSheetRepository {
	return &GSheetRepository{gsheetCli, gsheetID}
}

func (g *GSheetRepository) GetSheetLink() string {
	return fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit", g.GSheetID)
}

func (g *GSheetRepository) CheckSheetIfExist(sheetName string) bool {
	resp, err := g.gsheetCli.Spreadsheets.Get(g.GSheetID).Do()
	if err != nil {
		return false
	}

	for _, sheet := range resp.Sheets {
		if sheetName == sheet.Properties.Title {
			return true
		}
	}

	return false
}

func (g *GSheetRepository) CreateSheet(sheetName string) error {
	_, err := g.gsheetCli.Spreadsheets.BatchUpdate(g.GSheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AddSheet: &sheets.AddSheetRequest{
					Properties: &sheets.SheetProperties{
						Title: sheetName,
					},
				},
			},
		},
	}).Do()
	if err != nil {
		return err
	}

	return nil
}

func (g *GSheetRepository) GetSheetData(sheetId string) ([][]string, error) {
	resp, err := g.gsheetCli.Spreadsheets.Values.Get(g.GSheetID, sheetId+"!A1:Z").Do()
	if err != nil {
		return nil, err
	}

	sheetData := [][]string{}
	for _, row := range resp.Values {
		rowData := []string{}
		for _, column := range row {
			rowData = append(rowData, fmt.Sprint(column))
		}
		sheetData = append(sheetData, rowData)
	}

	return sheetData, nil
}

func (g *GSheetRepository) UpdateSheetData(sheetId string, data [][]string) error {
	sheetData := [][]interface{}{}
	for _, row := range data {
		rowData := []interface{}{}
		for _, column := range row {
			rowData = append(rowData, column)
		}
		sheetData = append(sheetData, rowData)
	}

	rb := &sheets.ClearValuesRequest{}

	_, err := g.gsheetCli.Spreadsheets.Values.Clear(g.GSheetID, sheetId+"!A1", rb).Do()
	if err != nil {
		return err
	}

	_, err = g.gsheetCli.Spreadsheets.Values.Append(g.GSheetID, sheetId+"!A1", &sheets.ValueRange{
		Values: sheetData,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}
