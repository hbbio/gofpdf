// Package test provides some helper routines for the test packages of
// gofpdf and its various contributed packages located beneath the contrib
// directory.
package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var gofpdfDir string

func init() {
	setRoot()
}

// setRoot assigns the relative path to the gofpdfDir directory based on current working
// directory
func setRoot() {
	wdStr, err := os.Getwd()
	if err == nil {
		gofpdfDir = ""
		list := strings.Split(filepath.ToSlash(wdStr), "/")
		for j := len(list) - 1; j >= 0 && list[j] != "gofpdf"; j-- {
			gofpdfDir = filepath.Join(gofpdfDir, "..")
		}
	} else {
		panic(err)
	}
}

// ImageFile returns a qualified filename in which the path to the image
// directory is prepended to the specified filename.
func ImageFile(fileStr string) string {
	return filepath.Join(gofpdfDir, "test/image", fileStr)
}

// FontDir returns the path to the font directory.
func FontDir() string {
	return filepath.Join(gofpdfDir, "test/font")
}

// FontFile returns a qualified filename in which the path to the font
// directory is prepended to the specified filename.
func FontFile(fileStr string) string {
	return filepath.Join(FontDir(), fileStr)
}

// TextFile returns a qualified filename in which the path to the text
// directory is prepended to the specified filename.
func TextFile(fileStr string) string {
	return filepath.Join(gofpdfDir, "test/text", fileStr)
}

// PdfDir returns the path to the PDF output directory.
func PdfDir() string {
	return filepath.Join(gofpdfDir, "test/pdf")
}

// PdfFile returns a qualified filename in which the path to the PDF output
// directory is prepended to the specified filename.
func PdfFile(fileStr string) string {
	return filepath.Join(PdfDir(), fileStr)
}

// Filename returns a qualified filename in which the example PDF directory
// path is prepended and the suffix ".pdf" is appended to the specified
// filename.
func Filename(baseStr string) string {
	return PdfFile(baseStr + ".pdf")
}

// Summary generates a predictable report for use by test examples. If the
// specified error is nil, the filename delimiters are normalized and the
// filename printed to standard output with a success message. If the specified
// error is not nil, its String() value is printed to standard output.
func Summary(err error, fileStr string) {
	if err == nil {
		fileStr = filepath.ToSlash(fileStr)
		fmt.Printf("Successfully generated %s\n", fileStr)
	} else {
		fmt.Println(err)
	}
}
