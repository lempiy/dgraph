package core

import "fmt"

type state struct {
	mtx   *Matrix
	queue *TraverseQueue
	x     int
	y     int
}

type GraphMatrix struct {
	*GraphBasic
}

func NewGraphMatrix(list []NodeInput) (*GraphMatrix, error) {
	g, err := NewGraphBasic(list)
	if err != nil {
		return nil, err
	}
	return &GraphMatrix{
		GraphBasic: g,
	}, nil
}

func (g *GraphMatrix) joinHasUnresolvedIncomes(item *NodeOutput) bool {
	return len(item.PassedIncomes) != len(g.incomes(item.Id))
}

func (g *GraphMatrix) insertOrSkipNodeOnMatrix(item *NodeOutput, state *state, checkCollision bool) (err error) {
	mtx := state.mtx
	if checkCollision && mtx.hasHorizontalCollision(state.x, state.y) {
		mtx.insertRowBefore(state.y)
	}
	mtx.insert(state.x, state.y, item)
	err = g.markIncomesAsPassed(mtx, item)
	return
}

func (g *GraphMatrix) getLowestYAmongIncomes(item *NodeOutput, mtx *Matrix) (y int, err error) {
	incomes := item.PassedIncomes
	if len(incomes) != 0 {
		var items []int
		for _, id := range incomes {
			coords := mtx.find(func(item *NodeOutput) bool {
				return item.Id == id
			})
			if len(coords) != 2 {
				err = fmt.Errorf("cannot find coordinates for passed income: `%s`", id)
			}
			items = append(items, coords[1])
		}
		y = min(items...)
	}
	return
}

func (g *GraphMatrix) processOrSkipNodeOnMatrix(item *NodeOutput, state *state) (isOk bool, err error) {
	mtx, queue := state.mtx, state.queue
	if len(item.PassedIncomes) != 0 {
		state.y, err = g.getLowestYAmongIncomes(item, mtx)
		if err != nil {
			return
		}
	}
	hasLoops := g.hasLoops(item)
	var loopNodes []loopNode
	if hasLoops {
		loopNodes, err = g.handleLoopEdges(item, state)
		if err != nil {
			return
		}
	}
	var needsLoopSkip bool
	if hasLoops && len(loopNodes) == 0 {
		needsLoopSkip = true
	}
	if mtx.hasVerticalCollision(state.x, state.y) || needsLoopSkip {
		queue.push(item)
		return
	}
	err = g.insertOrSkipNodeOnMatrix(item, state, false)
	if err != nil {
		return
	}
	if len(loopNodes) != 0 {
		err = g.insertLoopEdges(item, state, loopNodes)
		if err != nil {
			return
		}
	}
	isOk = true
	return
}

func (g *GraphMatrix) handleLoopEdges(item *NodeOutput, state *state) (loopNodes []loopNode, err error) {
	mtx := state.mtx
	loops := g.loops(item.Id)
	if len(loops) == 0 {
		err = fmt.Errorf("no loops found for node: `%s`", item.Id)
		return
	}
	for _, incomeId := range loops {
		if item.Id == incomeId {
			loopNodes = append(loopNodes, loopNode{
				id:         incomeId,
				node:       item,
				x:          state.x,
				y:          state.y,
				isSelfLoop: true,
			})
			continue
		}
		coords := mtx.find(func(n *NodeOutput) bool {
			return n.Id == incomeId
		})
		if len(coords) != 2 {
			err = fmt.Errorf("loop target `%s` not found on matrix", incomeId)
			return
		}
		node := mtx.getByCoords(coords[0], coords[1])
		if node == nil {
			err = fmt.Errorf("loop target node `%s` not found on matrix", incomeId)
			return
		}
		loopNodes = append(loopNodes, loopNode{
			id:   incomeId,
			node: node,
			x:    coords[0],
			y:    coords[1],
		})
	}
	var skip bool
	for _, income := range loopNodes {
		checkY := 0
		if income.y != 0 {
			checkY = income.y - 1
		}
		skip = mtx.hasVerticalCollision(state.x, checkY)
		if skip {
			break
		}
	}
	if skip {
		loopNodes = nil
	}
	return
}

func (g *GraphMatrix) hasLoops(item *NodeOutput) bool {
	return len(g.loops(item.Id)) != 0
}

func (g *GraphMatrix) insertLoopEdges(item *NodeOutput, state *state, loopNodes []loopNode) (err error) {
	mtx, initialX, initialY := state.mtx, state.x, state.y
	for _, income := range loopNodes {
		id, node, renderIncomeId := income.id, income.node, item.Id
		if income.isSelfLoop {
			state.x, state.y = initialX+1, initialY
			selfLoopId := fmt.Sprintf("%s-self", id)
			renderIncomeId = selfLoopId
			err = g.insertOrSkipNodeOnMatrix(&NodeOutput{
				NodeInput: &NodeInput{
					Id:   selfLoopId,
					Next: []string{id},
				},
				Anchor: &Anchor{
					Type:   AnchorLoop,
					Margin: AnchorMarginLeft,
					From:   item.Id,
					To:     id,
				},
				IsAnchor:         true,
				RenderIncomes:    []string{node.Id},
				PassedIncomes:    []string{item.Id},
				ChildrenOnMatrix: 0,
			}, state, false)
			if err != nil {
				return
			}
		}
		initialHeight := mtx.Height()
		fromId := fmt.Sprintf("%s-%s-from", id, item.Id)
		toId := fmt.Sprintf("%s-%s-to", id, item.Id)
		node.RenderIncomes = append(node.RenderIncomes, fromId)
		err = g.insertOrSkipNodeOnMatrix(&NodeOutput{
			NodeInput: &NodeInput{
				Id:   toId,
				Next: []string{id},
			},
			Anchor: &Anchor{
				Type:   AnchorLoop,
				Margin: AnchorMarginLeft,
				From:   item.Id,
				To:     id,
			},
			IsAnchor:         true,
			RenderIncomes:    []string{renderIncomeId},
			PassedIncomes:    []string{item.Id},
			ChildrenOnMatrix: 0,
		}, state, true)
		if err != nil {
			return
		}
		if initialHeight != mtx.Height() {
			initialY++
		}
		state.x = income.x
		err = g.insertOrSkipNodeOnMatrix(&NodeOutput{
			NodeInput: &NodeInput{
				Id:   fromId,
				Next: []string{id},
			},
			Anchor: &Anchor{
				Type:   AnchorLoop,
				Margin: AnchorMarginRight,
				From:   item.Id,
				To:     id,
			},
			IsAnchor:         true,
			RenderIncomes:    []string{toId},
			PassedIncomes:    []string{item.Id},
			ChildrenOnMatrix: 0,
		}, state, false)
		if err != nil {
			return
		}
		state.x = initialX
	}
	state.y = initialY
	return
}

func (g *GraphMatrix) insertSplitOutcomes(item *NodeOutput, state *state, levelQueue *TraverseQueue) (err error) {
	queue := state.queue
	outcomes := g.outcomes(item.Id)
	if len(outcomes) == 0 {
		err = fmt.Errorf("split `%s` has no outcomes", item.Id)
		return
	}
	firstOutcomeId := outcomes[0]
	outcomes = append([]string{}, outcomes[1:]...)
	first := g.node(firstOutcomeId)
	queue.add(&item.Id, levelQueue, &NodeInput{
		Id:   first.Id,
		Next: first.Next,
	})
	for _, outcomeId := range outcomes {
		state.y++
		id := fmt.Sprintf("%s-%s", item.Id, outcomeId)

		err = g.insertOrSkipNodeOnMatrix(&NodeOutput{
			NodeInput: &NodeInput{
				Id:   id,
				Next: []string{outcomeId},
			},
			Anchor: &Anchor{
				Type:   AnchorSplit,
				Margin: AnchorMarginRight,
				From:   item.Id,
				To:     outcomeId,
			},
			IsAnchor:         true,
			RenderIncomes:    []string{item.Id},
			PassedIncomes:    []string{item.Id},
			ChildrenOnMatrix: 0,
		}, state, true)
		if err != nil {
			return
		}
		v := g.node(outcomeId)
		queue.add(&id, levelQueue, &v)
	}
	return
}

func (g *GraphMatrix) insertJoinIncomes(item *NodeOutput, state *state, levelQueue *TraverseQueue, addItemToQueue bool) (err error) {
	mtx, queue, incomes := state.mtx, state.queue, item.PassedIncomes
	lowestY, err := g.getLowestYAmongIncomes(item, mtx)
	if err != nil {
		return
	}
	for _, incomeId := range incomes {
		found := mtx.findNode(func(n *NodeOutput) bool {
			return n.Id == incomeId
		})
		if found == nil {
			err = fmt.Errorf("income `%s` is not on matrix yet", incomeId)
			return
		}
		y, income := found.coords[1], found.item
		if lowestY == y {
			item.RenderIncomes = append(item.RenderIncomes, incomeId)
			income.ChildrenOnMatrix = min(income.ChildrenOnMatrix+1, len(income.Next))
			continue
		}
		state.y = y
		id := fmt.Sprintf("%s-%s", incomeId, item.Id)
		item.RenderIncomes = append(item.RenderIncomes, id)
		err = g.insertOrSkipNodeOnMatrix(&NodeOutput{
			NodeInput: &NodeInput{
				Id:   id,
				Next: []string{item.Id},
			},
			Anchor: &Anchor{
				Type:   AnchorJoin,
				Margin: AnchorMarginLeft,
				From:   incomeId,
				To:     item.Id,
			},
			IsAnchor:         true,
			RenderIncomes:    []string{incomeId},
			PassedIncomes:    []string{incomeId},
			ChildrenOnMatrix: 1,
		}, state, false)
		if err != nil {
			return
		}
	}
	if addItemToQueue {
		queue.add(&item.Id, levelQueue, g.getOutcomesArray(item.Id)...)
	}
	return
}

func (g *GraphMatrix) markIncomesAsPassed(mtx *Matrix, item *NodeOutput) (err error) {
	for _, incomeId := range item.RenderIncomes {
		found := mtx.findNode(func(n *NodeOutput) bool {
			return n.Id == incomeId
		})
		if found == nil {
			err = fmt.Errorf("income `%s` is not on matrix yet", incomeId)
			return
		}
		coords := found.coords
		income := found.item
		income.ChildrenOnMatrix = min(income.ChildrenOnMatrix+1, len(income.Next))
		mtx.insert(coords[0], coords[1], income)
	}
	return
}

func (g *GraphMatrix) resolveCurrentJoinIncomes(mtx *Matrix, join *NodeOutput) (err error) {
	err = g.markIncomesAsPassed(mtx, join)
	if err != nil {
		return
	}
	join.RenderIncomes = []string{}
	return
}

func min(numbers ...int) (min int) {
	min = numbers[0]
	for _, n := range numbers {
		if n < min {
			min = n
		}
	}
	return
}
