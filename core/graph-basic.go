package core

import (
	"fmt"
)

func isMultiple(m map[string][]string, id string) bool {
	list, ok := m[id]
	return ok && len(list) > 1
}

type Set map[string]struct{}

func (s *Set) Add(key string) {
	(*s)[key] = struct{}{}
}

func (s *Set) Has(key string) bool {
	_, ok := (*s)[key]
	return ok
}

func (s *Set) Copy() *Set {
	copy := make(Set)
	m := *s
	for v := range m {
		copy[v] = struct{}{}
	}
	return &copy
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
		if g.IsLoopEdge(node.Id, outcomeId) {
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

func (g *GraphBasic) IsLoopEdge(nodeId string, outcomeId string) bool {
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
