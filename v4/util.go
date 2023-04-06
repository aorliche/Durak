package main

import (
    "reflect"
)

// Count number meeting condition
func Count[T any](slice []T, fn func (T) bool) int {
    n := 0
    for _,v := range slice {
        if fn(v) {
            n++
        }
    }
    return n
}

// Remove object
func Remove[T any](slice []T, val T) []T {
    idx := IndexOf(slice, val)
    last := len(slice)-1
    if idx != -1 {
        slice[idx] = slice[last]
        slice = slice[:last]
    }
    return slice
}

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

// Transformation
func Apply[S any, T any](slice []S, filter func(S) T) []T {
    ts := make([]T, len(slice))
    for i,s := range slice {
        ts[i] = filter(s)
    }
    return ts
}

// Filter
func Filter[T any](slice []T, include func(T) bool) []T {
    ts := make([]T, 0)
    for _,v := range slice {
        if include(v) {
            ts = append(ts, v)
        }
    }
    return ts
}

// Returns not nil values
func NotNil[T any](slice []*T) []*T {
    return Filter(slice, func(val *T) bool {return val != nil})
}
