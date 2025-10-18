# Composite Types

本章節深入探討了 Go 語言中用於組合基本資料類型的核心結構，包括陣列、切片、字串、映射和結構體，並特別強調了在實務開發中，如何基於性能和慣用法（Idiomatic Go）來選擇和操作這些類型。

複合類型是 Go 語言中內建的資料結構，用於儲存一組值。理解它們的底層機制和慣用法對於編寫高效能且可維護的 Go 程式至關重要。

## 陣列 (Arrays) - 剛性容器

陣列是具有固定長度和相同元素類型的值序列。然而，由於它們固有的剛性，在 Go 程式設計中**很少直接使用**。

* **類型與大小 (Size as Part of Type):** 陣列的大小是其類型定義的一部分。例如，`int` 和 `int` 被視為兩種不同的類型，它們之間不能直接轉換或賦值。
* **宣告:** 必須在編譯時確定長度。

    ```go
    var xint                 // 宣告一個長度為 3 的 int 陣列，元素皆初始化為 0
    var x = [3]int{10, 20, 30} // 宣告一個長度為 3 的 int 陣列，並給元素值
    var y = [...]int{10, 20, 30} // 編譯器推斷長度為 3
    ```

* **比較 (Comparison):** 兩個陣列若長度和元素值均相等，則它們相等。
* **用途:** 陣列的主要職責是作為切片 (Slice) 的底層儲存結構 (backing store)。

## 切片 (Slices) - 彈性與效率的基礎

切片 (Slice) 是 Go 語言中最常用的序列數據結構。與陣列不同，切片的長度不是其類型的一部分，使其具有動態擴展的能力。

#### 核心概念

* **長度 (Length, `len()`):** 切片中當前擁有的元素數量。
  * 將 nil 切片傳遞給 len 會返回 0。
* **容量 (Capacity, `cap()`):** 底層陣列中從切片的起始點算起，仍然可用的連續內存空間的總數。
  * 相比於 `len`，它的使用頻率要低得多。大多數時候，`cap` 用來檢查切片是否足夠大以存放新數據，或者是否需要調用 `make` 來創建一個新的切片。
* **零值 (Zero Value):** 切片的零值是 `nil`。`nil` 切片長度和容量都是 0，且可與 `nil` 比較。
* **不可比較性 (Non-comparable):** 切片之間不允許直接使用 `==` 或 `!=` 進行值比較，只能與 `nil` 比較。從 Go 1.21 開始，應使用 `slices.Equal` 或 `slices.EqualFunc` 進行值比較。
* **宣告**:

    ```go
    var x = []int{10, 20, 30}
    var x [][]int // 多維
    ```

#### 操作與性能優化

1. **增長 (Growth - `append`)**

* `append` 函數用於向切片追加元素。由於 Go 採用傳值呼叫 (Call by Value)，`append` 可能會導致底層陣列重新分配 (當長度達到容量時)，因此**必須將結果重新賦值給原變數**。

* **容量增長策略 (Capacity Growth Strategy):** 在 Go 1.18 之後，當容量小於 256 時，通常會將容量翻倍。對於更大的切片，增長率會逐漸收斂到 25% 左右。

```go
    var x = []int{1, 2, 3}
    x = append(x, 4)
    x = append(x, 5, 6, 7)
    y := []int{20, 30, 40}
    x = append(x, y...) // Slice 附加至另一個 Slice 使用 ...
```

> 如果你忘記分配 `append` 返回的值，將會是編譯時錯誤。

2. **預先分配 (Pre-allocation - `make`):**

如果*預先知道切片所需的元素數量*，應使用 `make` 函數來設定初始容量，以減少運行時內存重新分配 (reallocation) 的開銷，進而**降低垃圾回收器 (GC) 的工作負載**。

```go
    // 創建一個長度為 0，但容量為 10 的切片，適用於後續使用 append
    x := make([]int, 0, 10)

    // 創建一個長度和容量均為 5 的切片，元素初始化為零值
    y := make([]int, 5)
    // 警告：如果對此切片使用 append，新元素將排在 5 個零值之後。
```

3. **切片表達式 (Slice expression):**

* 使用切片表達式 (如 `x[:2]`) 創建的子切片會與原始切片**共享底層儲存**。修改任一切片的元素，都會影響其他共享內存的切片。

* **高階警告 (Senior Warning - Capacity Issue):** 為了防止對子切片執行 `append` 時意外覆寫到原始切片後續的數據，建議使用三部分切片表達式 (Full Slice Expression) 來限制子切片的容量，使其容量等於其長度。
  * 完整的切片表達式包含第三部分，它表示父切片容量中可供子切片使用的最後位置。

```go
    x := make([]string, 0, 5)
    x := []string{"a", "b", "c", "d"}
    // y 的長度為 2，容量也被限制為 2 (4 - 2 = 2)，append 時不會影響 x 的底層數據
    y := x[2:4:4]
```

```go
    x := make([]
    string, 0, 5)
    x = append(x, "a", "b", "c", "d")
    y := x[:2]
    z := x[2:]
    fmt.Println(cap(x), cap(y), cap(z))
    y = append(y, "i", "j", "k")
    x = append(x, "x")
    z = append(z, "y")
    fmt.Println("x:", x)
    fmt.Println("y:", y)
    fmt.Println("z:", z)
    //5 2 2
    //x: [a b c d x]
    //y: [a b i j k]
    //z: [c d y]
```

4. **複製 (Copying - `copy`):**

若需要一個與原始切片完全獨立的資料結構，必須使用內建的 `copy(destination slice, source slice)` 函數進行深層複製。

```go
x := []int{1, 2, 3, 4}
y := make([]int, 4)
num := copy(y, x)
fmt.Println(y, num)
//  [1 2 3 4] 4
```

> 這個函數會將盡可能多的值從來源複製到目的地，數量以較短的切片為限，並返回複製的元素數量。x 和 y 的容量並不重要，重要的是長度。

5. Converting Arrays to Slices

如果你有一個陣列，你可以使用切片表達式從中取得一個切片。這是一種將陣列傳給只接受切片的函數的有效方法。可以使用 `[:]`

```go
xArray := [4]int{5, 6, 7, 8}
xSlice := xArray[:]
```

相反，slice 可以轉換成 array。此時，*當你將切片轉換為陣列時，切片中的資料會被複製到新的記憶體。*

```go
xSlice := []
int{1, 2, 3, 4}
xArray := [4]int(xSlice)
smallArray := [2]int(xSlice)
xSlice[0] = 10
fmt.Println(xSlice)
fmt.Println(xArray)
fmt.Println(smallArray)
// [10 2 3 4] 
// [1 2 3 4] 
// [1 2]
```

### 3. 字串、符文與位元組 (Strings, Runes, and Bytes)

* **字串本質 (String Anatomy):** 字串在 Go 中是不可變的 (immutable)，其底層由一系列位元組 (bytes) 組成，通常假設採用 UTF-8 編碼。
  * 由於字串是不可變的，它們不會像切片的切片那樣有修改上的問題
* **符文 (Rune):** `rune` 是 `int32` 的別名，用於表示單個 Unicode 碼點 (char)。字串面值預設類型為 `string`，符文面值預設類型為 `rune`。
* **危險操作 (Indexing Hazard):** 由於字串的索引和切片操作都是**以位元組為單位**，直接使用索引或切片可能會在處理多位元組字符 (如中文、表情符號) 時發生錯誤的截斷，導致非法的 UTF-8 序列。
  * **慣用法 (Idiomatic Usage):** 應使用 `for-range` 迴圈來安全地迭代字串中的**符文(rune)**（碼點），而非位元組(bytes)。
* **轉換 (Conversion):** 字串可以與 `[]byte` 或 `[]rune` 互相轉換。

    ```go
    var a rune    = 'x'
    var s string  = string(a)
    var b byte    = 'y'
    var s2 string = string(b)
    ```

    ```go
    var s string = "Hello, 🌞"
    var bs [] byte = []byte(s)
    var rs [] rune = []rune(s)
    fmt.Println(bs)
    fmt.Println(rs)
    //[72 101 108 108 111 44 32 240 159 140 158] 
    //[72 101 108 108 111 44 32 127774]
    ```

### 4. 映射 (Maps) - 鍵值對儲存

Map 是 Go 內建的無序鍵值對集合，採用哈希表 (Hash Map) 實現。其類型表示為 `map[keyType]valueType`。

* **零值與初始化:** Map 的零值是 `nil`。**對 `nil` Map 進行寫入操作會導致 `panic`**。因此，Map 應使用字面量 `{}` 或 `make()` 進行初始化。
  * 可以使用 `make` 來創建一個具有默認大小的映射，但可以超過最初指定的大小增長

    ```go
    var nilMap map[string]int //nilMap 被宣告為一個鍵為字串、值為整數的 map。
    teams := map[string][] string {
    "Orcas": []
    string{"Fred", "Ralph", "Bijou"},
    "Lions": []
    string{"Sarah", "Peter", "Billie"},
    "Kittens": []
    string{"Waldo", "Raul", "Ze"},
    }
    ```
* **鍵限制:** 鍵的類型必須是**可比較類型** (Comparable Type)。因此，`Slice` 或 `Map` 本身不能作為 `Map` 的鍵。
* **讀取操作:** 讀取 `Map` 中不存在的鍵會返回其值類型的零值。
* **Comma Ok Idiom:** Map 讀取操作返回兩個值：實際值和一個布林值 `ok`。這用於區分鍵值對是否存在，即便值為零值。

    ```go
    v, ok := m["key_name"] // ok 為 true 則表示鍵存在
    ```

* **刪除操作:** `delete` 函數接受一個 `map` 和一個 `key`，然後刪除具有指定鍵的鍵值對。
  * 不返回任何值
* **清空操作:** `clear` 函數，將其 `map` 長度設為零
* **Map 作為集合 (Sets):** Go 沒有內建 Set 類型，通常使用 Map 模擬，將要儲存的元素作為鍵，值類型設為 `bool` 或零位元組的 `struct{}` (用於內存優化)。

* **不可比較性與迭代順序:** Map 不可比較，但 Go 1.21 提供了 `maps.Equal` 函數。此外，Map 的迭代順序是**不確定且隨機的**，這是一個安全特性，旨在防止 Hash DoS 攻擊。

### 5. 結構體 (Structs) - 數據聚合

結構體是用來聚合多個不同類型欄位 (fields) 的自定義複合類型。

* **宣告與初始化:** 必須使用 `type` 關鍵字定義。初始化時，欄位若未明確賦值，則取其零值。

```go
type person struct {
    name string
    age  int
    pet  string
}

```

* **賦值慣用法:** 建議使用**帶有欄位名稱**的結構體字面量進行初始化，這提高了程式碼的可讀性和可維護性。

```go
    beth := person{
        age:  30,
        name: "Beth", // 使用欄位名賦值，順序可變
    }
```

* **匿名結構體 (Anonymous Structs):** 可以在不預先給予類型名稱的情況下定義和實例化結構體，常用於臨時數據結構 (如 JSON unmarshaling/marshaling)。

```go
var person struct {
 name string
 age  int
 pet  string
} 
person.name = "bob"
person.age = 50
person.pet = "dog" 
pet := struct {
 name string
 kind string
}{
    name: "Fido",
    kind: "dog"
}
```

* **比較與轉換:** 結構體只有在所有欄位都可比較時才可比較。結構體之間的類型轉換要求**欄位名稱、順序和類型**完全一致。

> 在設計 API 時，應當優先使用 Struct 而非 Map 來傳遞複雜數據。Struct 提供了強類型保障和明確的欄位定義，有助於數據流分析和提升代碼可讀性，從而解決 Map 缺乏結構和類型約束的限制。