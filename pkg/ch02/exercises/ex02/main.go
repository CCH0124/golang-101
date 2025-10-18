package main

import "fmt"

func main() {
	var message = "Hi 👩 and 👨"
	var rs []rune = []rune(message)
	// fmt.Println(rs[3]) // number
	fmt.Println(string(rs[3]))
}
