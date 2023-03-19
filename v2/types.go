package main

// Data types
// Type means bool, string, int, Slice, Pair, Prod
type Object struct {
    Type string
    Val interface{}
    Props map[string]*Object
}

type Atom func([]*Object) *Object

// Args: types of args
type Fn struct {
    A Atom
    N *Node
    F *Filter
    Args []string
    Commutes bool
}

type Filter struct {
    Pred *Fn
}

// BindTree is ordered union of binds of self and children
// Type: return type, filled in during search
type Node struct {
    F *Fn
    Children []*Node
    Val *Object
    Type string
    Bind int
    BindTree []int
    SavHash uint32
}
