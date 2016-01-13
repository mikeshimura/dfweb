package excel

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"strings"
	"time"
)

//border thin medium double など
type ExcelService struct {
	File  *xlsx.File
	Sheet *xlsx.Sheet
	Row   int
}

func (e *ExcelService) GetRow(no int) *xlsx.Row {
	addRowNo := no - e.Sheet.MaxRow + 1
	fmt.Printf("addRowNo %v \n", addRowNo)
	for i := 0; i < addRowNo; i++ {
		e.Sheet.AddRow()
		fmt.Printf("i %v \n", i)
	}

	return e.Sheet.Rows[no]
}
func (e *ExcelService) GetCell(rowno int, colno int) *xlsx.Cell {
	row := e.GetRow(rowno)
	addColNo := colno - len(row.Cells) + 1
	fmt.Printf("addColNo %v \n", addColNo)
	for i := 0; i < addColNo; i++ {
		row.AddCell()
		fmt.Printf("i %v \n", i)
	}

	return row.Cells[colno]
}
func (e *ExcelService) AddSheetPanic(name string) {
	sheet, err := e.File.AddSheet(name)
	if err != nil {
		panic(err)
	}
	e.Sheet = sheet
}

func (e *ExcelService) SetSheetPanic(name string) {
	sheet := e.File.Sheet[name]
	if sheet != nil {
		e.Sheet = sheet
	} else {
		panic("Sheet :" + name + " が見つかりません")
	}
}
func (e *ExcelService) SetString(rowno int, colno int, s string) {
	cell := e.GetCell(rowno, colno)
	cell.SetString(s)
}
func (e *ExcelService) SetInt(rowno int, colno int, i int) {
	cell := e.GetCell(rowno, colno)
	fmt := cell.GetNumberFormat()
	cell.SetInt(i)
	if fmt != "" {
		cell.NumFmt = fmt
	}
}
func (e *ExcelService) SetIntFormat(rowno int, colno int, i int, fmt string) {
	cell := e.GetCell(rowno, colno)
	cell.SetInt(i)
	cell.NumFmt = fmt
}
func (e *ExcelService) SetFloat(rowno int, colno int, f float64) {
	cell := e.GetCell(rowno, colno)
	fmt := cell.GetNumberFormat()
	cell.SetFloat(f)
	if fmt != "" {
		cell.NumFmt = fmt
	}
}
func (e *ExcelService) SetFloatFormat(rowno int, colno int, f float64, fmt string) {
	cell := e.GetCell(rowno, colno)
	cell.SetFloatWithFormat(f, fmt)
}
func (e *ExcelService) SetDate(rowno int, colno int, t time.Time) {
	cell := e.GetCell(rowno, colno)
	fmt := cell.GetNumberFormat()
	cell.SetDate(t)
	if fmt != "" {
		cell.NumFmt = fmt
	}
}
func (e *ExcelService) SetDateFormat(rowno int, colno int, t time.Time, fmt string) {
	cell := e.GetCell(rowno, colno)
	cell.SetDate(t)
	if fmt != "" {
		cell.NumFmt = fmt
	}
}
func (e *ExcelService) SetDateTime(rowno int, colno int, t time.Time) {
	cell := e.GetCell(rowno, colno)
	fmt := cell.GetNumberFormat()
	cell.SetDateTime(t)
	if fmt != "" {
		cell.NumFmt = fmt
	}
}
func (e *ExcelService) SetDateTimeFormat(rowno int, colno int, t time.Time, fmt string) {
	cell := e.GetCell(rowno, colno)
	cell.SetDateTime(t)
	if fmt != "" {
		cell.NumFmt = fmt
	}
}

func (e *ExcelService) SetStyle(rowno int, colno int, style *xlsx.Style) {
	cell := e.GetCell(rowno, colno)
	cell.SetStyle(style)
}
func (e *ExcelService) GetStyle(rowno int, colno int) *xlsx.Style {
	cell := e.GetCell(rowno, colno)
	return cell.GetStyle()
}
func (e *ExcelService) CreateStyle(font string, size int, border string,
	borderName string) *xlsx.Style {
	style := new(xlsx.Style)
	style.Font.Name = font
	style.Font.Size = size
	if len(border) > 0 {
		style.ApplyBorder = true
	}
	border = strings.ToUpper(border)
	if strings.Contains(border, "T") {
		style.Border.Top = borderName
	}
	if strings.Contains(border, "B") {
		style.Border.Bottom = borderName
	}
	if strings.Contains(border, "R") {
		style.Border.Right = borderName
	}
	if strings.Contains(border, "L") {
		style.Border.Left = borderName
	}
	return style
}
