package main

import (
    "fmt"
    "reflect"
)

func main() {
    a := make([]int, 0)
    b := make([]string, 0)
    fmt.Println(reflect.TypeOf(a).Kind() == reflect.TypeOf(b).Kind())
    fmt.Println(reflect.TypeOf(a).Kind() == reflect.Slice)
    fmt.Println(reflect.TypeOf(a) == reflect.TypeOf(a))
}
