package main

import (
    "fmt"
    "math/rand"
    "reflect"
    . "gorgonia.org/gorgonia"
    T "gorgonia.org/tensor"
)

type AddWeights struct {
    A1 float64
    A2 float64
}

type MultWeights struct {
    A1 float64
}

type NN struct {
    g *ExprGraph
    A, M *Node
    WA, WM *Node
    Res *Node
    ResVal Value
}

func InitNN(g *ExprGraph, nA int, nM int) *NN {
    a := NewVector(g, Float64, WithShape(nA), WithName("A"))
    m := NewVector(g, Float64, WithShape(nM), WithName("M"))
    wA := NewVector(g, Float64, WithShape(nA), WithName("WA"), WithInit(GlorotN(1.0)))
    wM := NewVector(g, Float64, WithShape(nM), WithName("WM"), WithInit(GlorotN(1.0)))
    return &NN{g: g, A: a, M: m, WA: wA, WM: wM}
}

func (nn *NN) Eval() (err error) {
    aa := Must(Mul(nn.WA, nn.A))
    mm := Must(Mul(nn.WM, nn.M))
    b := Must(Sum(mm))
    c := Must(Mul(aa, b))
    nn.Res = Must(Sum(c))
    Read(nn.Res, &nn.ResVal)
    return nil
}

func (nn *NN) Learnables() Nodes {
    return Nodes{nn.WA, nn.WM}
}

func ToNode(g *ExprGraph, vec T.Tensor, name string) *Node {
    node := NewVector(g, T.Float64, WithName(name), WithShape(vec.Shape()...), WithValue(vec))
    return node
}

func (w *AddWeights) ToVec() T.Tensor {
    ww := reflect.ValueOf(*w)
    return ToVec(&ww)
}

func (w *MultWeights) ToVec() T.Tensor {
    ww := reflect.ValueOf(*w)
    return ToVec(&ww)
}

func ToVec(st *reflect.Value) T.Tensor {
    v := make([]float64, st.NumField())
    for i := 0; i < st.NumField(); i++ {
        v[i] = st.Field(i).Interface().(float64)
    }
    return T.New(T.WithShape(st.NumField()), T.WithBacking(v))
}

func GenData(n int) ([]*AddWeights, []*MultWeights) {
    a := make([]*AddWeights, n)
    m := make([]*MultWeights, n)
    for i := 0; i < n; i++ {
        aa := AddWeights{rand.NormFloat64(), rand.NormFloat64()}
        mm := MultWeights{rand.NormFloat64()}
        a[i] = &aa
        m[i] = &mm
    }
    return a, m
}

func main() {
    a := AddWeights{2, 4}
    m := MultWeights{2}
    ap := &a
    mp := &m
    av := ap.ToVec()
    mv := mp.ToVec()
    fmt.Println(av)
    fmt.Println(mv)
    aData, mData := GenData(100)
    fmt.Println(aData[0].ToVec())
    fmt.Println(mData[0].ToVec())
    // Make graph
    g := NewGraph()
    nn := InitNN(g, 2, 1)
    nn.Eval()
    Let(nn.WA, av)
    Let(nn.WM, mv)
    // Evaluate ground truth
    y := make([]float64, len(aData))
    vm := NewTapeMachine(g)
    defer vm.Close()
    for i := 0; i < len(aData); i++ {
        vm.Reset()
        Let(nn.A, aData[i].ToVec())
        Let(nn.M, mData[i].ToVec())
        vm.RunAll()
        y[i] = nn.ResVal.Data().(float64)
        fmt.Println(y[i])
    }
    fmt.Println("---")
    // Learn
    h := NewGraph()
    nn = InitNN(h, 2, 1)
    nn.Eval()
    yv := NewScalar(h, T.Float64, WithName("y"), WithShape(1))
    losses := Must(Sub(yv, nn.Res))
    sq := Must(Square(losses))
    cost := Must(Mean(sq))
    Grad(cost, nn.Learnables()...)
    vm = NewTapeMachine(h)
    defer vm.Close()
    solver := NewVanillaSolver(WithLearnRate(0.01))
    for i := 0; i < len(aData); i++ {
        Let(nn.A, aData[i].ToVec())
        Let(nn.M, mData[i].ToVec())
        Let(yv, y[i])
        if err := vm.RunAll(); err != nil {
            fmt.Println(err)
        }
        solver.Step(NodesToValueGrads(nn.Learnables()))
        vm.Reset()
        y[i] = nn.ResVal.Data().(float64)
        fmt.Println(y[i])
    }
}
