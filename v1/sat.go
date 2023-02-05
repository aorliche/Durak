package main

func satisfy(int nargs, exs []*example, table [][]*nodeArgs) *disj {
        
}

func argComb(int nargs, row []*nodeArgs) [][]int {
    sets := make([][]int, nargs)
    for i:=0; i<nargs; i++ {
        sets[i] = make([]int, 0)
    }
    for i,n := range row {
        for _,j := range n.args {
            sets[j] = append(sets[j], i)
        }
    }
    num := 1
    for i:=0; i<nargs; i++ {
        num *= len(sets[i])
    }
    res := make([][]int, num)
    for i:=0; i<num; i++ {
        mod := 1
        res[i] = make([]int, nargs) 
        for j := 0; j < nargs; j++ {
            idx = sets[j][(i/mod)%len(sets[j])]
            res[i] = append(res[i], idx)
            mod *= len(cargs[j])
        }
    }
    return res
}
