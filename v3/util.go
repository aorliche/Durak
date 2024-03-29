package main

import (
    "reflect"
)

// Generic IndexOf
func IndexOf[T any](slice []T, val T) int {
    for idx, v := range slice {
        if reflect.DeepEqual(v, val) {
            return idx
        }
    }
    return -1
}

// Using function
func IndexOfFn[T any](slice []T, fn func(T) bool) int {
    for idx, v := range slice {
        if fn(v) {
            return idx
        }
    }
    return -1
}

// Concatenate multiple slices into one
func Cat[T any](slices ...[]T) []T {
    res := make([]T, 0)
    for _,sl := range slices {
        res = append(res, sl...)
    }
    return res
}

// Convenience
func Ternary[T any](cond bool, a T, b T) T {
    if cond {
        return a
    }
    return b
}

// Filter
func Apply[S any, T any](slice []S, filter func(S) T) []T {
    ts := make([]T, len(slice))
    for i,s := range slice {
        ts[i] = filter(s)
    }
    return ts
}
