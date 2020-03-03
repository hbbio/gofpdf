package gofpdf

import "github.com/hbbio/gofpdf/svg"

// SVGBasicWrite renders the paths encoded in the basic SVG image specified by
// sb. The scale value is used to convert the coordinates in the path to the
// unit of measure specified in New(). The current position (as set with a call
// to SetXY()) is used as the origin of the image. The current line cap style
// (as set with SetLineCapStyle()), line width (as set with SetLineWidth()),
// and draw color (as set with SetDrawColor()) are used in drawing the image
// paths.
func (f *Fpdf) SVGBasicWrite(sb *svg.BasicType, scale float64) {
	originX, originY := f.GetXY()
	var x, y, newX, newY float64
	var cx0, cy0, cx1, cy1 float64
	var path []svg.BasicSegmentType
	var seg svg.BasicSegmentType
	var startX, startY float64
	sval := func(origin float64, arg int) float64 {
		return origin + scale*seg.Arg[arg]
	}
	xval := func(arg int) float64 {
		return sval(originX, arg)
	}
	yval := func(arg int) float64 {
		return sval(originY, arg)
	}
	val := func(arg int) (float64, float64) {
		return xval(arg), yval(arg + 1)
	}
	for j := 0; j < len(sb.Segments) && f.Ok(); j++ {
		path = sb.Segments[j]
		for k := 0; k < len(path) && f.Ok(); k++ {
			seg = path[k]
			switch seg.Cmd {
			case 'M':
				x, y = val(0)
				startX, startY = x, y
				f.SetXY(x, y)
			case 'L':
				newX, newY = val(0)
				f.Line(x, y, newX, newY)
				x, y = newX, newY
			case 'C':
				cx0, cy0 = val(0)
				cx1, cy1 = val(2)
				newX, newY = val(4)
				f.CurveCubic(x, y, cx0, cy0, newX, newY, cx1, cy1, "D")
				x, y = newX, newY
			case 'Q':
				cx0, cy0 = val(0)
				newX, newY = val(2)
				f.Curve(x, y, cx0, cy0, newX, newY, "D")
				x, y = newX, newY
			case 'H':
				newX = xval(0)
				f.Line(x, y, newX, y)
				x = newX
			case 'V':
				newY = yval(0)
				f.Line(x, y, x, newY)
				y = newY
			case 'Z':
				f.Line(x, y, startX, startY)
				x, y = startX, startY
			default:
				f.SetErrorf("Unexpected path command '%c'", seg.Cmd)
			}
		}
	}
}
