package core

import (
	"fmt"
)

func isMultiple(m map[string][]string, id string) bool {
	list, ok := m[id]
	return ok && len(list) > 1
}

type relationMap map[string][]string
type nodesMap map[string]NodeInput

type GraphBasic struct {
	list                []NodeInput
	nodesMap            nodesMap
	incomesByNodeIdMap  relationMap
	outcomesByNodeIdMap relationMap
	loopsByNodeIdMap    relationMap
}

func NewGraphBasic(list []NodeInput) (*GraphBasic, error) {
	g := &GraphBasic{
		list:                list,
		nodesMap:            make(nodesMap),
		incomesByNodeIdMap:  make(relationMap),
		outcomesByNodeIdMap: make(relationMap),
		loopsByNodeIdMap:    make(relationMap),
	}
	for _, node := range g.list {
		if _, ok := g.nodesMap[node.Id]; ok {
			return nil, fmt.Errorf("duplicate node id `%s`", node.Id)
		}
		g.nodesMap[node.Id] = node
	}
	err := g.detectIncomesAndOutcomes()
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *GraphBasic) detectIncomesAndOutcomes() (err error) {
	totalSet := make(Set)
	for _, node := range g.list {
		if totalSet.Has(node.Id) {
			continue
		}
		branchSet := make(Set)
		_, err = g.traverseVertically(node, &branchSet, &totalSet)
		if err != nil {
			return
		}
	}
	return
}

func (g *GraphBasic) traverseVertically(node NodeInput, branchSet *Set, totalSet *Set) (*Set, error) {
	var err error
	if branchSet.Has(node.Id) {
		return nil, fmt.Errorf("duplicate incomes for node id `%s`", node.Id)
	}
	branchSet.Add(node.Id)
	totalSet.Add(node.Id)
	for _, outcomeId := range node.Next {
		// skip loops which are already detected
		if g.isLoopEdge(node.Id, outcomeId) {
			continue
		}
		// detect loops
		if branchSet.Has(outcomeId) {
			addUniqueRelation(&g.loopsByNodeIdMap, node.Id, outcomeId)
			continue
		}
		addUniqueRelation(&g.incomesByNodeIdMap, outcomeId, node.Id)
		addUniqueRelation(&g.outcomesByNodeIdMap, node.Id, outcomeId)
		totalSet, err = g.traverseVertically(g.nodesMap[outcomeId], branchSet.Copy(), totalSet)
		if err != nil {
			return nil, err
		}
	}
	return totalSet, nil
}

func (g *GraphBasic) isLoopEdge(nodeId, outcomeId string) bool {
	loops, ok := g.loopsByNodeIdMap[nodeId]
	if !ok {
		return false
	}
	for _, id := range loops {
		if id == outcomeId {
			return true
		}
	}
	return false
}

func (g *GraphBasic) roots() (roots []NodeInput) {
	for _, node := range g.list {
		if g.isRoot(node.Id) {
			roots = append(roots, node)
		}
	}
	return
}

func (g *GraphBasic) isRoot(id string) bool {
	incomes, ok := g.incomesByNodeIdMap[id]
	return !ok || len(incomes) == 0
}

func (g *GraphBasic) isSplit(id string) bool {
	return isMultiple(g.outcomesByNodeIdMap, id)
}

func (g *GraphBasic) isJoin(id string) bool {
	return isMultiple(g.outcomesByNodeIdMap, id)
}

func (g *GraphBasic) loops(id string) []string {
	return g.loopsByNodeIdMap[id]
}

func (g *GraphBasic) outcomes(id string) []string {
	return g.outcomesByNodeIdMap[id]
}

func (g *GraphBasic) incomes(id string) []string {
	return g.incomesByNodeIdMap[id]
}

func (g *GraphBasic) node(id string) NodeInput {
	return g.nodesMap[id]
}

func (g *GraphBasic) nodeType(id string) NodeType {
	switch {
	case g.isRoot(id) && g.isSplit(id):
		return NodeTypeRootSplit
	case g.isRoot(id):
		return NodeTypeRootSimple
	case g.isSplit(id) && g.isJoin(id):
		return NodeTypeSplitJoin
	case g.isSplit(id):
		return NodeTypeSplit
	case g.isJoin(id):
		return NodeTypeJoin
	default:
		return NodeTypeSimple
	}
}

func (g *GraphBasic) getOutcomesArray(itemId string) (result []*NodeInput) {
	outcomes := g.outcomes(itemId)
	if len(outcomes) == 0 {
		return
	}
	for _, id := range outcomes {
		node := g.node(id)
		result = append(result, &node)
	}
	return
}

func addUniqueRelation(rm *relationMap, key, value string) {
	relations, ok := (*rm)[key]
	if !ok {
		(*rm)[key] = []string{value}
	}
	isUnique := true
	for _, v := range relations {
		if v == value {
			isUnique = false
			break
		}
	}
	if isUnique {
		(*rm)[key] = append(relations, value)
	}
}
