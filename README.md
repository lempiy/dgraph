## ↬ dGraph ↫

#### Draw direct graphs with ascii symbols using Golang and more...

##### Usage

```go
package main

import (
	"encoding/json"
	"log"
	"fmt"
	"github.com/lempiy/dgraph"
	"github.com/lempiy/dgraph/core"
)

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
		log.Fatal(err)
		return
	}
	canvas, err := dgraph.DrawGraph(list)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("%s\n", canvas)
}
```
**Output**:
```
┌─────────┐   ┌─────────┐       ┌─────┐       ┌─────────────┐   ┌─────┐   ┌─────────────┐                     ┌───────────┐
│ Client  ├───┤ Route53 ├───────┤ ELB ├───────┤ WebServer1  ├───┤ LB1 ├───┤ AppServer1  ├─────────────────────┤ DBMaster  │
└─────────┘   └────┬────┘       └──┬──┘       └─────────────┘   └──┬──┘   └──────┬──────┘                     └─────┬─────┘
                   │               │                               │             │                                  │      
                   │               │                               │             │                                  │      
                   │               │                               │             │                                  │      
                   │               │          ┌─────────────┐      │             │          ┌─────────────┐         │      
                   │               ├──────────┤ WebServer2  ├──────┤             └──────────┤ DBReplica1  ├─────────┤      
                   │               │          └─────────────┘      │                        └──────┬──────┘         │      
                   │               │                               │                               │                │      
                   │               │                               │                               │                │      
                   │               │                               │                               │                │      
                   │               │                               │      ┌─────────────┐          │                │      
                   │               │                               └──────┤ AppServer2  ├──────────┼────────────────┤      
                   │               │                                      └──────┬──────┘          │                │      
                   │               │                                             │                 │                │      
                   │               │                                             │                 │                │      
                   │               │                                             │                 │                │      
                   │               │                                             │                 │                │      
                   │               │                                             └─────────────────┘                │      
                   │               │                                                                                │      
                   │               │                                                                                │      
                   │               │                                                                                │      
                   │               │                                                                                │      
                   │               │          ┌─────────────┐   ┌─────┐   ┌─────────────┐                           │      
                   │               ├──────────┤ WebServer3  ├───┤ LB2 ├───┤ AppServer3  ├───────────────────────────┤      
                   │               │          └─────────────┘   └──┬──┘   └──────┬──────┘                           │      
                   │               │                               │             │                                  │      
                   │               │                               │             │                                  │      
                   │               │                               │             │                                  │      
                   │               │          ┌─────────────┐      │             │          ┌─────────────┐         │      
                   │               └──────────┤ WebServer4  ├──────┤             └──────────┤ DBReplica2  ├─────────┤      
                   │                          └─────────────┘      │                        └──────┬──────┘         │      
                   │                                               │                               │                │      
                   │                                               │                               │                │      
                   │                                               │                               │                │      
                   │        ┌─────────────┐       ┌─────┐          │      ┌─────────────┐          │                │      
                   └────────┤ CloudFront  ├───────┤ S3  │          └──────┤ AppServer4  ├──────────┼────────────────┘      
                            └─────────────┘       └─────┘                 └──────┬──────┘          │                       
                                                                                 │                 │                       
                                                                                 │                 │                       
                                                                                 │                 │                       
                                                                                 │                 │                       
                                                                                 └─────────────────┘                       
                                                                                                                           

```
If you need to render your graph in other format. You can use `core` package to get low level `core.Matrix` struct.
