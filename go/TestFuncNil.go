package main

import "fmt"

type fn func() int

func a() int {
    return 1
}

func main() {
    var f fn
    f = nil
    fmt.Println(f)
}
