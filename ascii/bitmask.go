package ascii

const (
	HorizontalLine uint8 = 1 << iota
	VerticalLine
	HorizontalBorder
	VerticalBorder
	LeftTopCorner
	RightTopCorner
	LeftBottomCorner
	RightBottomCorner
)

var AsciiBitmask = [131]rune{
	HorizontalLine:   '─',
	HorizontalBorder: '─',
	VerticalLine:     '│',
	VerticalBorder:   '│',

	LeftTopCorner:     '┌',
	RightTopCorner:    '┐',
	LeftBottomCorner:  '└',
	RightBottomCorner: '┘',

	HorizontalLine | VerticalLine: '┼',

	HorizontalBorder | VerticalLine:    '┬',
	HorizontalLine | LeftTopCorner:     '┬',
	HorizontalLine | RightTopCorner:    '┬',
	VerticalBorder | HorizontalLine:    '┴',
	HorizontalLine | LeftBottomCorner:  '┴',
	HorizontalLine | RightBottomCorner: '┴',

	VerticalLine | RightBottomCorner: '┤',
	VerticalLine | RightTopCorner:    '┤',
	VerticalLine | LeftTopCorner:     '├',
	VerticalLine | LeftBottomCorner:  '├',
}

func Set(b, flag uint8) uint8    { return b | flag }
func Clear(b, flag uint8) uint8  { return b &^ flag }
func Toggle(b, flag uint8) uint8 { return b ^ flag }
func Has(b, flag uint8) bool     { return b&flag != 0 }
