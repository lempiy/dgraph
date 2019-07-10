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
	fmt.Println("before")
	isInserted, err := g.processOrSkipNodeOnMatrix(item, state)
	if err != nil {
		return
	}
	fmt.Println("after")
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

func (g *Graph) traverseLevel(iterations int, state *state) (n int, err error) {
	queue := state.queue
	levelQueue := queue.drain()
	for levelQueue.length() != 0 {
		iterations++
		item := levelQueue.shift()
		err = g.traverseItem(item, state, levelQueue)
		if err != nil {
			return
		}
		if iterations > MaxIterations {
			err = fmt.Errorf("infinite loop")
			return
		}
	}
	n = iterations
	return
}

func (g *Graph) traverseList(state *state) (mtx *Matrix, err error) {
	safe := 0
	mtx, queue := state.mtx, state.queue
	for queue.length() != 0 {
		safe, err = g.traverseLevel(safe, state)
		if err != nil {
			return
		}
		state.x++
	}
	return
}

func (g *Graph) Traverse() (mtx *Matrix, err error) {
	roots := g.roots()
	state := &state{
		mtx:   NewMatrix(),
		queue: NewTraverseQueue(),
	}
	if len(roots) == 0 {
		err = fmt.Errorf("no graph roots found")
		return
	}
	mtx, queue := state.mtx, state.queue
	queue.add(
		nil, nil, roots...)
	mtx, err = g.traverseList(state)
	return
}
