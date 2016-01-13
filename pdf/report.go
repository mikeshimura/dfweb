package pdf

import (
	"bytes"
	"fmt"
	"github.com/mikeshimura/dbflute/df"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	PageHeader   = "PageHeader"
	PageFooter   = "PageFooter"
	Detail       = "Detail"
	Summary      = "Summary"
	GroupHeader  = "GroupHeader"
	GroupSummary = "GroupSummary"
)

type Report struct {
	Data         *df.List
	DataPos      int
	Bands        map[string]*Band
	Pdf          *DfPdf
	PageX        float64
	PageY        float64
	CurrY        float64
	MaxY         float64
	MaxGroup     int
	Page         int
	PageTotal    bool
	NewPageForce bool
}

func (r *Report) Execute(filename string) {
	r.Page = 1
	r.CurrY = 0
	r.ExecutePageHeader()
	r.AddLine("v\tPAGE\t" + strconv.Itoa(r.Page))
	for r.DataPos = 0; r.DataPos < r.Data.Size(); r.DataPos++ {
		r.ExecuteDetail()
	}
	r.ExecuteSummary()
	r.ExecutePageFooter()
	r.ReplacePageTotal()
	r.Pdf.Execute()
	r.Pdf.Pdf.WritePdf(filename)
}
func (r *Report) ReplacePageTotal() {
	if r.PageTotal == false {
		return
	}
	lines := strings.Split(r.Pdf.Doc, "\n")
	list := new(df.List)
	for i, line := range lines {
		if len(line) < 8 {
			continue
		}
		if line[0:7] == "v\tPAGE\t" {
			h := new(pagehist)
			h.line = i
			h.page = AtoiPanic(line[7:])
			list.Add(h)
			fmt.Printf("hist %v \n", h)
		}
	}
	for i, line := range lines {
		if strings.Index(line, "{#TotalPage#}") > -1 {
			total := r.getTotalValue(i, list)
			fmt.Printf("total :%v\n", total)
			lines[i] = strings.Replace(lines[i], "{#TotalPage#}", strconv.Itoa(total), -1)
		}
	}
	buf := new(bytes.Buffer)
	for _, line := range lines {
		buf.WriteString(line + "\n")
	}
	r.Pdf.Doc = buf.String()
}
func (r *Report) getTotalValue(lineno int, list *df.List) int {
	count := 0
	page := 0
	for i, l := range list.GetAsArray() {
		if l.(*pagehist).line >= lineno {
			count = i
			break
		}
	}
	for i := count; i < list.Size(); i++ {
		newpage := list.Get(i).(*pagehist).page
		if newpage <= page {
			return page
		}
		page = newpage
		fmt.Printf("page :%v\n", page)
	}
	return page
}

type pagehist struct {
	line int
	page int
}

func (r *Report) PageBreak() {
	r.ExecutePageFooter()
	r.AddLine("NP")
	r.Page++
	r.CurrY = 0
	r.ExecutePageHeader()
	r.AddLine("v\tPAGE\t" + strconv.Itoa(r.Page))
}
func (r *Report) PageBreakCheck(height float64) {
	if r.CurrY+height > r.MaxY {
		r.PageBreak()
	}
}
func (r *Report) ExecutePageFooter() {
	h := r.Bands[PageFooter]
	if h != nil {
		(*h).Execute()
		r.CurrY += (*h).GetHeight()
	}
}
func (r *Report) ExecuteSummary() {
	h := r.Bands[Summary]
	if h != nil {
		r.PageBreakCheck((*h).GetHeight())
		(*h).Execute()
		r.CurrY += (*h).GetHeight()
	}
}
func (r *Report) ExecutePageHeader() {
	h := r.Bands[PageHeader]
	if h != nil {
		(*h).Execute()
		r.CurrY += (*h).GetHeight()
	}
}
func (r *Report) ExecuteGroupHeader(level int) {
	for l := level; l > 0; l-- {
		h := r.Bands[GroupHeader+strconv.Itoa(l)]
		if h != nil {
			r.PageBreakCheck((*h).GetHeight())
			(*h).Execute()
			r.CurrY += (*h).GetHeight()
		}
	}
}
func (r *Report) ExecuteGroupSummary(level int) {
	for l := 1; l <= level; l++ {
		h := r.Bands[GroupSummary+strconv.Itoa(l)]
		if h != nil {
			r.PageBreakCheck((*h).GetHeight())
			(*h).Execute()
			r.CurrY += (*h).GetHeight()
		}
	}
}
func (r *Report) ExecuteDetail() {
	h := r.Bands[Detail]
	if h != nil {
		if r.NewPageForce {
			fmt.Println("NewPageForce")
			r.PageBreak()
			r.NewPageForce = false
		}
		var deti interface{} = *h
		if r.MaxGroup > 0 {
			bfr := reflect.ValueOf(deti).MethodByName("BreakCheckBefore")
			if bfr.IsValid() == false {
				panic("BreakCheckBefore functionがDetailにありません")
			}
			res := bfr.Call([]reflect.Value{reflect.ValueOf(r)})
			level := res[0].Int()
			if level > 0 {
				r.ExecuteGroupHeader(int(level))
			}
		}
		r.PageBreakCheck((*h).GetHeight())
		(*h).Execute()
		r.CurrY += (*h).GetHeight()
		aft := reflect.ValueOf(deti).MethodByName("BreakCheckAfter")
		if aft.IsValid() == false {
			panic("BreakCheckAfter functionがDetailにありません")
		}
		res := aft.Call([]reflect.Value{reflect.ValueOf(r)})
		level := res[0].Int()
		if level > 0 {
			r.ExecuteGroupSummary(int(level))
		}
	}
}
func (r *Report) RegisterBand(band *Band, name string) {
	(*band).SetReport(r)
	r.Bands[name] = band
}
func (r *Report) RegisterGroupBand(band *Band, name string, level int) {
	(*band).SetReport(r)
	r.Bands[name+strconv.Itoa(level)] = band
	if r.MaxGroup < level {
		r.MaxGroup = level
	}
}
func Ftoa(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
func (r *Report) AddLine(s string) {
	r.Pdf.Doc += s + "\n"
}
func (r *Report) Font(fontName string, size int, style string) {
	r.AddLine("F\t" + fontName + "\t" + style + "\t" + strconv.Itoa(size))
}
func (r *Report) Cell(x float64, y float64, content string) {
	r.AddLine("C1\t" + Ftoa(x) + "\t" + Ftoa(r.CurrY+y) + "\t" + content)
}
func (r *Report) Var(name string, val string) {
	r.AddLine("V\t" + name + "\t" + val)
}
func (r *Report) SetPage(size string, orientation string) {
	switch size {
	case "A4":
		switch orientation {
		case "P":
			r.AddLine("P\tA4\tP")
			r.PageX = 210
			r.PageY = 297
		case "L":
			r.AddLine("P\tA4\tL")
			r.PageX = 297
			r.PageY = 210
		}
	}
}
func (r *Report) SaveDoc(fileName string) {
	ioutil.WriteFile(fileName, []byte(r.Pdf.Doc), os.ModePerm)
}

type Band interface {
	GetHeight() float64
	Execute()
	SetReport(r *Report)
}

func CreateReport() *Report {
	report := new(Report)
	report.Bands = make(map[string]*Band)
	report.Pdf = new(DfPdf)
	return report
}
