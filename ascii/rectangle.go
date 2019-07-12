package ascii

const MaxWidth = 18
const Padding = 1

type Rectangle struct {
	x1     int
	y1     int
	width  int
	height int
	title  string
}

func GetWidthFromTitle(title string) int {
	if len(title)+2 > MaxWidth {
		return MaxWidth
	}
	width := len(title) + 2
	if width%2 == 0 {
		width += 1
	}
	return width + Padding*2
}

func NewRectangle(x1, y1, width, height int, title string) *Rectangle {
	return &Rectangle{
		x1:     x1,
		y1:     y1,
		width:  width,
		height: height,
		title:  title,
	}
}

func (r *Rectangle) Draw(canvas *Canvas) (err error) {
	err = r.drawCorners(canvas)
	if err != nil {
		return
	}
	err = r.drawBorders(canvas)
	if err != nil {
		return
	}
	err = r.drawTitle(canvas)
	if err != nil {
		return
	}
	return
}

func (r *Rectangle) drawTitle(canvas *Canvas) (err error) {
	coords := getCellCenter(r.width, r.height, r.x1, r.y1)
	t := r.truncateTitle()
	halfLength := len(t) / 2
	x := coords[0] - halfLength
	y := coords[1]
	for _, rn := range t {
		err = canvas.drawLetter(x, y, rn)
		if err != nil {
			return
		}
		x++
	}
	return
}

func (r *Rectangle) truncateTitle() string {
	max := MaxWidth - 2 - Padding*2
	if len(r.title) > max {
		var title []rune
		for _, rn := range r.title {
			if len(title) == max-3 {
				break
			}
			title = append(title, rn)
		}
		title = append(title, '.', '.', '.')
		return string(title)
	}
	return r.title
}

func (r *Rectangle) drawBorders(canvas *Canvas) (err error) {
	err = r.drawTopBorder(canvas)
	if err != nil {
		return
	}
	err = r.drawBottomBorder(canvas)
	if err != nil {
		return
	}
	err = r.drawLeftBorder(canvas)
	if err != nil {
		return
	}
	err = r.drawRightBorder(canvas)
	if err != nil {
		return
	}
	return
}

func (r *Rectangle) drawTopBorder(canvas *Canvas) (err error) {
	x1 := r.x1 + 1
	y := r.y1
	for x := x1; x < r.x1+r.width-1; x++ {
		err = canvas.drawPixel(x, y, HorizontalBorder)
		if err != nil {
			return
		}
	}
	return
}

func (r *Rectangle) drawBottomBorder(canvas *Canvas) (err error) {
	x1 := r.x1 + 1
	y := r.y1 + r.height - 1
	for x := x1; x < r.x1+r.width-1; x++ {
		err = canvas.drawPixel(x, y, HorizontalBorder)
		if err != nil {
			return
		}
	}
	return
}

func (r *Rectangle) drawLeftBorder(canvas *Canvas) (err error) {
	y1 := r.y1 + 1
	x := r.x1
	for y := y1; y < r.y1+r.height-1; y++ {
		err = canvas.drawPixel(x, y, VerticalBorder)
		if err != nil {
			return
		}
	}
	return
}

func (r *Rectangle) drawRightBorder(canvas *Canvas) (err error) {
	y1 := r.y1 + 1
	x := r.x1 + r.width - 1
	for y := y1; y < r.y1+r.height-1; y++ {
		err = canvas.drawPixel(x, y, VerticalBorder)
		if err != nil {
			return
		}
	}
	return
}

func (r *Rectangle) GetEntryToVector(direction Direction, isTarget bool) [2]int {
	if isTarget {
		return getToCellEntry(direction, r.width, r.height, r.x1, r.y1)
	}
	return getFromCellEntry(direction, r.width, r.height, r.x1, r.y1)
}

func (r *Rectangle) drawCorners(canvas *Canvas) (err error) {
	err = canvas.drawPixel(r.x1, r.y1, LeftTopCorner)
	if err != nil {
		return
	}
	err = canvas.drawPixel(r.x1+r.width-1, r.y1, RightTopCorner)
	if err != nil {
		return
	}
	err = canvas.drawPixel(r.x1, r.y1+r.height-1, LeftBottomCorner)
	if err != nil {
		return
	}
	err = canvas.drawPixel(r.x1+r.width-1, r.y1+r.height-1, RightBottomCorner)
	if err != nil {
		return
	}
	return
}
