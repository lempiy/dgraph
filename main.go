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

//const data = `[
//    {
//        "id": "1",
//        "next": ["2"]
//    },
//    {
//        "id": "2",
//        "next": ["3", "4"]
//    },
//    {
//        "id": "3",
//        "next": ["5"]
//    },
//    {
//        "id": "6",
//        "next": []
//    },
//    {
//        "id": "7",
//        "next": ["4"]
//    },
//    {
//        "id": "4",
//        "next": ["8"]
//    },
//    {
//        "id": "8",
//        "next": ["9"]
//    },
//    {
//        "id": "10",
//        "next": ["11"]
//    },
//    {
//        "id": "11",
//        "next": ["12", "13"]
//    },
//    {
//        "id": "12",
//        "next": ["14"]
//    },
//    {
//        "id": "14",
//        "next": ["15", "16"]
//    },
//    {
//        "id": "15",
//        "next": ["17", "12"]
//    },
//    {
//        "id": "16",
//        "next": ["17"]
//    },
//    {
//        "id": "17",
//        "next": []
//    },
//    {
//        "id": "13",
//        "next": []
//    },
//    {
//        "id": "9",
//        "next": ["18"]
//    },
//    {
//        "id": "5",
//        "next": ["19"]
//    },
//    {
//        "id": "19",
//        "next": ["20"]
//    },
//    {
//        "id": "18",
//        "next": ["20"]
//    },
//    {
//        "id": "20",
//        "next": ["21", "6"]
//    },
//    {
//        "id": "21",
//        "next": ["21"]
//    }
//]`

func main() {
	//var list []core.NodeInput
	//err := json.Unmarshal([]byte(data), &list)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
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
	fmt.Printf("%s\n", mtx)
}
