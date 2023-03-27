package main

import (
    "fmt"
    "math/rand"
    "reflect"
    "gonum.org/v1/gonum/mat"
)

type AddWeights struct {
    A1 float64
    A2 float64
}

type MultWeights struct {
    A1 float64
}

func (w *AddWeights) ToVec() *mat.VecDense {
    ww := reflect.ValueOf(*w)
    return ToVec(&ww)
}

func (w *MultWeights) ToVec() *mat.VecDense {
    ww := reflect.ValueOf(*w)
    return ToVec(&ww)
}

func ToVec(st *reflect.Value) *mat.VecDense {
    v := make([]float64, st.NumField())
    for i := 0; i < st.NumField(); i++ {
        v[i] = st.Field(i).Interface().(float64)
    }
    vec := mat.NewVecDense(st.NumField(), v)
    return vec
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

func Eval(a *AddWeights, m *MultWeights, aData []*AddWeights, mData []*MultWeights) []float64 {
    av := a.ToVec()
    mv := m.ToVec()
    res := make([]float64, len(aData))
    for i := 0; i < len(aData); i++ {
        adv := aData[i].ToVec()
        adv.MulElemVec(av, adv)
        mdv := mData[i].ToVec()
        mdv.MulElemVec(mv, mdv)
        mdv.ScaleVec(mat.Sum(adv), mdv)
        res[i] = mat.Sum(mdv)
    }
    return res
}

func UpdateWeights(a *AddWeights, m *MultWeights, aData []*AddWeights, mData []*MultWeights, y []float64, lr float64) {
    yhat := Eval(a, m, aData, mData)
    av := a.ToVec()
    mv := m.ToVec()
    dyda := mat.NewVecDense(av.Len(), nil)
    dyda.ScaleVec(0, dyda)
    dydm := mat.NewVecDense(mv.Len(), nil)
    dydm.ScaleVec(0, dydm)
    n := len(aData)
    for i := 0; i < n; i++ {
        adv := aData[i].ToVec()
        wa := mat.NewVecDense(adv.Len(), nil)
        wa.MulElemVec(av, adv)
        mdv := mData[i].ToVec()
        wm := mat.NewVecDense(mdv.Len(), nil)
        wm.MulElemVec(mv, mdv)
        // dy/dm
        mdv.ScaleVec((y[i]-yhat[i])*lr/float64(n)*mat.Sum(wa), mdv)
        dydm.AddVec(mdv, dydm)
        // dy/dw
        adv.ScaleVec((y[i]-yhat[i])*lr/float64(n)*mat.Sum(wm), adv)
        dyda.AddVec(adv, dyda)
    }
    // Update weights
    fmt.Println(mat.Formatted(dyda))
    fmt.Println(mat.Formatted(dydm))
}

func main() {
    a := AddWeights{2, 4}
    m := MultWeights{2}
    ap := &a
    mp := &m
    av := ap.ToVec()
    mv := mp.ToVec()
    fmt.Println(mat.Formatted(av))
    fmt.Println(mat.Formatted(mv))
    addData, multData := GenData(100)
    y := Eval(&a, &m, addData, multData)
    fmt.Println(y)
    fmt.Println(mat.Formatted(addData[0].ToVec()))
    fmt.Println(mat.Formatted(multData[0].ToVec()))
    UpdateWeights(&a, &m, addData, multData, y, 0.01)
}

// Inputs: a, d
// Weights: m, b, w
// Eq: y = Sum((m*d+b)*w*a)
/*func UpdateWeights(y float64, m *MultWeights, b *MultWeights, w *Weights, a *Weights, d *MultWeights, lr float64) {
    mm := reflect.ValueOf(m)
    bb := reflect.ValueOf(b)
    ww := reflect.ValueOf(w)
    aa := reflect.ValueOf(a)
    dd := reflect.ValueOf(d)
    yhat := 0
    for i := 0; i < mm.NumField(); i++ {
        mv := mm.Field(i).Interface().(float64)
        bv := bb.Field(i).Interface().(float64)
        dv := dd.Field(i).Interface().(float64)
        mult := mv*bv+dv
        for j := 0; j < ww.NumField(); j++ {
            wv := ww.Field(j).Interface().(float64)
            av := aa.Field(j).Interface().(float64)
            val := wv*av
            yhat += mult*val
        }
    }
    err := y-yhat
    // Get partials
    for i := 0; i < mm.NumField(); i++ {
        mv := mm.Field(i).Interface().(float64)
        bv := bb.Field(i).Interface().(float64)
        dv := dd.Field(i).Interface().(float64)
        mult := mv*bv+dv
        for j := 0; j < ww.NumField(); j++ {
            wv := ww.Field(j).Interface().(float64)
            av := aa.Field(j).Interface().(float64)
            val := wv*av
            yhat += mult*val
        }
    }
}*/
