/*
 * Copyright (c) 2015 Kurt Jung (Gmail: kurt.w.jung)
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package gofpdf

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type sortType struct {
	length int
	less   func(int, int) bool
	swap   func(int, int)
}

func (s *sortType) Len() int {
	return s.length
}

func (s *sortType) Less(i, j int) bool {
	return s.less(i, j)
}

func (s *sortType) Swap(i, j int) {
	s.swap(i, j)
}

func gensort(Len int, Less func(int, int) bool, Swap func(int, int)) {
	sort.Sort(&sortType{length: Len, less: Less, swap: Swap})
}

func writeBytes(leadStr string, startPos int, sl []byte) {
	var pos, max int
	var b byte
	fmt.Printf("%s %07x", leadStr, startPos)
	max = len(sl)
	for pos < max {
		fmt.Printf(" ")
		for k := 0; k < 8; k++ {
			if pos < max {
				fmt.Printf(" %02x", sl[pos])
			} else {
				fmt.Printf("   ")
			}
			pos++
		}
	}
	fmt.Printf("  |")
	pos = 0
	for pos < max {
		b = sl[pos]
		if b < 32 || b >= 128 {
			b = '.'
		}
		fmt.Printf("%c", b)
		pos++
	}
	fmt.Printf("|\n")
}

func checkBytes(pos int, sl1, sl2 []byte, printDiff bool) (eq bool) {
	eq = bytes.Equal(sl1, sl2)
	if !eq && printDiff {
		writeBytes("<", pos, sl1)
		writeBytes(">", pos, sl2)
	}
	return
}

// CompareBytes compares the bytes referred to by sl1 with those referred to by
// sl2. Nil is returned if the buffers are equal, otherwise an error.
func CompareBytes(sl1, sl2 []byte, printDiff bool) (err error) {
	var posStart, posEnd, len1, len2, length int
	var diffs bool

	len1 = len(sl1)
	len2 = len(sl2)
	length = len1
	if length > len2 {
		length = len2
	}
	for posStart < length-1 {
		posEnd = posStart + 16
		if posEnd > length {
			posEnd = length
		}
		if !checkBytes(posStart, sl1[posStart:posEnd], sl2[posStart:posEnd], printDiff) {
			diffs = true
		}
		posStart = posEnd
	}
	if diffs {
		err = fmt.Errorf("documents are different")
	}
	return
}

// ComparePDFs reads and compares the full contents of the two specified
// readers byte-for-byte. Nil is returned if the buffers are equal, otherwise
// an error.
func ComparePDFs(rdr1, rdr2 io.Reader, printDiff bool) (err error) {
	var b1, b2 *bytes.Buffer
	_, err = b1.ReadFrom(rdr1)
	if err == nil {
		_, err = b2.ReadFrom(rdr2)
		if err == nil {
			err = CompareBytes(b1.Bytes(), b2.Bytes(), printDiff)
		}
	}
	return
}

// ComparePDFFiles reads and compares the full contents of the two specified
// files byte-for-byte. Nil is returned if the file contents are equal, or if
// the second file is missing, otherwise an error.
func ComparePDFFiles(file1Str, file2Str string, printDiff bool) (err error) {
	var sl1, sl2 []byte
	sl1, err = ioutil.ReadFile(file1Str)
	if err == nil {
		sl2, err = ioutil.ReadFile(file2Str)
		if err == nil {
			err = CompareBytes(sl1, sl2, printDiff)
		} else {
			// Second file is missing; treat this as success
			err = nil
		}
	}
	return
}

// referenceCompare compares the specified file with the file's reference copy
// located in the 'reference' subdirectory. All bytes of the two files are
// compared except for the value of the /CreationDate field in the PDF. This
// function succeeds if both files are equivalent except for their
// /CreationDate values or if the reference file does not exist.
func referenceCompare(fileStr string) (err error) {
	var refFileStr, refDirStr, dirStr, baseFileStr string
	dirStr, baseFileStr = filepath.Split(fileStr)
	refDirStr = filepath.Join(dirStr, "reference")
	err = os.MkdirAll(refDirStr, 0755)
	if err == nil {
		refFileStr = filepath.Join(refDirStr, baseFileStr)
		err = ComparePDFFiles(fileStr, refFileStr, false)
	}
	return
}

// SummaryCompare generates a predictable report for use by test examples. If
// the specified error is nil, the generated file is compared with a reference
// copy for byte-for-byte equality. If the files match, then the filename
// delimiters are normalized and the filename printed to standard output with a
// success message. If the files do not match, this condition is reported on
// standard output. If the specified error is not nil, its String() value is
// printed to standard output.
func SummaryCompare(err error, fileStr string) {
	if err == nil {
		err = referenceCompare(fileStr)
	}
	if err == nil {
		fileStr = filepath.ToSlash(fileStr)
		fmt.Printf("Successfully generated %s\n", fileStr)
	} else {
		fmt.Println(err)
	}
}
