## 1. 宣告與呼叫函式 (Declaring and Calling Functions)

Go 語言的函式宣告語法與其他 C 衍生語言相似，但有其獨特之處。

### 函式結構

1. **宣告格式**：使用 `func` 關鍵字，後接函式名稱、輸入參數列表（在括號內），以及回傳型別。
2. **參數順序**：在參數列表中，參數名稱在前，型別在後。Go 是一種嚴格型別語言，必須指定參數的型別。
3. **回傳型別**：回傳型別寫在輸入參數的結束括號和函式主體的開頭大括號之間。

### 慣例與簡化

* **連續相同型別參數**：當有多個連續輸入參數的型別相同時，可以只在最後一個參數後指定型別，以簡化宣告。

```go
// 慣例寫法：只有最後一個參數需指定型別
func div(num, denom int) int {
    if denom == 0 {
        return 0
    }
    return num / denom
}
//
```

## 2. 模擬具名與可選參數 (Simulating Named and Optional Parameters)

Go 語言不支援具名（Named）和可選（Optional）輸入參數。

* **模擬方法**：若函式需要模擬具名或可選參數，工程師通常會定義一個 `struct`，其欄位對應所需的參數。將此 `struct` 傳入函式，可以透過結構體字面量（struct literal）來選擇性地賦值。
* **設計原則**：如果一個函式參數過多，可能意味著該函式過於複雜。

```go
type MyFuncOpts struct {
    FirstName string
    LastName  string
    Age       int
}

func MyFunc(opts MyFuncOpts) error {
    // 執行邏輯
}

func main() {
    // 呼叫時只傳入部分需要的參數，其他參數保留為零值
    MyFunc(MyFuncOpts{
        LastName: "Patel",
        Age:      50,
    }) 
    //
}
```

## 3. 可變參數與切片 (Variadic Input Parameters and Slices)

Go 支援可變參數（Variadic Parameters），允許函式接受任意數量的特定型別輸入。

* **宣告與使用**：可變參數必須是參數列表中的**最後一個**參數（或唯一一個）。使用三個點 `...` 放在型別之前來宣告 (e.g., `...int`)。
* **內部機制**：在函式內部，可變參數會被視為指定型別的切片（slice）。
* **切片合併**：當將一個切片附加（append）到另一個切片時，也需要使用 `...` 運算子，將源切片展開成獨立的值。

```go
func addTo(base int, vals ...int) []int {
    out := make([]int, 0, len(vals))
    for _, v := range vals {
        out = append(out, base+v)
    }
    return out
}

func main() {
    fmt.Println(addTo(3))
    fmt.Println(addTo(3, 2))
    fmt.Println(addTo(3, 2, 4, 6, 8))
    a := []int{4, 3}
    fmt.Println(addTo(3, a...))
    fmt.Println(addTo(3, []int{1, 2, 3, 4, 5}...))
 }
```

## 4. 多重回傳值 (Multiple Return Values)

Go 語言最顯著的特性之一是其支援多重回傳值，這對於錯誤處理至關重要。

### 標準多重回傳值

* **型別清單**：多個回傳值的型別必須用括號 `()` 包裹，並用逗號分隔。
* **錯誤慣例**：錯誤（`error`）型別應**永遠**作為函式的最後一個回傳值。成功時回傳 `nil`。
* **賦值**：呼叫多重回傳值的函式時，通常需要使用多個變數來接收結果。

```go
func divAndRemainder(num, denom int) (int, int, error) {
    if denom == 0 {
        return 0, 0, errors.New("cannot divide by zero")
    }
    return num / denom, num % denom, nil
}

func main() {
    result, remainder, err := divAndRemainder(5, 2)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println(result, remainder)
}
```

### 忽略回傳值

* **顯式忽略**：若不需要接收某個回傳值，應使用底線 `_` 符號來顯式忽略，這是慣例做法。
* **隱式忽略**：Go 允許隱式忽略所有回傳值 (e.g., `divAndRemainder(5,2)` 即可)，但通常只有在呼叫 `fmt.Println` 這類函式時才建議這麼做。

### 具名回傳值 (Named Return Values)

* **機制**：在回傳型別列表中為回傳值指定名稱，這些名稱會成為函式內的**預先宣告變數**，並初始化為各自的零值。
* **作用域**：具名回傳值的名稱僅限於函式內部有效。
* **潛在問題**：具名回傳值可能被函式內部的變數宣告遮蔽（shadowing），導致邏輯混亂。

```go
func divAndRemainder(num, denom int) (result int, remainder int, err error) {
    if denom == 0 {
        err = errors.New("cannot divide by zero")
        return result, remainder, err
    }
    result, remainder = num/denom, num%denom
    return result, remainder, err
}
```

### 空白回傳 (Blank Returns)

* **語法**：如果函式使用了具名回傳值，可以只寫 `return` 而不指定回傳的變數。這會回傳具名變數的當前值。
* **嚴格避免**：Go 社群通常認為空白回傳是一種**嚴重的設計缺陷（severe misfeature）**，應嚴格避免使用，以保持程式碼的清晰度。

## 5. 函式是數值 (Functions Are Values)

Go 語言中的函式是一級公民（first-class concept），可以像任何其他數值一樣被賦值、儲存和傳遞。

### 函式型別與宣告

* **型別簽章**：函式的型別由 `func` 關鍵字、參數型別和回傳值型別組成，這被稱為函式的「簽章」（signature）。

```go
var myFuncVariable func(string) int
func f1(a string) int {
    return len(a)
} 
 
func f2(a string) int {
    total := 0
    for _, v := range a {
        total += int(v)
    }
    return total
} 
func main() {
    var myFuncVariable func(string) int
    myFuncVariable = f1
    result := myFuncVariable("Hello")
    fmt.Println(result) 
 
    myFuncVariable = f2
    result = myFuncVariable("Hello")
    fmt.Println(result)
}
```

* **零值**：函式變數的零值是 `nil`。嘗試執行一個 `nil` 函式變數將導致 `panic`。
* **匿名函式**：可以在函式內部定義新的函式並將其賦值給變數，這類函式沒有名稱。

```go
func main() {
    // anonymous
    f := func(j int) {
        fmt.Println("printing", j, "from inside of an anonymous function")
    }
    for i := 0; i < 5; i++ {
        f(i)
    }
}

func main() {
    for i := 0; i < 5; i++ {
        func(j int) {
            fmt.Println("printing", j, "from inside of an anonymous function")
        }(i) // call them immediately
    }
 }
```

* **Functoion Type**: 

```go
type opFuncType func(int, int) int
```

### 閉包 (Closures)

在函式內部宣告的函式被稱為「閉包」（closure）。

* **特性**：閉包可以存取並修改其外部函式中宣告的變數（即捕獲了外部作用域的變數）。
* **應用**：閉包經常用於限制函式的作用域，或作為參數傳遞給其他函式，以達到複雜的資料抽象與處理。

```go
func main() {
    a := 20
    f := func() {
        fmt.Println(a) // 閉包讀取外部 a
        a = 30         // 閉包修改外部 a
    }
    f() // 輸出 20
    fmt.Println(a) // 輸出 30
}
//
```

> 在閉包內使用 `:=` 而不是 `=` 會創建一個新的變數 a，當閉包結束時該變數也隨之消失。在使用內部函數時，要小心使用正確的賦值運算符，特別是當左側有多個變數時。

#### Passing Functions as Parameters

由於函數本身也是*值*，並且你可以透過其參數和返回類型來指定函數的類型，因此你可以將函數作為參數傳入其他函數。

```go
type Person struct {
 FirstName string
 LastName  string
 Age       int
} 

people := []Person{
 {"Pat", "Patterson", 37},
 {"Tracy", "Bobdaughter", 23},
 {"Fred", "Fredson", 18},
}
fmt.Println(people)

sort.Slice(people, func(i, j int) bool {
    return people[i].LastName < people[j].LastName
})
fmt.Println(people)
```

#### Returning Functions from Functions

* 從一個函數返回一個閉包
* Go 也使用閉包來實現資源清理，透過 defer 關鍵字來完成。

```go
func makeMult(base int) func(int) int {
    return func(factor int) int {
        return base * factor
    }
}
```

## 6. defer 機制

`defer` 關鍵字用於將清理程式碼（cleanup code）延遲到包含它的函式即將退出時才執行。

* **用途**：確保資源（如檔案、網路連線、鎖定）在函式結束時被釋放，無論函式是正常返回還是發生錯誤。
* **執行順序**：如果一個函式內有多個 `defer` 陳述式，它們將以後進先出（LIFO）的順序執行。
* **參數求值**：當 `defer` 語句被執行時，其所引用的函式參數會**立即**被求值，但函式體本身會被延遲執行。
* **與具名回傳值結合**：`defer` 的閉包可以檢查或修改外部函式的具名回傳值，這常用於資料庫事務的提交/回滾或錯誤資訊的上下文添加。

```go
func main() {
    // 參數在 defer 執行時立即求值 (a=10)
    a := 10
    defer func(val int) { 
        fmt.Println("first:", val) // 輸出 10
    }(a)
    
    // 參數在 defer 執行時立即求值 (a=20)
    a = 20
    defer func(val int) {
        fmt.Println("second:", val) // 輸出 20
    }(a)
    
    a = 30
    fmt.Println("exiting:", a) // 輸出 30
    // 輸出順序: exiting: 30, second: 20, first: 10
}
//
```

## 7. Go 是傳值呼叫 (Go Is Call by Value)

Go 是一種嚴格的「傳值呼叫」（call-by-value）語言。

* **核心原則**：當變數作為參數傳遞給函式時，Go 總是會建立該變數的**副本**。
* **值型別行為**：對於基本型別（如 `int`, `string`）和結構體（`struct`），傳入函式的是其副本。因此，在函式內對副本的修改，不會影響到呼叫者作用域的原始變數。
* **切片與映射的特殊行為**：儘管 Go 堅持傳值呼叫，但 `map` 和 `slice` 的行為看起來像是傳參考（pass-by-reference）。
  * **原因分析**：這是因為 `map` 和 `slice` 內部是以指標（pointers）實現的。傳遞參數時，傳遞的是包含長度、容量和指向底層數據結構的**指標副本**。
    * **結果**：
      * **Map**：對 map 參數的修改（更新或刪除鍵值）會反映在原始變數中。
      * **Slice**：可以修改 slice 內部的元素值，但不能透過 `append` 改變其長度，並讓變動反映在原始變數上（因為改變長度或容量會產生新的 slice 結構副本，與原始結構脫鉤）。

**範例：值型別不變**

```go
type person struct {
    age  int
    name string
}

func modifyFails(i int, s string, p person) {
    i = i * 2      // 僅修改 i 的副本
    s = "Goodbye"  // 僅修改 s 的副本
    p.name = "Bob" // 僅修改 p (struct) 的副本
}

func main() {
    p := person{}
    i := 2
    s := "Hello"
    modifyFails(i, s, p)
    fmt.Println(i, s, p) // 輸出: 2 Hello {0 } (原始值未變)
}
//,
```