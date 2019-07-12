package ascii

import "github.com/lempiy/dgraph/core"

type Corner struct {
	x           int
	y           int
	orientation core.AnchorOrientation
	rune        rune
	flag        uint16
}

func NewCorner(x, y int, orientation core.AnchorOrientation) *Corner {
	return &Corner{
		x:           x,
		y:           y,
		orientation: orientation,
		rune:        AsciiBitmask[getFlagFromOrientation(orientation)],
		flag:        getFlagFromOrientation(orientation),
	}
}

func (c *Corner) draw(canvas *Canvas) (err error) {
	err = canvas.drawPixel(c.x, c.y, c.flag)
	return
}

func getFlagFromOrientation(orientation core.AnchorOrientation) uint16 {
	switch orientation {
	case core.AnchorOrientationTopLeft:
		return LeftTopCorner
	case core.AnchorOrientationTopRight:
		return RightTopCorner
	case core.AnchorOrientationBottomLeft:
		return LeftBottomCorner
	case core.AnchorOrientationBottomRight:
		return RightBottomCorner
	}
	return 0
}

func (c *Corner) GetEntryToVector(direction Direction, isTarget bool) [2]int {
	if isTarget {
		return c.getEntryToVectorTarget(direction)
	}
	return c.getEntryToVectorSource(direction)
}

func (c *Corner) getEntryToVectorTarget(direction Direction) [2]int {
	switch direction {
	case DirectionLeft:
		return [2]int{c.x + 1, c.y}
	case DirectionRight:
		return [2]int{c.x - 1, c.y}
	case DirectionBottom:
		return [2]int{c.x, c.y - 1}
	case DirectionTop:
		return [2]int{c.x, c.y + 1}
	}
	return [2]int{}
}

func (c *Corner) getEntryToVectorSource(direction Direction) [2]int {
	switch direction {
	case DirectionLeft:
		return [2]int{c.x - 1, c.y}
	case DirectionRight:
		return [2]int{c.x + 1, c.y}
	case DirectionBottom:
		return [2]int{c.x, c.y + 1}
	case DirectionTop:
		return [2]int{c.x, c.y - 1}
	}
	return [2]int{}
}
