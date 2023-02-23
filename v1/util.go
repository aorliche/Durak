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
func Unique[T comparable](t []T) []T {
    uniq := make([]T, 0)
    outer:
    for i:=0; i<len(t); i++ {
        for j:=0; j<len(uniq); j++ {
            if t[i] == uniq[j] {
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

// Array Minus
// Assumes sorted keys
/*func StrArrMinus(keys1 []string, keys2 []string) []string {
    set := make([]string, 0)
    for i,j := 0,0; i < len(keys1); i++ {
        for j < len(keys2) && keys2[j] < keys1[i] {
            j++
        }
        if j < len(keys2) && keys1[i] == keys2[j] {
            j++
            continue
        }
        set = append(set, keys1[i])
    }
    return set
}

// Array Intersection
// Assumes sorted keys
func StrArrInt(keys1 []string, keys2 []string) []string {
    set := make([]string, 0)
    for i,j := 0,0; i < len(keys1); i++ {
        for j < len(keys2) && keys2[j] < keys1[i] {
            j++
        }
        if j == len(keys2) {
            break
        }
        if keys1[i] == keys2[j] {
            set = append(set, keys1[i])
            j++
        }
    }
    return set
}*/

// Hashing functions
/*func AppendPtr[T any](b []byte, ptr *T) []byte {
    addr := uint64(uintptr(unsafe.Pointer(ptr)))
    return binary.LittleEndian.AppendUint64(b, addr)
}

func AppendBool(b []byte, flag bool) []byte {
    ui := uint32(0)
    if flag {
        ui = 1
    }
    return binary.LittleEndian.AppendUint32(b, ui)
}

func AppendAny(b []byte, val interface{}) []byte {
    switch val.(type) {
        case bool: return AppendBool(b, val.(bool))
        case int: return binary.LittleEndian.AppendUint32(b, uint32(val.(int)))
        case string: return append(b, []byte(val.(string))...)
        case *object: return AppendPtr(b, val.(*object))
        case []*object: return AppendPtr(b, &val.([]*object)[0])
    }
    return b
}

func Hash(n *node) uint32 {
    if n.vhash != 0 {
        return n.vhash
    }
    c := crc32.NewIEEE()
    b := make([]byte, 0)
    if n.f != nil {
        b = AppendPtr(b, n.f)
    }
    b = AppendAny(b, n.val)
    for _,m := range n.children {
        b = AppendAny(b, m.val)
    }
    c.Write(b)
    n.vhash = c.Sum32()
    return n.vhash
}*/
