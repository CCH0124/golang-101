package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var store = make([]int, 0, 100)
	for i := 0; i < 100; i++ {
		store = append(store, rand.Intn(100))
	}
	fmt.Println(store)
}
