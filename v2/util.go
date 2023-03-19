package main

import (
    "reflect"
    "runtime"
    "sort"
    //"unsafe"
)

// https://stackoverflow.com/questions/7052693/how-to-get-the-name-of-a-function-in-go
func GetFunctionName(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// Convenience
func Ternary[T any](cond bool, a T, b T) T {
    if cond {
        return a
    }
    return b
}

// Unique elements of array
func Unique[T any](t []T) []T {
    uniq := make([]T, 0)
    outer:
    for i:=0; i<len(t); i++ {
        for j:=0; j<len(uniq); j++ {
            if reflect.DeepEqual(t[i], uniq[j]) {
                continue outer
            }
        }
        uniq = append(uniq, t[i])
    }
    return uniq
}

// My own integer power function
func IntPow(base int, exp int) int {
    if exp == 0 {
        return 1
    }
    res := 1
    for i := 0; i<exp; i++ {
        res *= base
    }
    return res
}

// Generic Index of
func IndexOf[T any](slice []T, val T) int {
    for idx, v := range slice {
        if (reflect.DeepEqual(v, val)) {
            return idx
        }
    }
    return -1
}

// Keys from map
func GetKeys(props map[string]interface{}) []string {
    keys := make([]string, len(props))
    i := 0
    for k := range props {
        keys[i] = k
        i++
    }
    sort.Strings(keys)
    return keys
}

// Concatenate multiple slices into one
func Cat[T any](slices ...[]T) []T {
    res := make([]T, 0)
    for _,sl := range slices {
        res = append(res, sl...)
    }
    return res
}
