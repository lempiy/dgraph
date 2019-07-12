package core

type Matrix struct {
	s [][]*NodeOutput
}

func NewMatrix() *Matrix {
	return new(Matrix)
}

func (m *Matrix) Width() (width int) {
	for _, row := range m.s {
		if len(row) > width {
			width = len(row)
		}
	}
	return
}

func (m *Matrix) Height() int {
	return len(m.s)
}

func (m *Matrix) hasHorizontalCollision(x, y int) bool {
	if len(m.s) == 0 || y >= len(m.s) {
		return false
	}
	row := m.s[y]
	for _, point := range row {
		if point != nil && !m.isAllChildrenOnMatrix(point) {
			return true
		}
	}
	return false
}

func (m *Matrix) hasVerticalCollision(x, y int) bool {
	if x >= m.Width() {
		return false
	}
	for index, row := range m.s {
		if index >= y && x < len(row) && row[x] != nil {
			return true
		}
	}
	return false
}

func (m *Matrix) getFreeRowForColumn(x int) int {
	if m.Height() == 0 {
		return 0
	}
	y := -1
	for index, row := range m.s {
		if len(row) == 0 || x >= len(row) || row[x] == nil {
			y = index
			break
		}
	}
	if y == -1 {
		y = m.Height()
	}
	return y
}

func (m *Matrix) extendHeight(toValue int) {
	for m.Height() < toValue {
		m.s = append(m.s, make([]*NodeOutput, m.Width()))
	}
}

func (m *Matrix) extendWidth(toValue int) {
	for i := range m.s {
		for len(m.s[i]) < toValue {
			m.s[i] = append(m.s[i], nil)
		}
	}
}

func (m *Matrix) insertRowBefore(y int) {
	row := make([]*NodeOutput, m.Width())
	m.s = append(m.s[:y], append([][]*NodeOutput{row}, m.s[y:]...)...)
}

func (m *Matrix) insertColumnBefore(x int) {
	for _, row := range m.s {
		row = append(row[:x], append([]*NodeOutput{nil}, row[x:]...)...)
	}
}

func (m *Matrix) find(callback func(item *NodeOutput) bool) []int {
	for y, row := range m.s {
		for x, point := range row {
			if point == nil {
				continue
			}
			if callback(point) {
				return []int{x, y}
			}
		}
	}
	return nil
}

type findNodeResult struct {
	coords []int
	item   *NodeOutput
}

func (m *Matrix) findNode(callback func(item *NodeOutput) bool) *findNodeResult {
	for y, row := range m.s {
		for x, point := range row {
			if point == nil {
				continue
			}
			if callback(point) {
				return &findNodeResult{
					coords: []int{x, y},
					item:   point,
				}
			}
		}
	}
	return nil
}

func (m *Matrix) GetByCoords(x, y int) *NodeOutput {
	// TODO: remove dat check
	if x >= m.Width() || y >= m.Height() {
		return nil
	}
	return m.s[y][x]
}

func (m *Matrix) insert(x, y int, item *NodeOutput) {
	if m.Height() <= y {
		m.extendHeight(y + 1)
	}
	if m.Width() <= x {
		m.extendWidth(x + 1)
	}
	m.s[y][x] = item
}

func (m *Matrix) isAllChildrenOnMatrix(item *NodeOutput) bool {
	return len(item.Next) == item.ChildrenOnMatrix
}

func (m *Matrix) String() (result string) {
	max := 0
	for _, row := range m.s {
		for _, cell := range row {
			if cell == nil {
				continue
			}
			if len(cell.Id) > max {
				max = len(cell.Id)
			}
		}
	}
	for _, row := range m.s {
		for _, cell := range row {
			if cell == nil {
				result += fillWithSpaces(" ", max)
				result += "│"
				continue
			}
			result += fillWithSpaces(cell.Id, max)
			result += "│"
		}
		result += "\n"
	}
	return result
}

func (m *Matrix) Normalize() map[string]*MatrixNode {
	acc := make(map[string]*MatrixNode)
	for y, row := range m.s {
		for x, item := range row {
			if item != nil {
				acc[item.Id] = &MatrixNode{NodeOutput: item, X: x, Y: y}
			}
		}
	}
	return acc
}

func fillWithSpaces(str string, l int) string {
	for len(str) < l {
		str += " "
	}
	return str
}
