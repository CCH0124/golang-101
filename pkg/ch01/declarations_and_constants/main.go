package main

import "fmt"

const CompileTimeVal = 100 // 非型別常數，可靈活使用

func main() {
	// 1. 函式內常用的短宣告
	x := 5
	s := "Go Language"

	// 2. 宣告零值：明確使用零值時的最佳寫法
	var counter int // counter = 0

	// 3. 需指定非預設型別時
	var b byte = 255 // 避免使用 b := byte(255)

	// 4. 變數遮蔽範例 (Bad Practice!)
	if x > 0 {
		x, y := 20, 30             // 在內部區塊 (Block) 宣告了新的 x 和 y
		fmt.Println("Inner x:", x) // 20
	}
	fmt.Println("Outer x:", x) // 5 (外部的 x 未被改變)
}
