# Chapter 4. Blocks, Shadows, and Control Structures

本章節主要探討 Go 語言的程式碼組織結構（區塊與作用域）以及控制流程語法，這對於編寫具備清晰邏輯和高可維護性的程式碼至關重要。

### 區塊（Blocks）

Go 語言透過「區塊」（blocks）的概念來管理變數、常數、型別和函式的宣告範圍。

1. **區塊定義**：每個發生宣告的位置都被稱為一個區塊。
    * **Package Block（套件區塊）**：在任何函式外部宣告的變數、常數、型別和函式都位於此區塊。
    * **File Block（檔案區塊）**：`import` 陳述式定義的名稱存在於檔案區塊，這些名稱在包含 `import` 陳述式的檔案內有效。
    * **函式級區塊**：函式頂層宣告的所有變數（包括函式參數）都屬於此區塊。
    * **內部區塊**：在函式內，每組大括號 `{}` 都定義了一個新的區塊。控制結構（如 `if`、`for`、`switch`）也定義了自己的區塊。
2. **存取規則**：內部區塊可以存取在任何外部區塊中定義的識別符（identifier）。

### 變數遮蔽（Shadowing Variables）

當內部區塊中的宣告與外部區塊中的識別符具有相同的名稱時，即發生了「遮蔽」（shadowing）現象。

1. **影響**：只要遮蔽變數存在，被遮蔽的外部變數就無法被存取。
2. **`:=` 的風險**：使用短宣告運算子 `:=` 很容易意外造成變數遮蔽。`:=` 運算子可以同時宣告和賦值給多個變數，但它只會重複使用在**當前區塊**中已宣告的變數。如果左側存在至少一個新變數，而其他同名變數是在外部區塊中宣告的，那麼這些外部變數將被遮蔽。

    ```go
            func main() {
                x := 10 // 外部 x
                if x > 5 {
                    x, y := 5, 20 // 內部區塊宣告了新的 x，遮蔽了外部 x
                    fmt.Println(x, y) // 輸出: 5 20
                }
                fmt.Println(x) // 輸出: 10 (外部 x 未被改變)
            }
    ```

3. **套件名稱遮蔽**：宣告與匯入套件同名的變數（例如 `fmt`）會使該套件在該區塊內無法使用。

4. **Universe Block 識別符遮蔽**：Go 的內建型別（如 `int`, `string`）、常數（如 `true`, `false`）和函式（如 `make`, `close`）定義在 `universe block` 中。雖然可以被遮蔽，但**應嚴格避免**這樣做，因為這會導致程式行為怪異。

### if 條件判斷

Go 語言中的 `if` 語句與大多數程式語言相似，但有其獨特之處。

1. **語法**：條件判斷式周圍**不使用**括號。
2. **變數作用域增強**：Go 允許在 `if` 條件判斷前宣告變數，這些變數的作用域限定在 `if`、`else if` 和 `else` 的整個結構中。

    ```go
        // 變數 n 僅在整個 if/else if/else 結構中有效
        if n := rand.Intn(10); n == 0 {
            fmt.Println("That's too low") 
        } else if n > 5 {
            fmt.Println("That's too big:", n) 
        } else {
            fmt.Println("That's a good number:", n) 
        }
        // 在此之後，n 是未定義的。
    ```

### for 迴圈 (for, Four Ways)

`for` 是 Go 語言中唯一的迴圈關鍵字，具備四種基本形式。

#### 1. 完整的 C 樣式 for 語句 (The Complete for Statement)

這是傳統的初始化、條件判斷和增減量結構。

* **初始化部分**：必須使用 `:=` 進行變數初始化；`var` 在此處不合法。
* **分號**：三個部分之間用分號分隔。
* **靈活性**：可以省略其中一個或多個部分。

```go
i := 0
for ; i < 10; i++ {
    fmt.Println(i)
}
```

```go
for i := 0; i < 10; {
    fmt.Println(i)
    if i % 2 == 0 {
        i++
    } else {
        i+=2
    }
}
```

#### 2. 僅條件判斷 for 語句 (The Condition-Only for Statement)

當省略初始化和增減量部分時，不應包含分號。此形式等同於其他語言中的 `while` 語句。

```go
i := 1
for i < 100 {
    fmt.Println(i)
    i = i * 2
}
```

#### 3. 無限 for 語句 (The Infinite for Statement)

省略所有部分，迴圈將永遠執行。應在迴圈體內使用 `break` 或 `return` 來終止執行。

#### 4. for-range 語句 (The for-range Statement)

用於迭代複合型別（如strings, arrays, slices, and maps ...）中的元素。

* **返回變數**：通常返回兩個值：位置/鍵（index/key）和該位置的值（value）。
* **忽略變數**：若不需要位置/鍵，應使用底線 `_` 忽略，以避免編譯器報告「變數未讀取」的錯誤。
* **迭代特性**：
  * **Map**：迭代順序是隨機變化的（這是一個安全特性），不應依賴特定的順序。
  * **String**：`for-range` 會迭代 **rune**（Unicode 碼點），而非位元組（byte）。第一個變數是位元組偏移量，第二個變數是 `rune` 型別。
  * **值複製**：`for-range` 迴圈中的值變數 (`v`) 是被迭代元素的**副本**。修改這個副本不會影響原始複合型別中的元素。
  * **Go 1.22 變更**：自 Go 1.22 版本起，預設行為是在每次迭代中建立新的索引和值變數，這解決了一個常見的併發錯誤（如果您在 `go.mod` 中指定了 `go 1.22` 或更高版本）。

```go
evenVals := []int{2, 4, 6, 8, 10, 12}
for i, v := range evenVals {
    fmt.Println(i, v)
}
```

#### break 和 continue

* `continue` 可用於避免巢狀的 `if/else` 結構，使程式碼更清晰易讀。
* 預設情況下，`break` 和 `continue` 僅作用於直接包含它們的 `for` 迴圈。

#### 迴圈標籤 (Labeling)

當需要跳出或跳過巢狀迴圈中的外層迴圈迭代時，可以使用標籤（Label）。標籤必須緊接著 `for` 關鍵字，且 `go fmt` 會將標籤縮排至與外部大括號相同的層級。

### switch 選擇結構

Go 語言中的 `switch` 語句相比於 C 衍生語言更加實用和安全。

#### 1. 表達式 switch (Expression Switch)

* **語法**：不使用括號圍繞被比較的表達式。
* **短變數宣告**：與 `if` 語句相似，可以在 `switch` 關鍵字後宣告變數，其作用域涵蓋所有 `case` 分支。
  * **範例**：`switch size := len(word); size { ... }`
* **無貫穿（No Fall-through）**：Go 的 `case` 預設不會貫穿到下一個 `case`。雖然存在 `fallthrough` 關鍵字，但應盡量避免使用，以保持邏輯清晰。
* **可比較性**：`switch` 語句可以對任何可使用 `==` 比較的型別進行操作（slice、map、channel或包含這些型別的 struct 除外）。
* **在迴圈內跳出**：在 `for` 迴圈內，如果 `switch` 語句中的 `break` 沒有指定標籤，它只會跳出 `case` 區塊，不會跳出 `for` 迴圈。若要跳出外部 `for` 迴圈，必須使用標籤，例如 `break loop`。

#### 2. 空白 switch (Blank Switches)

當省略 `switch` 後的比較值時，每個 `case` 必須包含一個布林表達式。

* **優勢**：允許使用任何布林比較邏輯，而不僅限於相等性檢查。
* **規範**：如果所有 `case` 都是對同一變數進行相等性比較，則應使用表達式 `switch` 而非空白 `switch`。

#### `if` 與 `switch` 的選擇

從功能上講，一系列的 `if/else` 陳述式和空白 `switch` 陳述式功能重疊。

* **慣例**：當存在多個相關的條件判斷時（例如互斥或依賴關係），應優先選擇 **`blank switch`**。這使比較邏輯更具可見性，並強調整個結構是相關聯的一組關注點。

### goto — 對，就是 goto

`goto` 關鍵字在 Go 語言中存在，但建議幾乎不要使用。

1. **限制**：Go 限制了 `goto` 的跳轉範圍，**禁止**跳過變數宣告（variable declarations）或跳入內部區塊。
    * **編譯錯誤範例**：`goto skip jumps over declaration of b`。
2. **有效用途**：
    * 在函式結尾處執行複雜的清理邏輯或最終步驟，避免重複程式碼。
    * 儘管如此，通常應嘗試使用標籤化的 `break` 或 `continue` 替代。
