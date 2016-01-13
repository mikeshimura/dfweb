package pdf

import (
	"fmt"
	"github.com/signintech/gopdf"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)
const(
	Gothic="IPAexゴシック"
	GothicBold="IPAexゴシック太字"
	Mincho="IPAex明朝"
	MinchoBold="IPAex明朝太字"
)

var ConvPtMm float64 = 2.834645669

type DfPdf struct {
	Pdf     *gopdf.GoPdf
	Doc     string
	FontDir string
}

func (p *DfPdf) ReadFile(fileName string) error {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	p.Doc = strings.Replace(string(buf), "\r", "", -1)
	var UTF8_BOM = []byte{239, 187, 191}
	if p.Doc[0:3] == string(UTF8_BOM) {
		fmt.Println("BOM found")
		p.Doc = p.Doc[3:]
	}
	return nil
}
func (p *DfPdf) SetFontDirFromEnv() {
	p.FontDir = os.Getenv("FONTDIR")
}
func (p *DfPdf) Execute() {
	lines := strings.Split(p.Doc, "\n")
	for _, line := range lines {
		//fmt.Println("line:" + line)
		eles := strings.Split(line, "\t")
		//fmt.Printf("eles[0]:%v:len %v\n",eles[0],len(eles[0]))
		switch eles[0] {
		case "P":
			p.Page(line, eles)
		case "NP":
			p.NewPage(line, eles)
		case "F":
			p.Font(line, eles)
		case "C", "C1", "CR":
			p.Cell(line, eles)
		case "M":
			p.Move(line, eles)
		default:
			fmt.Println("default:" + line + ":")
		}
	}
}
func (p *DfPdf) AddFont() {
	err := p.Pdf.AddTTFFont(Gothic, p.FontDir+"\\ipaexg.ttf")
	err = p.Pdf.AddTTFFont(GothicBold, p.FontDir+"\\ipaexg-bold.ttf")
	err = p.Pdf.AddTTFFont(Mincho, p.FontDir+"\\ipaexm.ttf")
	if err != nil {
		panic(err)
	}
}
func (p *DfPdf) Page(line string, eles []string) {
	p.Pdf = new(gopdf.GoPdf)
	switch eles[0] {
	case "P":
		CheckLength(line, eles, 3)
		switch eles[1] {
		case "A4":
			if eles[2] == "P" {
				p.Start(595.28, 841.89)
			} else if eles[2] == "L" {
				p.Start(841.89, 595.28)
			} else {
				panic("Page Orientationは PかLです")
			}
		default:
			panic("このサイズはサポートされていません:" + eles[1])
		}
	case "P1":
		CheckLength(line, eles, 3)
		p.Start(ParseFloatPanic(eles[1])*ConvPtMm, ParseFloatPanic(eles[2])*ConvPtMm)
	}
	p.AddFont()
	p.Pdf.AddPage()
}
func (p *DfPdf) NewPage(line string, eles []string) {
	p.Pdf.AddPage()
}
func (p *DfPdf) Start(w float64, h float64) {
	p.Pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: w, H: h}}) //595.28, 841.89 = A4
}
func (p *DfPdf) Font(line string, eles []string) {
	CheckLength(line, eles, 4)
	err := p.Pdf.SetFont(eles[1], eles[2], AtoiPanic(eles[3]))
	if err != nil {
		panic(err)
	}
}
func CheckLength(line string, eles []string, no int) {
	if len(eles) < no {
		panic("項目が不足です。:" + line)
	}
}
func (p *DfPdf) Cell(line string, eles []string) {
	switch eles[0] {
	case "C":
		CheckLength(line, eles, 6)
		err := p.Pdf.SetFont(eles[1], "", AtoiPanic(eles[2]))
		if err != nil {
			panic(err)
		}
		p.MoveSub(eles[3], eles[4])
		p.Pdf.Cell(nil, eles[5])
	case "C1":
		CheckLength(line, eles, 4)
		p.MoveSub(eles[1], eles[2])
		p.Pdf.Cell(nil, eles[3])
	case "CR":
		CheckLength(line, eles, 5)
		tw, err := p.Pdf.MeasureTextWidth(eles[4])
		if err != nil {
			panic(err)
		}
		x := ParseFloatPanic(eles[1]) * ConvPtMm
		y := ParseFloatPanic(eles[2]) * ConvPtMm
		w := ParseFloatPanic(eles[3]) * ConvPtMm
		finalx := x + w - tw
		p.Pdf.SetX(finalx)
		p.Pdf.SetY(y)
		p.Pdf.Cell(nil, eles[4])
	}
}
func (p *DfPdf) Move(line string, eles []string) {
	CheckLength(line, eles, 3)
	p.MoveSub(eles[1], eles[2])
}
func (p *DfPdf) MoveSub(sx string, sy string) {
	p.Pdf.SetX(ParseFloatPanic(sx) * ConvPtMm)
	p.Pdf.SetY(ParseFloatPanic(sy) * ConvPtMm)
}
func AtoiPanic(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(s + "は整数ではありません")
	}
	return i
}
func ParseFloatPanic(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(s + "は数値ではありません")
	}
	return f
}
