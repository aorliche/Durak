package main

// Data types
// Either List or Props are used
type Object struct {
    Type string
    Props map[string] interface{}
    PropTypes map[string] string
}

type FnSig func([]interface{}) interface{}

type Fn struct {
    F FnSig
    Args []string
    Commutes bool
}

// Args is union of binds of self and children
type Node struct {
    F *Fn
    Children []*Node
    Val interface{}
    Bind int
    Args []int
    SavHash uint32
}

// For serialization
type NodeRec struct {
    Name string         `json:"Name,omitempty"`
    Children []*NodeRec `json:"Children,omitempty"`
    Val string          `json:"Val,omitempty"`
    Bind int
}

// Canonicalize argument order for commuting Fns like Equal
type CanonArg struct {
    Arg interface{}
    Child *Node
}

type Example struct {
    Val bool
    Args []interface{}
}

type Pred struct {
    Nodes []*Node
    Idx int
    Name string
    Args []string
}

type PredOrNode struct {
    n *Node
    p *Pred
}

// For serialization
type PredRec struct {
    Nodes []*NodeRec
    Idx int
    Name string
    Args []string
}

// For satisfiability search
type Index struct {
    NumTerms int
    NumNodes int
    CurIdx int
    Idcs [4]int
}
