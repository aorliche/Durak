package main

// Data types
type object struct {
    typ string
    props map[string] interface{}
    propTypes map[string] string
}

type fnSig func([]interface{}) interface{}

type fn struct {
    f fnSig
    args []string
    commutes bool
}

type node struct {
    f *fn
    children []*node
    val interface{}
    bind int
    args []int
    vhash uint32
    fhash uint32
}

// Canonicalize argument order for commuting fns like equal
type canonArg struct {
    arg interface{}
    child *node
}

type example struct {
    val bool
    args []interface{}
}

type history struct {
    nodes []*node
    idx int
    exs []*example
}

// Number of nodes and idx determines how predicate is evaluated
// evalCombo?() methods in sat.go
// e.g. A~B, AB+~C, etc
type pred struct {
    nodes []*node
    idx int
    name string
    argTypes []string
    exs []*example
    hist []*history
    fns []*fn
}

type index struct {
    nTerms int
    nNodes int
    curIdx int
    idcs [4]int
}
