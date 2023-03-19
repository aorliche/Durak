package main

// Data types
// Type means bool, string, int, slice, pair, prod, Card, Player
// Object types capitalized, "base" types lowercase
type Object struct {
    Type string
    Val interface{}
    Props map[string]*Object
}

type Atom func([]*Object) *Object

// Args: types of args
type Fn struct {
    Name string
    A Atom
    N *Node
    F *Filter
    Args []string
    Commutes bool
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

type Filter struct {
    N *Node
    C []*Filter
    D []*Filter
    Out [][]*Object
    NumObjects int // for ranking in final solve step
}
