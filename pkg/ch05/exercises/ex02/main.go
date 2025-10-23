package main

import (
	"fmt"
)

func UpdateSlice(s []string, v string) {
	s[len(s)-1] = v
	fmt.Println("in UpdateSlice:", s)
}
func GrowSlice(s []string, v string) {
	s = append(s, v)
	fmt.Println("in GrowSlice:", s)
}
func main() {
	s := []string{
		"Itachi",
		"Madara",
		"Naruto",
	}
	UpdateSlice(s, "Java")
	fmt.Println("after UpdateSlice:", s)
	GrowSlice(s, "cilium")
	fmt.Println("after GrowSlice:", s)
}
