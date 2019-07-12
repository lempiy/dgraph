package ascii

import (
	"fmt"
	"github.com/lempiy/dgraph/core"
)

const CellHeight = 3
const MinCellWidth = 3

type geometry interface {
	GetEntryToVector(direction Direction, isTarget bool) [2]int
}

type context2D struct {
	figures map[string]geometry
	canvas  *Canvas
	nodeMap map[string]*core.MatrixNode
}

func newContext2D(c *Canvas, nodeMap map[string]*core.MatrixNode) *context2D {
	return &context2D{
		figures: make(map[string]geometry),
		canvas:  c,
		nodeMap: nodeMap,
	}
}

func DrawAsciiMatrix(mtx *core.Matrix) (c *Canvas, err error) {
	rows := mtx.Height()
	colWidths := getColumnWidths(mtx)
	width := sum(colWidths...) + CellHeight*(mtx.Width()-1)
	height := mtx.Height()*CellHeight*2 - CellHeight
	c = NewCanvas(width, height)
	cX := 0
	ctx := newContext2D(c, mtx.Normalize())
	for x, colWidth := range colWidths {
		for y := 0; y < rows; y++ {
			cY := 2 * y * CellHeight
			node := mtx.GetByCoords(x, y)
			// TODO handle anchor on lines
			switch {
			case node == nil:
			case node.IsAnchor:
				err = ctx.drawCorner(c, cX, cY, colWidth, CellHeight, node.Id, node.Orientation)
			default:
				err = ctx.drawNode(c, cX, cY, colWidth, CellHeight, node)
			}
			if err != nil {
				return
			}
		}
		cX += colWidth + CellHeight
	}
	err = ctx.drawLines(c)
	return
}

func (s *context2D) drawNode(c *Canvas, x1, y1, cellWidth, cellHeight int, node *core.NodeOutput) (err error) {
	w := GetWidthFromTitle(node.Id)
	x1 += (cellWidth - w) / 2
	r := NewRectangle(x1, y1, w, cellHeight, node.Id)
	err = r.Draw(c)
	if err != nil {
		return
	}
	s.figures[node.Id] = r

	return
}

func (s *context2D) drawCorner(c *Canvas, x1, y1, cellWidth, cellHeight int,
	id string, orientation core.AnchorOrientation) (err error) {
	coords := getCellCenter(cellWidth, cellHeight, x1, y1)
	crn := NewCorner(coords[0], coords[1], orientation)
	err = crn.draw(c)
	if err != nil {
		return
	}
	s.figures[id] = crn
	return
}

func (s *context2D) drawLines(c *Canvas) (err error) {
	var all []Polyline
	for id, node := range s.nodeMap {
		var polylines []Polyline
		if node.IsAnchor {
			continue
		}
		figure := s.figures[id]
		rect, ok := figure.(*Rectangle)
		if !ok {
			err = fmt.Errorf("figure for node `%s` is not Box", id)
			return
		}
		polylines, err = s.getNodeBranches(rect, node.NodeOutput)
		all = append(all, polylines...)
	}
	for _, polyline := range all {
		err = polyline.draw(c)
		if err != nil {
			return
		}
	}
	return
}

func getColumnWidths(mtx *core.Matrix) (widths []int) {
	columns := mtx.Width()
	rows := mtx.Height()
	for x := 0; x < columns; x++ {
		columnWidth := MinCellWidth
		for y := 0; y < rows; y++ {
			node := mtx.GetByCoords(x, y)
			if node == nil || node.IsAnchor {
				continue
			}
			candidate := GetWidthFromTitle(node.Id)
			if candidate > columnWidth {
				columnWidth = candidate
			}
		}
		widths = append(widths, columnWidth)
	}
	return
}

func (s *context2D) getNodeBranches(start *Rectangle, node *core.NodeOutput) (result []Polyline, err error) {
	n, ok := s.nodeMap[node.Id]
	if !ok {
		err = fmt.Errorf("node with id `%s` not found on map", node.Id)
		return
	}
	incomes, err := s.getAllIncomes(n)
	if err != nil {
		return
	}
	for _, income := range incomes {
		var poly Polyline
		poly, err = s.diveToNodeIncome(n, income)
		if err != nil {
			return
		}
		result = append(result, poly)
	}
	return
}

func (s *context2D) getAllIncomes(source *core.MatrixNode) (result []*core.MatrixNode, err error) {
	for _, incomeId := range source.RenderIncomes {
		target, ok := s.nodeMap[incomeId]
		if !ok {
			err = fmt.Errorf("income with id `%s` not found on map", source.Id)
			return
		}
		result = append(result, target)
	}
	return
}

func (s *context2D) diveToNodeIncome(source *core.MatrixNode, target *core.MatrixNode) (poly Polyline, err error) {
	first, err := s.getLine(source, target)
	if err != nil {
		return
	}
	poly = append(poly, &first)
	for {
		if !target.IsAnchor {
			break
		}
		source = target
		var incomes []*core.MatrixNode
		incomes, err = s.getAllIncomes(source)
		if err != nil {
			return
		}
		if len(incomes) != 1 {
			err = fmt.Errorf("anchor `%s` has more then 1 income", source.Id)
			return
		}
		target = incomes[0]
		var l Line
		l, err = s.getLine(source, target)
		if err != nil {
			return
		}
		poly = append(poly, &l)
	}
	return
}

func (s *context2D) getLine(source *core.MatrixNode, target *core.MatrixNode) (line Line, err error) {
	direction := getVectorDirection(source.X, source.Y, target.X, target.Y)
	src, ok := s.figures[source.Id]
	if !ok {
		err = fmt.Errorf("source figure with id `%s` not found on canvas", source.Id)
		return
	}
	sourceCoords := src.GetEntryToVector(direction, false)
	line.x1, line.y1 = sourceCoords[0], sourceCoords[1]
	tgt, ok := s.figures[target.Id]
	if !ok {
		err = fmt.Errorf("target figure with id `%s` not found on canvas", source.Id)
		return
	}
	if !source.IsAnchor {
		line.markerStart = getCorrectMarker(false)
	}
	targetCoords := tgt.GetEntryToVector(direction, true)
	line.x2, line.y2 = targetCoords[0], targetCoords[1]
	line.direction = direction
	if !target.IsAnchor {
		line.markerEnd = getCorrectMarker(true)
	}

	return
}

func getCorrectMarker(isTarget bool) uint16 {
	if isTarget {
		return EndMarker
	}
	return StartMarker
}

func sum(num ...int) (sum int) {
	for _, n := range num {
		sum += n
	}
	return
}
