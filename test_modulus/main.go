package main

import (
    "fmt"
)

func main() {
    a := 0
    b := 1
    c := 2
    ab := (a-b)%3
    ac := (a-c)%3
    fmt.Println(ab)
    fmt.Println(ac)
}
