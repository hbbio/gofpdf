package ttf

import (
	"fmt"

	"github.com/hbbio/gofpdf/test"
)

func ExampleParse() {
	font, err := Parse(test.FontDir() + "/calligra.ttf")
	if err == nil {
		fmt.Printf("Postscript name:  %s\n", font.PostScriptName)
		fmt.Printf("unitsPerEm:       %8d\n", font.UnitsPerEm)
		fmt.Printf("Xmin:             %8d\n", font.Xmin)
		fmt.Printf("Ymin:             %8d\n", font.Ymin)
		fmt.Printf("Xmax:             %8d\n", font.Xmax)
		fmt.Printf("Ymax:             %8d\n", font.Ymax)
	} else {
		fmt.Printf("%s\n", err)
	}
	// Output:
	// Postscript name:  CalligrapherRegular
	// unitsPerEm:           1000
	// Xmin:                 -173
	// Ymin:                 -234
	// Xmax:                 1328
	// Ymax:                  899
}
