package ascii

import (
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
}

func newContext2D(c *Canvas) *context2D {
	return &context2D{
		figures: make(map[string]geometry),
		canvas:  c,
	}
}

func DrawAsciiMatrix(mtx *core.Matrix) (c *Canvas, err error) {
	rows := mtx.Height()
	colWidths := getColumnWidths(mtx)
	width := sum(colWidths...) + CellHeight*(mtx.Width()-1)
	height := mtx.Height()*CellHeight*2 - CellHeight
	c = NewCanvas(width, height)
	cX := 0
	ctx := newContext2D(c)
	for x, colWidth := range colWidths {
		for y := 0; y < rows; y++ {
			cY := 2 * y * CellHeight
			node := mtx.GetByCoords(x, y)
			// TODO handle anchor on lines
			if node == nil || node.IsAnchor {
				continue
			} else {
				err = ctx.drawNode(c, cX, cY, colWidth, CellHeight, node.Id, node.Id)
			}
			if err != nil {
				return
			}
		}
		cX += colWidth + CellHeight
	}
	return
}

func (s *context2D) drawNode(c *Canvas, x1, y1, cellWidth, cellHeight int, id, title string) (err error) {
	w := GetWidthFromTitle(title)
	x1 += (cellWidth - w) / 2
	r := NewRectangle(x1, y1, w, cellHeight, title)
	err = r.Draw(c)
	if err != nil {
		return
	}
	s.figures[id] = r
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

func sum(num ...int) (sum int) {
	for _, n := range num {
		sum += n
	}
	return
}
