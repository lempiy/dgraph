package core

type NodeType int

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeRootSimple
	NodeTypeRootSplit
	NodeTypeSimple
	NodeTypeSplit
	NodeTypeJoin
	NodeTypeSplitJoin
)

type AnchorType int

const (
	AnchorUnknown AnchorType = iota
	AnchorJoin
	AnchorSplit
	AnchorLoop
)

type AnchorMarginType int

const (
	AnchorMarginNone AnchorMarginType = iota
	AnchorMarginLeft
	AnchorMarginRight
)

type NodeInput struct {
	Id   string
	Next []string
}

type NodeOutput struct {
	*NodeInput
	*Anchor
	IsAnchor         bool
	PassedIncomes    []string
	RenderIncomes    []string
	ChildrenOnMatrix int
}

type Anchor struct {
	Type   AnchorType
	From   string
	To     string
	Margin AnchorMarginType
}

type MatrixNode struct {
	*NodeOutput
	X int
	Y int
}

type loopNode struct {
	id         string
	node       *NodeOutput
	x          int
	y          int
	isSelfLoop bool
}
