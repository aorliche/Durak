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
}

type example struct {
    val bool
    args []interface{}
}

type history struct {
    nodes []*node
    idx int
}

// nNodes and idx determine how predicate is evaluated
// e.g. A~B, AB+~C, etc
type pred struct {
    nodes []*node
    idx int
    name string
    argTypes []string
    exs []*example
    hist []*history
}

