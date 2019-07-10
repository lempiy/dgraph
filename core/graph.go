package core

import "fmt"

const MaxIterations = 100000

type Graph struct {
	*GraphMatrix
}

func NewGraph(list []NodeInput) (*Graph, error) {
	g, err := NewGraphMatrix(list)
	if err != nil {
		return nil, err
	}
	return &Graph{
		GraphMatrix: g,
	}, nil
}

func (g *Graph) handleSplitNode(item *NodeOutput, state *state, levelQueue *TraverseQueue) (err error) {
	isInserted, err := g.processOrSkipNodeOnMatrix(item, state)
	if err != nil {
		return
	}
	if isInserted {
		err = g.insertSplitOutcomes(item, state, levelQueue)
		if err != nil {
			return
		}
	}
	return
}

func (g *Graph) handleSplitJoinNode(item *NodeOutput, state *state, levelQueue *TraverseQueue) (err error) {
	queue, mtx := state.queue, state.mtx
	if g.joinHasUnresolvedIncomes(item) {
		queue.push(item)
		return
	}
	err = g.resolveCurrentJoinIncomes(mtx, item)
	if err != nil {
		return
	}
	isInserted, err := g.processOrSkipNodeOnMatrix(item, state)
	if err != nil {
		return
	}
	if isInserted {
		err = g.insertJoinIncomes(item, state, levelQueue, false)
		if err != nil {
			return
		}
		err = g.insertSplitOutcomes(item, state, levelQueue)
		if err != nil {
			return
		}
	}
	return
}

func (g *Graph) handleJoinNode(item *NodeOutput, state *state, levelQueue *TraverseQueue) (err error) {
	queue, mtx := state.queue, state.mtx
	if g.joinHasUnresolvedIncomes(item) {
		queue.push(item)
		return
	}
	err = g.resolveCurrentJoinIncomes(mtx, item)
	if err != nil {
		return
	}
	isInserted, err := g.processOrSkipNodeOnMatrix(item, state)
	if err != nil {
		return
	}
	if isInserted {
		err = g.insertJoinIncomes(item, state, levelQueue, true)
	}
	return
}

func (g *Graph) handleSimpleNode(item *NodeOutput, state *state, levelQueue *TraverseQueue) (err error) {
	queue := state.queue
	isInserted, err := g.processOrSkipNodeOnMatrix(item, state)
	if err != nil {
		return
	}
	if isInserted {
		queue.add(&item.Id, levelQueue, g.getOutcomesArray(item.Id)...)
	}
	return
}

func (g *Graph) traverseItem(item *NodeOutput, state *state, levelQueue *TraverseQueue) (err error) {
	mtx := state.mtx
	switch g.nodeType(item.Id) {
	case NodeTypeRootSimple:
		state.y = mtx.getFreeRowForColumn(0)
		fallthrough
	case NodeTypeSimple:
		err = g.handleSimpleNode(item, state, levelQueue)
	case NodeTypeRootSplit:
		state.y = mtx.getFreeRowForColumn(0)
		fallthrough
	case NodeTypeSplit:
		err = g.handleSplitNode(item, state, levelQueue)
	case NodeTypeJoin:
		err = g.handleJoinNode(item, state, levelQueue)
	case NodeTypeSplitJoin:
		err = g.handleSplitJoinNode(item, state, levelQueue)
	default:
		err = fmt.Errorf("unknown node type - %d", g.nodeType(item.Id))
	}
	return
}
