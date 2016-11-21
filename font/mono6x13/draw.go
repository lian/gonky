package mono6x13

import (
	"image/color"
	"unicode/utf8"

	"github.com/pbnjay/pixfont"
)

var charNewLine rune
var charSpace rune
var charTab rune

func init() {
	charNewLine, _ = utf8.DecodeRuneInString("\n")
	charSpace, _ = utf8.DecodeRuneInString(" ")
	charTab, _ = utf8.DecodeRuneInString("\t")
}

const Width = 6

func DrawString(dr pixfont.Drawable, x, y int, s string, clr color.Color) (int, int) {
	sx := x
	for _, c := range s {
		switch c {
		case charNewLine:
			x = sx
			y += 12
			break
		case charSpace:
			x += 6
			break
		case charTab:
			x += 6 * 2
			break
		default:
			haveChar, w := Font.DrawRune(dr, x, y, c, clr)
			if haveChar {
				x += w
			}
		}
	}
	return x, y
}
