package main

import (
    "fmt"
    "hash/crc32"
    "encoding/binary"
    "unsafe"
)

func main() {
    b := make([]byte, 0)
    c := crc32.NewIEEE()
    a := 1
    aa := &a
    aaa := uint64(uintptr(unsafe.Pointer(aa)))
    b = binary.LittleEndian.AppendUint64(b, aaa)
    c.Write(b)
    fmt.Println(b)
}
