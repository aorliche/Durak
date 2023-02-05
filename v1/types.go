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
}

// 0: game
type nodeArgs struct {
    n *node
    args []int
}

type conj struct {
    terms []*node
    neg []bool
}

type disj struct {
    terms []*conj
}

type example struct {
    val bool
    args []interface{}
}

type history struct {
    dis *disj
    ex *example
}

type pred struct {
    dis *disj
    name string
    argTypes []string
    hist []*history
    fns []*fn
}

