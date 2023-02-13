package main

import "fmt"

func A(a int) {
    fmt.Println(a)
}

func A(a int, b int) {
    fmt.Println(a,b)
}

func main() {
    A(1)
    A(2,3)
}
