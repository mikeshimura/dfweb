package pdf

import (
	"fmt"
	"github.com/mikeshimura/dbflute/df"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func TestReport(t *testing.T) {
	r := CreateReport()
	r.Pdf.FontDir="D:\\AW\\picorre\\wkpicorre\\to-bi-develop\\font"
	r.PageTotal=true
	h := new(TestPageHeader)
	var hdr Band = h
	r.RegisterBand(&hdr, PageHeader)
	d := new(TestDetail)
	var det Band = d
	r.RegisterBand(&det, Detail)
	gh := new(TestBeforeGroup2)
	var ghb Band = gh
	r.RegisterGroupBand(&ghb, GroupHeader,2)
	ga := new(TestAfterGroup1)
	var gab Band = ga
	r.RegisterGroupBand(&gab, GroupSummary,1)
	ga2 := new(TestAfterGroup2)
	var gab2 Band = ga2
	r.RegisterGroupBand(&gab2, GroupSummary,2)
	r.Data = getData()
	fmt.Printf("Data %v \n", r.Data)
	r.SetPage("A4", "L")
	r.MaxY = 190
	r.Execute("pdfreport.pdf")
	r.SaveDoc("savedoc.txt")
}

type TestPageHeader struct {
	Report *Report
}

func (h *TestPageHeader) GetHeight() float64 {
	return 50
}
func (h *TestPageHeader) Execute() {
	h.Report.Pdf.FontDir = "D:\\AW\\picorre\\wkpicorre\\to-bi-develop\\font"
	h.Report.Font(Gothic, 14, "")
	h.Report.Cell(15, 20, "REPORT TITLE")
	h.Report.Cell(105, 20, "Page")
	h.Report.Cell(120, 20, strconv.Itoa(h.Report.Page))
	h.Report.Cell(140, 20, "of")
	h.Report.Cell(150, 20, "{#TotalPage#}")
	h.Report.Font(Gothic, 12, "")
	h.Report.Cell(15, 30, "G2HDR")
	e := h.Report.Data.Get(h.Report.DataPos).(*Entx)
	h.Report.Cell(40, 30, e.G2)
	h.Report.Cell(15, 40, "G2Title")
	h.Report.Cell(55, 40, "G1Title")
	h.Report.Cell(95, 40, "DetTitle")
}
func (h *TestPageHeader) SetReport(r *Report) {
	h.Report = r
}

type TestDetail struct {
	Report *Report
}

func (h *TestDetail) GetHeight() float64 {
	return 10
}
func (h *TestDetail) Execute() {
	h.Report.Font(Gothic, 12, "")
	e := h.Report.Data.Get(h.Report.DataPos).(*Entx)
	h.Report.Cell(15, 2, e.G2)
	h.Report.Cell(55, 2, e.G1)
	h.Report.Cell(95, 2, e.Det)
}
func (h *TestDetail) SetReport(r *Report) {
	h.Report = r
}

func (h *TestDetail) BreakCheckBefore(report *Report) int {
	if report.DataPos == 0 {
		//max no
		return 2
	}
	curr := report.Data.Get(report.DataPos).(*Entx)
	before := report.Data.Get(report.DataPos - 1).(*Entx)
	if curr.G2 != before.G2 {
		return 2
	}
	if curr.G1 != before.G1 {
		return 1
	}
	return 0
}
func (h *TestDetail) BreakCheckAfter(report *Report) int {
	if report.DataPos == report.Data.Size()-1 {
		//max no
		return 2
	}
	curr := report.Data.Get(report.DataPos).(*Entx)
	after := report.Data.Get(report.DataPos + 1).(*Entx)
	if curr.G2 != after.G2 {
		return 2
	}
	if curr.G1 != after.G1 {
		return 1
	}
	return 0
}
type TestBeforeGroup2 struct {
	Report *Report
}

func (h *TestBeforeGroup2) GetHeight() float64 {
	return 10
}
func (h *TestBeforeGroup2) Execute() {
	h.Report.Font(Gothic, 12, "")
	e := h.Report.Data.Get(h.Report.DataPos).(*Entx)
	h.Report.Cell(15, 2, "G2HDR")
	h.Report.Cell(40, 2, e.G2)

}
func (h *TestBeforeGroup2) SetReport(r *Report) {
	h.Report = r
}
type TestAfterGroup1 struct {
	Report *Report
}

func (h *TestAfterGroup1) GetHeight() float64 {
	return 10
}
func (h *TestAfterGroup1) Execute() {
	h.Report.Font(Gothic, 12, "")
	e := h.Report.Data.Get(h.Report.DataPos).(*Entx)
	h.Report.Cell(15, 2, "G1Sum")
	h.Report.Cell(40, 2, e.G1)
}
func (h *TestAfterGroup1) SetReport(r *Report) {
	h.Report = r
}
type TestAfterGroup2 struct {
	Report *Report
}

func (h *TestAfterGroup2) GetHeight() float64 {
	return 10
}
func (h *TestAfterGroup2) Execute() {
	h.Report.Font(Gothic, 12, "")
	e := h.Report.Data.Get(h.Report.DataPos).(*Entx)
	h.Report.Cell(15, 2, "G2Sum")
	h.Report.Cell(40, 2, e.G2)
	h.Report.NewPageForce=true
	h.Report.Page=0
}
func (h *TestAfterGroup2) SetReport(r *Report) {
	h.Report = r
}
func getData() *df.List {
	res, _ := ioutil.ReadFile("d:\\temp\\testtab.txt")
	lines := strings.Split(string(res), "\r\n")
	list := new(df.List)
	for _, line := range lines {
		cols := strings.Split(line, "\t")
		if len(cols) < 3 {
			continue
		}
		e := new(Entx)
		e.G2 = cols[0]
		e.G1 = cols[1]
		e.Det = cols[2]
		list.Add(e)
	}
	return list
}

type Entx struct {
	G1  string
	G2  string
	Det string
}
