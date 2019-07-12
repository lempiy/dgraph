package ascii

type Line struct {
	x1          int
	y1          int
	markerStart uint16
	x2          int
	y2          int
	markerEnd   uint16
	direction   Direction
}

func (l *Line) draw(c *Canvas) (err error) {
	var startX, startY, endX, endY int
	var startMarker, endMarker uint16
	switch l.direction {
	case DirectionLeft, DirectionTop:
		startX, startY = l.x2, l.y2
		endX, endY = l.x1, l.y1
		startMarker = l.markerEnd
		endMarker = l.markerStart
	case DirectionRight, DirectionBottom:
		startX, startY = l.x1, l.y1
		endX, endY = l.x2, l.y2
		startMarker = l.markerStart
		endMarker = l.markerEnd
	}
	switch l.direction {
	case DirectionTop:
		err = drawVerticalLine(c, startX, startY, endY, TopVector, startMarker, endMarker)
	case DirectionBottom:
		err = drawVerticalLine(c, startX, startY, endY, BottomVector, startMarker, endMarker)
	case DirectionRight:
		err = drawHorizontalLine(c, startX, endX, startY, RightVector, startMarker, endMarker)
	case DirectionLeft:
		err = drawHorizontalLine(c, startX, endX, startY, LeftVector, startMarker, endMarker)
	}
	return
}

func drawHorizontalLine(c *Canvas, startX, endX, y int, flag, startMarker, endMarker uint16) (err error) {
	for x := startX; x <= endX; x++ {
		f := flag
		if x == startX {
			f |= startMarker
		} else if x == endX {
			f |= endMarker
		}
		err = c.drawPixel(x, y, f)
		if err != nil {
			return
		}
	}
	return
}

func drawVerticalLine(c *Canvas, x, startY, endY int, flag, startMarker, endMarker uint16) (err error) {
	for y := startY; y <= endY; y++ {
		f := flag
		if y == startY {
			f |= startMarker
		} else if y == endY {
			f |= endMarker
		}
		err = c.drawPixel(x, y, f)
		if err != nil {
			return
		}
	}
	return
}

type Polyline []*Line

func (p *Polyline) draw(c *Canvas) (err error) {
	for _, line := range *p {
		err = line.draw(c)
		if err != nil {
			return
		}
	}
	return
}
