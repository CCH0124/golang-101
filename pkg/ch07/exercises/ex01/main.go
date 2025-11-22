package main

import "fmt"

type AllowTypes interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

func DoublerNumber[T AllowTypes](number T) T {
	return number * 2
}

func main() {
	fmt.Println(DoublerNumber(2))
	fmt.Println(DoublerNumber(10.45))
}
