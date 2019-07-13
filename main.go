package dgraph

import (
	"github.com/lempiy/dgraph/ascii"
	"github.com/lempiy/dgraph/core"
)

func DrawGraph(list []core.NodeInput) (canvas *ascii.Canvas, err error) {
	g, err := core.NewGraph(list)
	if err != nil {
		return
	}
	mtx, err := g.Traverse()
	if err != nil {
		return
	}
	canvas, err = ascii.DrawAsciiMatrix(mtx)
	if err != nil {
		return
	}
	return
}
