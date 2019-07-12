package main

import (
	"encoding/json"
	"fmt"
	"github.com/lempiy/dgraph/ascii"
	"github.com/lempiy/dgraph/core"
)

//var list = []core.NodeInput{
//	{
//		Id:   "A",
//		Next: []string{"B"},
//	},
//	{
//		Id:   "B",
//		Next: []string{"AWS", "D", "S"},
//	},
//	{
//		Id:   "D",
//		Next: []string{"E"},
//	},
//	{
//		Id:   "AWS",
//		Next: []string{"F"},
//	},
//	{
//		Id:   "S",
//		Next: []string{"DEL"},
//	},
//	{
//		Id:   "E",
//		Next: []string{"DEL"},
//	},
//	{
//		Id:   "F",
//		Next: []string{"DEL"},
//	},
//	{
//		Id:   "DEL",
//		Next: []string{},
//	},
//}

//const data = `[
//   {
//       "id": "1",
//       "next": ["2"]
//   },
//   {
//       "id": "2",
//       "next": ["3", "4"]
//   },
//   {
//       "id": "3",
//       "next": ["5"]
//   },
//   {
//       "id": "6",
//       "next": []
//   },
//   {
//       "id": "7",
//       "next": ["4"]
//   },
//   {
//       "id": "4",
//       "next": ["8"]
//   },
//   {
//       "id": "8",
//       "next": ["9"]
//   },
//   {
//       "id": "10",
//       "next": ["11"]
//   },
//   {
//       "id": "11",
//       "next": ["12", "13"]
//   },
//   {
//       "id": "12",
//       "next": ["14"]
//   },
//   {
//       "id": "14",
//       "next": ["15", "16"]
//   },
//   {
//       "id": "15",
//       "next": ["17", "12"]
//   },
//   {
//       "id": "16",
//       "next": ["17"]
//   },
//   {
//       "id": "17",
//       "next": []
//   },
//   {
//       "id": "13",
//       "next": []
//   },
//   {
//       "id": "9",
//       "next": ["18"]
//   },
//   {
//       "id": "5",
//       "next": ["19"]
//   },
//   {
//       "id": "19",
//       "next": ["20"]
//   },
//   {
//       "id": "18",
//       "next": ["20"]
//   },
//   {
//       "id": "20",
//       "next": ["21", "6"]
//   },
//   {
//       "id": "21",
//       "next": ["21"]
//   }
//]`

const data = `[
{
	"id": "Client",
    "next": ["Route53"]
},
{
	"id": "Route53",
    "next": ["ELB", "CloudFront"]
},
{
	"id": "CloudFront",
    "next": ["S3"]
},
{
	"id": "S3",
    "next": []
},
{
	"id": "ELB",
    "next": ["WebServer1", "WebServer2", "WebServer3", "WebServer4"]
},
{
	"id": "WebServer1",
    "next": ["LB1"]
},
{
	"id": "WebServer2",
    "next": ["LB1"]
},
{
	"id": "WebServer3",
    "next": ["LB2"]
},
{
	"id": "WebServer4",
    "next": ["LB2"]
},
{
	"id": "LB1",
    "next": ["AppServer1", "AppServer2"]
},
{
	"id": "LB2",
    "next": ["AppServer3", "AppServer4"]
},
{
	"id": "AppServer1",
    "next": ["DBMaster", "DBReplica1"]
},
{
	"id": "AppServer2",
    "next": ["DBMaster", "DBReplica1"]
},
{
	"id": "AppServer3",
    "next": ["DBMaster", "DBReplica2"]
},
{
	"id": "AppServer4",
    "next": ["DBMaster", "DBReplica2"]
},
{
	"id": "DBReplica1",
    "next": ["DBMaster"]
},
{
	"id": "DBReplica2",
    "next": ["DBMaster"]
},
{
	"id": "DBMaster",
    "next": []
}
]`

func main() {
	var list []core.NodeInput
	err := json.Unmarshal([]byte(data), &list)
	if err != nil {
		fmt.Println(err)
		return
	}
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
	c, err := ascii.DrawAsciiMatrix(mtx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", c)
}
