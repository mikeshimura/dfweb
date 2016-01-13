package pdf

import (
//	"fmt"
//	"github.com/signintech/gopdf"
	"testing"
)

func TestPdf(t *testing.T) {
	pdf:=new( DfPdf)
	pdf.FontDir="D:\\AW\\picorre\\wkpicorre\\to-bi-develop\\font"
	pdf.ReadFile("d:\\temp\\pdftest.txt")
	pdf.Execute()
	pdf.Pdf.WritePdf("pdftest.pdf")

//	pdf := gopdf.GoPdf{}
//	
//	pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
//	err := pdf.AddTTFFont("IPAexゴシック", "../../../../../font/ipaexg.ttf")
//	err = pdf.AddTTFFont("IPAex明朝", "../../../../../font/ipaexm.ttf")
//	if err != nil {
//		panic(err)
//	}
//	pdf.AddPage()
//
//	err = pdf.SetFont("IPAexゴシック", "", 14)
//	if err != nil {
//		panic(err)
//	}
//	pdf.SetX(100)
//	pdf.SetY(100)
//	pdf.Cell(nil, "志村正信")
//	pdf.Curr.Font_ISubset.AddChars("志村正信")
//	w, err := pdf.MeasureTextWidth("志村正信")
//	fmt.Printf("err %v w %v \n", err, w)
//	err = pdf.SetFont("IPAexゴシック", "BI", 14)
//	pdf.SetX(100)
//	pdf.SetY(140)
//	pdf.Cell(nil, "志村正信")
//	pdf.Curr.Font_ISubset.AddChars("志村正信")
//	w, err = pdf.MeasureTextWidth("志村正信")
//	fmt.Printf("err %v w %v \n", err, w)
//	pdf.SetLineType("dashed")
//	pdf.SetLineWidth(0.6)
//	pdf.Line(50, 200, 550, 200)
//	pdf.SetLineType("straight")
//	pdf.SetLineWidth(0.3)
//	pdf.Line(50, 210, 550, 210)
//	pdf.WritePdf("pdftest.pdf")
//		pdf2:=new( DfPdf)
//	pdf2.FontDir="D:\\AW\\picorre\\wkpicorre\\to-bi-develop\\font"
//	pdf2.ReadFile("d:\\temp\\pdftest.txt")
//	pdf2.Execute()
//	pdf2.Pdf.WritePdf("pdftest2.pdf")

}
