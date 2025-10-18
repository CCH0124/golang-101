package main

import "fmt"

func main() {
	greetings := []string{
		"Hello",
		"Hola",
		"नमस्कार",
		"こんにちは",
		"Привіт",
	}
	var slice1 = greetings[:2]
	var slice2 = greetings[1:4]
	var slice3 = greetings[3:]
	fmt.Println(slice1)
	fmt.Println(slice2)
	fmt.Println(slice3)
}
