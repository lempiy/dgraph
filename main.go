package main

import (
	"fmt"
	"github.com/lempiy/dgraph/core"
)

var list = []core.NodeInput{
	{
		Id:   "A",
		Next: []string{"B"},
	},
	{
		Id:   "B",
		Next: []string{"C", "D"},
	},
	{
		Id:   "D",
		Next: []string{"E"},
	},
	{
		Id:   "C",
		Next: []string{"F"},
	},
	{
		Id:   "E",
		Next: []string{"J"},
	},
	{
		Id:   "F",
		Next: []string{"J"},
	},
	{
		Id:   "J",
		Next: []string{},
	},
}

func main() {
	g, err := core.NewGraph(list)
	if err != nil {
		fmt.Println(err)
		return
	}
	mtx, err := g.Traverse()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", mtx.Normalize())
}
