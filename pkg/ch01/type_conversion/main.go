package main

import "fmt"

func main() {
	var total int = 50
	var avg float64 = 7.5

	// 錯誤：Go 不允許直接混合不同數值型別運算
	// sum := total + avg

	// 正確：必須顯式轉換型別
	// 將 int 轉為 float64 才能相加
	resultFloat := float64(total) + avg

	// 將 float64 轉為 int 相加（會截斷小數部分）
	resultInt := total + int(avg) // avg 7.5 變成 7

	fmt.Println(resultFloat) // 57.5
	fmt.Println(resultInt)   // 57

	// 嚴格性：不能將數值隱式視為布林值（無真值性）
	// if total { ... } // 錯誤：無法將 int 當作 bool
	if total != 0 {
		fmt.Println("Total is non-zero.")
	}
}
