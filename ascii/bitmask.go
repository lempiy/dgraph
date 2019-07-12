package ascii

const (
	RightVector uint16 = 1 << iota
	LeftVector
	TopVector
	BottomVector
	HorizontalBorder
	VerticalBorder
	LeftTopCorner
	RightTopCorner
	LeftBottomCorner
	RightBottomCorner
	StartMarker
	EndMarker
)

var AsciiBitmask = map[uint16]rune{
	RightVector:      '─',
	LeftVector:       '─',
	HorizontalBorder: '─',
	TopVector:        '│',
	BottomVector:     '│',
	VerticalBorder:   '│',

	LeftTopCorner:     '┌',
	RightTopCorner:    '┐',
	LeftBottomCorner:  '└',
	RightBottomCorner: '┘',

	RightVector | TopVector:    '┼',
	LeftVector | TopVector:     '┼',
	RightVector | BottomVector: '┼',
	LeftVector | BottomVector:  '┼',

	LeftVector | RightVector: '─',
	TopVector | BottomVector: '│',

	HorizontalBorder | TopVector | StartMarker:    '┴',
	HorizontalBorder | BottomVector | StartMarker: '┬',
	VerticalBorder | RightVector | StartMarker:    '├',
	VerticalBorder | LeftVector | StartMarker:     '┤',

	HorizontalBorder | TopVector | EndMarker:    '┬',
	HorizontalBorder | BottomVector | EndMarker: '┴',
	VerticalBorder | RightVector | EndMarker:    '┤',
	VerticalBorder | LeftVector | EndMarker:     '├',

	RightVector | LeftTopCorner: '┬',
	LeftVector | LeftTopCorner:  '┬',

	RightVector | RightTopCorner: '┬',
	LeftVector | RightTopCorner:  '┬',

	RightVector | LeftBottomCorner: '┴',
	LeftVector | LeftBottomCorner:  '┴',

	LeftVector | RightBottomCorner:  '┴',
	RightVector | RightBottomCorner: '┴',

	TopVector | RightBottomCorner:    '┤',
	BottomVector | RightBottomCorner: '┤',

	TopVector | RightTopCorner:    '┤',
	BottomVector | RightTopCorner: '┤',

	TopVector | LeftTopCorner:    '├',
	BottomVector | LeftTopCorner: '├',

	TopVector | LeftBottomCorner:    '├',
	BottomVector | LeftBottomCorner: '├',
}
