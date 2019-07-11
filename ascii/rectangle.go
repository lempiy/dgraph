package ascii

const MaxWidth = 11

type Rectangle struct {
	x1    uint
	y1    uint
	x2    uint
	y2    uint
	title string
}

func GetWidthFromTitle(title string) int {
	if len(title)+2 > MaxWidth {
		return MaxWidth
	}
	width := len(title) + 2
	if width%2 == 0 {
		width += 1
	}
	return width
}

func NewRectangle(x1, y1, x2, y2 uint) *Rectangle {

}
