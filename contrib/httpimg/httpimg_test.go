package httpimg_test

import (
	"github.com/hbbio/gofpdf"
	"github.com/hbbio/gofpdf/contrib/httpimg"
	"github.com/hbbio/gofpdf/test"
)

func ExampleRegister() {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetFillColor(200, 200, 220)
	pdf.AddPage()

	url := "https://github.com/hbbio/gofpdf/raw/master/image/logo_gofpdf.jpg?raw=true"
	httpimg.Register(pdf, url, "")
	pdf.Image(url, 15, 15, 267, 0, false, "", 0, "")
	fileStr := test.Filename("contrib_httpimg_Register")
	err := pdf.OutputFileAndClose(fileStr)
	test.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../test/pdf/contrib_httpimg_Register.pdf
}
