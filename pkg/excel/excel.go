package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Page struct {
	File  *excelize.File
	Sheet string
}

func CreateFile() *excelize.File {
	return excelize.NewFile()
}

func (e *Page) Title(title string, hCell, vCell string, style int) {
	e.File.MergeCell(e.Sheet, hCell, vCell)
	e.File.SetCellStr(e.Sheet, hCell, title)
	if style > 0 {
		e.File.SetCellStyle(e.Sheet, hCell, vCell, style)
	}
}

func (e *Page) Range(title string, hCell, vCell string, style int) {
	e.File.MergeCell(e.Sheet, hCell, vCell)
	e.File.SetCellStr(e.Sheet, hCell, title)
	if style > 0 {
		e.File.SetCellStyle(e.Sheet, hCell, vCell, style)
	}
}

func (e *Page) CellStyleAlignment(horizontal, vertical string, wrap bool) int {
	style, _ := e.File.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: horizontal,
			Vertical:   vertical,
			WrapText:   wrap,
		},
	})

	return style
}

func (e *Page) CellStyleTitle(horizontal, vertical string, wrap bool, size float64) int {
	style, _ := e.File.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: horizontal,
			Vertical:   vertical,
			WrapText:   wrap,
		},
		Font: &excelize.Font{
			Size: size,
		},
	})

	return style
}

func ToCharStrConst(i int) string {
	return abc[i-1 : i]
}

func CellExcel(nCol, nRow int) string {
	return fmt.Sprintf("%s%v", ToCharStrConst(nCol), nRow)
}

func (e *Page) CellStyleBody(fontSize float64) int {
	style, _ := e.File.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: false,
			Size: fontSize,
		},
	})

	return style
}

func (e *Page) CellStyleBodyColor(fontSize float64, color string) int {
	style, _ := e.File.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{color},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: false,
			Size: fontSize,
		},
	})

	return style
}

func (e *Page) CellStyleBody2(fontSize float64) int {
	style, _ := e.File.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: false,
			Size: fontSize,
		},
	})

	return style
}

func (e *Page) CellStyleHeader(fontSize float64) int {
	style, _ := e.File.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
			Size: fontSize,
		},
	})

	return style
}

func (e *Page) CellStyleHeader2(fontSize float64) int {
	style, _ := e.File.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Font: &excelize.Font{
			Bold: true,
			Size: fontSize,
		},
	})

	return style
}

func (e *Page) SetCellStr(nCol, nRow int, val string) error {
	return e.File.SetCellStr(e.Sheet, CellExcel(nCol, nRow), strings.Trim(val, " "))
}

func (e *Page) SetCellInt(nCol, nRow, val int) error {
	return e.File.SetCellInt(e.Sheet, CellExcel(nCol, nRow), val)
}

func (e *Page) SetCellStyle(hCell, vCell string, val int, size float64) error {
	return e.File.SetCellStyle(e.Sheet, hCell, vCell, e.CellStyleBody(size))
}
