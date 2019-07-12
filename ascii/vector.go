package ascii

type Direction int

const (
	DirectionTop = iota
	DirectionBottom
	DirectionRight
	DirectionLeft
)

func getXVertexDirection(x1 int, x2 int) Direction {
	if x1 < x2 {
		return DirectionRight
	}
	return DirectionLeft
}

func getYVertexDirection(y1 int, y2 int) Direction {
	if y1 < y2 {
		return DirectionBottom
	}
	return DirectionTop
}

func getVectorDirection(x1, y1, x2, y2 int) Direction {
	if y1 == y2 {
		return getXVertexDirection(x1, x2)
	}
	return getYVertexDirection(y1, y2)
}

func getCellCenter(width, height, x1, y1 int) (coords [2]int) {
	coords[0] = x1 + width/2
	coords[1] = y1 + height/2
	return coords
}

func getFromCellEntry(direction Direction, width, height, x1, y1 int) (coords [2]int) {
	switch direction {
	case DirectionTop:
		center := getCellCenter(width, height, x1, y1)
		coords[0] = center[0]
		coords[1] = y1
	case DirectionBottom:
		center := getCellCenter(width, height, x1, y1)
		coords[0] = center[0]
		coords[1] = y1 + height - 1
	case DirectionLeft:
		center := getCellCenter(width, height, x1, y1)
		coords[0] = x1
		coords[1] = center[1]
	case DirectionRight:
		center := getCellCenter(width, height, x1, y1)
		coords[0] = x1 + width - 1
		coords[1] = center[1]
	}
	return
}

func getToCellEntry(direction Direction, width, height, x1, y1 int) (coords [2]int) {
	switch direction {
	case DirectionBottom:
		center := getCellCenter(width, height, x1, y1)
		coords[0] = center[0]
		coords[1] = y1
	case DirectionTop:
		center := getCellCenter(width, height, x1, y1)
		coords[0] = center[0]
		coords[1] = y1 + height - 1
	case DirectionRight:
		center := getCellCenter(width, height, x1, y1)
		coords[0] = x1
		coords[1] = center[1]
	case DirectionLeft:
		center := getCellCenter(width, height, x1, y1)
		coords[0] = x1 + width - 1
		coords[1] = center[1]
	}
	return
}
