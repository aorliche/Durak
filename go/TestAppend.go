package main

import (
    "fmt"
)

func main() {
    a := make([]string, 0)
    a = append(a, "hello")
    b := append(a, "world")
    fmt.Println(a)
    fmt.Println(b)
}
