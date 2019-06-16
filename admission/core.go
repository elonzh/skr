package admission

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/unidoc/unioffice/document"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strings"
)

const TOKEN = "$"

// ConvertToPDF uses go-ole to convert a docx to a PDF using the Word application
func ConvertToPDF(source, destination string) {
	err := ole.CoInitialize(0)
	if err != nil {
		log.Fatalln(err)
	}
	defer ole.CoUninitialize()

	iunk, err := oleutil.CreateObject("Word.Application")
	if err != nil {
		log.Fatalf("error creating Word object: %s", err)
	}

	word := iunk.MustQueryInterface(ole.IID_IDispatch)
	defer word.Release()

	// opening then saving works due to the call to doc.Settings.SetUpdateFieldsOnOpen(true) above

	docs := oleutil.MustGetProperty(word, "Documents").ToIDispatch()
	wordDoc := oleutil.MustCallMethod(docs, "Open", source).ToIDispatch()

	// file format constant comes from https://msdn.microsoft.com/en-us/vba/word-vba/articles/wdsaveformat-enumeration-word
	const wdFormatPDF = 17
	oleutil.MustCallMethod(wordDoc, "SaveAs2", destination, wdFormatPDF)
	oleutil.MustCallMethod(wordDoc, "Close")
	oleutil.MustCallMethod(word, "Quit")
}

func MakeAdmissionNoticeFile(templatePath string, fields map[string]string, outputPath string) {
	doc, err := document.Open(templatePath)
	if err != nil {
		log.Fatalf("error opening Windows Word 2016 document: %s", err)
	}
	var paragraphs []document.Paragraph
	for _, p := range doc.Paragraphs() {
		paragraphs = append(paragraphs, p)
	}

	// This sample document uses structured document tags, which are not common
	// except for in document templates.  Normally you can just iterate over the
	// document's paragraphs.
	for _, sdt := range doc.StructuredDocumentTags() {
		for _, p := range sdt.Paragraphs() {
			paragraphs = append(paragraphs, p)
		}
	}

	for _, table := range doc.Tables() {
		for _, row := range table.Rows() {
			for _, cell := range row.Cells() {
				for _, p := range cell.Paragraphs() {
					paragraphs = append(paragraphs, p)
				}
			}
		}
	}

	fmt.Println("段落数：", len(paragraphs))
	for _, p := range paragraphs {
		// TODO: 实际上 Run 是内联文字, 一段文字可能会分割成多段, 这里为了节约时间通过直接修改 document.xml 确保变量在一个 Run 中
		for _, r := range p.Runs() {
			before := r.Text()
			after := before
			for k, v := range fields {
				f := TOKEN + k + TOKEN
				if strings.Contains(after, f) {
					after = strings.ReplaceAll(after, f, v)
				}
			}
			if before != after {
				fmt.Println("before:", before)
				r.ClearContent()
				r.AddText(after)
				fmt.Println("after :", after)
			}
		}
	}
	if err := doc.SaveToFile(outputPath); err != nil {
		log.Fatalln(err)
	}
}

func SendEmail() {
	//d := gomail.NewDialer("smtp.163.com", 25, "user", "123456")
	// 姓名
	// 缴费截止日期
	// 日期
	// 编号
	// 1909B135
	from := ""
	to := "tcrwaane@sharklasers.com"
	password := ""

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	m.Attach("杨罗罗.xlsx")

	d := gomail.NewDialer("smtp.qq.com", 465, from, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Send the email to Bob, Cora and Dan.
	fmt.Println("开始发送邮件")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func DisableUniofficeWatermark() {
	flag.Bool("test.v", true, "")
	err := flag.CommandLine.Set("test.v", "1")
	if err != nil {
		log.Fatalln(err)
	}
}

func Run() {
	DisableUniofficeWatermark()
	fields := map[string]string{
		"编号":     "1909B000",
		"姓名":     "张三",
		"缴费截止日期": "0000年6月7日",
		"日期":     "0000年6月1日",
	}
	//templatePath := "D:/i/Desktop/admission_notice/direct_template.docx"
	templatePath := "D:/i/Desktop/admission_notice/conditioned_template.docx"
	//templatePath := "test.docx"
	outputPath := "D:/i/Desktop/admission_notice/MakeAdmissionNoticeFile.docx"
	//outputPath := "MakeAdmissionNoticeFile.docx"
	pdfOutputPath := "D:/i/Desktop/admission_notice/MakeAdmissionNoticeFile.pdf"
	//pdfOutputPath := "MakeAdmissionNoticeFile.pdf"
	//SimpleMakeAdmissionNoticeFile(fields)
	if err := os.Remove(outputPath); err != nil {
		log.Println(err)
	}
	if err := os.Remove(pdfOutputPath); err != nil {
		log.Println(err)
	}
	MakeAdmissionNoticeFile(templatePath, fields, outputPath)
	ConvertToPDF(outputPath, pdfOutputPath)
}
