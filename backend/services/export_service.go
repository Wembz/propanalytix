package services

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Wembz/propanalytix/backend/config"
	"github.com/Wembz/propanalytix/backend/models"
	"github.com/jung-kurt/gofpdf"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
)

type ExportService struct{}

func NewExportService() *ExportService {
	return &ExportService{}
}

// ExportCalculationToPDF generates a simple PDF of the calculation result
func (s *ExportService) ExportCalculationToPDF(calc *models.Calculation) ([]byte, error) {
	pdf := gofpdf.New("p", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Property Analysis Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10, fmt.Sprintf("Calculation ID: %d\nDate: %s\n\nInput: \n%s\n\nResults:\n%s",
		calc.ID, calc.CreatedAt.Format(time.RFC1123), calc.InputData, calc.Results), "", "L", false)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ExportCalculationToExcel generates an XLSX version of the calculation
func (s *ExportService) ExportCalculationToExcel(calc *models.Calculation) ([]byte, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Calculation")
	if err != nil {
		return nil, err
	}

	sheet.AddRow().AddCell().SetString("property Analysis Report")
	sheet.AddRow()
	sheet.AddRow().AddCell().SetString("Calculation ID")
	sheet.AddRow().AddCell().SetString(fmt.Sprintf("%d", calc.ID))
	sheet.AddRow().AddCell().SetString("Created At")
	sheet.AddRow().AddCell().SetString(calc.CreatedAt.Format(time.RFC1123))
	sheet.AddRow().AddCell().SetString("Input Data")
	sheet.AddRow().AddCell().SetString(calc.InputData)
	sheet.AddRow().AddCell().SetString("Results")
	sheet.AddRow().AddCell().SetString(calc.Results)

	var buf bytes.Buffer
	err = file.Write(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *ExportService) ExportToExcel(userID , calculationID uint) (string, error) {
	var calc models.Calculation

	if err := config.DB.Where("id = ? AND user_id = ? ", userID, calculationID).Find(&calc).Error; err != nil {
		return "", fmt.Errorf("calculation not found")
	}
f := excelize.NewFile()
	sheet := "Summary"
	f.NewSheet(sheet)

	// Populate sample data — adapt this to your actual calculation fields
	f.SetCellValue(sheet, "A1", "Calculation ID")
	f.SetCellValue(sheet, "B1", calc.ID)
	f.SetCellValue(sheet, "A2", "User ID")
	f.SetCellValue(sheet, "B2", calc.UserID)
	f.SetCellValue(sheet, "A3", "Result")
	f.SetCellValue(sheet, "B3", calc.Results) // adjust this to match your schema

	// Create exports directory if needed
	outputDir := "./exports"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", err
	}

	// Generate filename
	filename := fmt.Sprintf("calc_%d_%d.xlsx", userID, time.Now().Unix())
	filepath := filepath.Join(outputDir, filename)

	// Save the file
	if err := f.SaveAs(filepath); err != nil {
		return "", err
	}

	return filepath, nil

}

func (s *ExportService) ExportToPDF(userID, calculationID uint) (string, error) {
	
}