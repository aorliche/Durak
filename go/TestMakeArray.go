package main

import (
    "fmt"
)

func main() {
    a := make([]int, 2)
    a[1] = 3
    fmt.Println(a[1])
    b := make(map[[2]int]bool)
    b[[2]int{0,1}] = true
    for k,v := range b {
        fmt.Println(k,v)
    }
}
