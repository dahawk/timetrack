//contains all functions and structs to render pdf
package main

import (
	"fmt"
	"io"

	"github.com/jung-kurt/gofpdf"
)

func generatePDF(w io.Writer, logs []DisplayLog, data anonStruct) error {
	Info.Println("generatePDF()")
	//print header
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetAutoPageBreak(true, 20)
	pdf.SetMargins(20, 20, 20)

	err := printHeader(pdf, data, logs)
	if err != nil {
		return err
	}

	//print table
	err = printTable(pdf, logs)
	if err != nil {
		return err
	}

	err = pdf.Output(w)
	if err != nil {
		return err
	}

	return nil
}

func printHeader(pdf *gofpdf.Fpdf, data anonStruct, logs []DisplayLog) error {
	err := printUserData(pdf, data)
	if err != nil {
		return err
	}

	err = printUserStats(pdf, logs, data)
	if err != nil {
		return err
	}
	pdf.Ln(-1)

	return pdf.Error()
}

func printUserData(pdf *gofpdf.Fpdf, data anonStruct) error {
	pdf.SetFont("Arial", "", 24)
	pdf.CellFormat(0, 20, "TimeTrack Report", "", 1, "", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 7, fmt.Sprintf("Name: %s", data.User.Name), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("Period: %s - %s", data.From, data.To), "", 1, "", false, 0, "")

	return pdf.Error()
}

func printUserStats(pdf *gofpdf.Fpdf, data []DisplayLog, anon anonStruct) error {
	stats := calculateStats(data, anon)
	pdf.Ln(-1)
	pdf.CellFormat(0, 7, fmt.Sprintf("Expected Work time: %s", stats.ExtectedWorkTime), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("Actual Work time: %s", stats.ActualWorkTime), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("Difference: %s", stats.Delta), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("Holidays: %d days", stats.Holidays), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("Sick leave: %d days", stats.Sickdays), "", 1, "", false, 0, "")

	return pdf.Error()
}

func printTable(pdf *gofpdf.Fpdf, data []DisplayLog) error {
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(30, 6, "", "B", 0, "", false, 0, "")
	pdf.CellFormat(20, 6, "From", "B", 0, "", false, 0, "")
	pdf.CellFormat(30, 6, "", "B", 0, "", false, 0, "")
	pdf.CellFormat(20, 6, "To", "B", 0, "", false, 0, "")
	pdf.CellFormat(30, 6, "Duration", "B", 0, "", false, 0, "")
	pdf.CellFormat(20, 6, "Type", "B", 0, "", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 12)
	for _, l := range data {
		if l.ToDate == nil {
			continue
		}
		pdf.CellFormat(30, 6, l.DateFrom, "", 0, "", false, 0, "")
		pdf.CellFormat(20, 6, l.TimeFrom, "", 0, "", false, 0, "")
		if l.Type == workTimeConst {
			pdf.CellFormat(30, 6, "", "", 0, "", false, 0, "")
		} else {
			pdf.CellFormat(30, 6, l.DateTo, "", 0, "", false, 0, "")
		}

		pdf.CellFormat(20, 6, l.TimeTo, "", 0, "", false, 0, "")
		pdf.CellFormat(30, 6, l.Duration, "", 0, "", false, 0, "")
		pdf.CellFormat(20, 6, l.Type, "", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	return pdf.Error()
}
