package main

import "fmt"

func main() {
	const value = 10.3
	i := value
	// 明確指定類型
	var f float64 = value
	fmt.Println(i)
	fmt.Println(f)
}
