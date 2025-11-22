# Generics

泛型是 Go 1.18 版本引入的重大特性，旨在提高程式碼的重用性（減少重複程式碼，DRY原則）並增強編譯時的型別安全檢查，這是從靜態型別語言設計角度來看極為關鍵的改進。

## 1. 泛型減少重複程式碼並提高型別安全性 (Generics Reduce Repetitive Code and Increase Type Safety)

在 Go 引入泛型之前，如果開發者需要為不同型別（如 `int` 和 `string`）實作相同的資料結構（例如二元樹），通常有兩種不理想的解決方案：

1. **複製程式碼：** 為每種型別都撰寫一套重複的邏輯，這冗長且容易出錯。
2. **使用介面搭配 `any` (或舊版 `interface{}`)：** 將資料儲存為 `any` 型別，並透過介面定義行為（如排序）。然而，這種方法會犧牲 Go 的核心優勢：**編譯時的型別安全**。如果嘗試將不相容的型別插入資料結構，編譯器無法報錯，將導致程式在執行時發生恐慌（Panic）。

此外，缺乏泛型也限制了函式的抽象化能力。例如，標準函式庫中的 `sort.Slice` 必須依賴反射（Reflection），這會犧牲性能和型別安全。泛型（即型別參數）正是解決這一根本限制的工具。

## 2. 泛型的引入 (Introducing Generics in Go)

Go 語言的泛型實作強調**編譯速度快**、**程式碼可讀性高**和**執行效率好**。

#### 核心語法與結構

* **型別參數 (Type Parameters)：** 在宣告型別或函式名稱後，使用方括號 `[]` 來定義型別參數。例如，`Stack[T any]`。
* **型別限制 (Type Constraints)：** 型別參數的定義必須包含一個介面，用於限制哪些型別可以被替換。這確保了型別安全。
* **`any` 限制：** 這是預定義的識別符，表示任何型別都可以使用（等同於舊版 `interface{}`）。
* **`comparable` 限制：** 如果泛型程式碼內部需要使用比較運算子 `==` 或 `!=`，型別參數必須使用內建的 `comparable` 介面進行限制。

#### 泛型與零值

對於泛型型別，如果需要返回零值，不能簡單地返回 `nil` (因為 `nil` 對於 `int` 等值型別是無效的)。標準的做法是使用 `var` 宣告一個變數，因為 `var` 總是將變數初始化為其型別的零值。

**範例 (Stack 結構)：**

```go
type Stack[T any] struct {
    vals []T
}
// 這裡 T 必須滿足 any 限制
func (s *Stack[T]) Push(val T) {
    s.vals = append(s.vals, val)
 } 
 
func (s *Stack[T]) Pop() (T, bool) {
    if len(s.vals) == 0 {
        var zero T
        return zero, false
    }
    top := s.vals[len(s.vals)-1]
    s.vals = s.vals[:len(s.vals)-1]
    return top, true
 }
```

## 3. 抽象化演算法的泛型函式 (Generic Functions Abstract Algorithms)

泛型函式將型別參數放在函式名稱之後、變數參數之前。這使得開發者可以輕鬆編寫適用於多種型別的通用演算法，例如函數式編程中常見的 `Map`、`Reduce` 和 `Filter`。

**範例(Map):**

```go
 // Map turns a []T1 to a []T2 using a mapping function.
 // This function has two type parameters, T1 and T2.
 // This works with slices of any type.
 func Map[T1, T2 any](s []T1, f func(T1) T2) []T2 {
    r := make([]T2, len(s))
    for i, v := range s {
        r[i] = f(v)
    }
    return r
 } 
```

**範例 (Filter 函式)：**

```go
// T 可以是任何型別 (any)
func Filter[T any](s []T, f func(T) bool) []T {
    var r []T
    for _, v := range s {
        if f(v) {
            r = append(r, v)
        }
    }
    return r
}
// 此函式在編譯時確定 T 的具體型別，保證了型別安全。
```

## 4. 泛型與介面 (Generics and Interfaces)

在 Go 泛型中，任何介面都可以作為型別限制，不僅僅是 `any` 和 `comparable`。

* **介面作為限制：** 可以要求型別參數滿足特定的介面，例如 `[T fmt.Stringer]`，以確保該型別具有 `String() string` 方法。
* **泛型介面：** 介面本身也可以定義型別參數（例如 `type Differ[T any] interface { ... Diff(T) ... }`），這在定義型別之間互動的行為時非常有用。

## 5. 使用型別項指定運算子 (Use Type Terms to Specify Operators)

為了讓泛型程式碼能夠使用運算子（如 `+`, `-`, `%` 或比較運算子），Go 引入了**型別元素 (Type Elements)**，它由一個或多個**型別項 (Type Terms)** 組成，這些型別項之間用 `|` 分隔，定義在介面內部。

* **運算子支援：** 介面中列出的所有具體型別都必須支援該運算子，泛型程式碼才能使用它。
* **型別項的限制：** 包含型別元素的介面僅能作為**型別限制**使用，不能作為變數、欄位、參數或回傳值的型別。
* **波浪符號 (`~`)：** 如果希望型別項涵蓋該基礎型別以及所有**底層型別 (underlying type)** 為該基礎型別的**使用者自定義型別**，需要在型別項前加上 `~`。例如 `~int`。

**範例 (Integer 介面與波浪符號)：**

```go
 type Integer interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
        ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
 }
// 使用 ~int，則自定義型別 MyInt (底層為 int) 也能滿足此限制。
```

Go 1.21 新增的 `cmp` 套件即定義了 `Ordered` 介面，涵蓋了所有支援比較運算子的基礎型別。

## 6. 型別推斷與泛型 (Type Inference and Generics)

Go 支援型別推斷，可以根據傳入的參數自動推導出型別參數，從而簡化泛型函式的呼叫。但如果型別參數僅作為回傳值使用，編譯器無法推斷，則必須明確指定所有型別參數。

## 7. 型別元素限制常數 (Type Elements Limit Constants)

泛型型別的變數能被賦予的常數，必須對該型別元素中列出的**所有型別項都有效**。例如，如果一個 `Integer` 介面包含 `int8`，則不能將超過 127 的常數賦予該泛型型別的變數，否則會導致編譯錯誤。

## 8. 結合泛型函式與泛型資料結構 (Combining Generic Functions with Generic Data Structures)

在複雜的泛型資料結構（如二元搜尋樹）中，通常會使用一個**函式型別**來處理操作細節（如排序或比較），然後將此函式型別作為參數傳遞給資料結構。

**範例：**
若要建立一個通用的二元樹 `Tree[T]`，可以定義一個比較函式簽名 `OrderableFunc[T]`，並將其實例化為 `Tree` 結構的欄位。樹的內部方法（如 `Add` 和 `Contains`）隨後呼叫此函式進行比較。
可以使用 `cmp.Compare` 或方法的**方法表示式 (Method Expression)** 來提供這個比較函式。

## 9. 關於 comparable 的更多細節 (More on comparable)

雖然 `comparable` 介面是強大的工具，但使用時必須謹慎，特別是當它與其他介面型別結合時。

* 如果一個泛型函式接受 `comparable` 型別參數，但傳入的變數是一個**介面**，而該介面底層的值是一個**不可比較型別**（例如 Slice），那麼程式碼雖然可以編譯，但在執行時仍會觸發 `panic`。
* **總結：** 雖然 `comparer` 限制了型別參數，但它無法在編譯時檢查所有潛在的執行時不可比較性，因此在使用介面作為 `comparable` 型別的替代時，開發者仍需注意。

## 10. 泛型未包含的特性 (Things That Are Left Out)

Go 語言的泛型實作選擇是克制的，並未包含其他語言常見的許多特性。

* **運算子重載 (Operator Overloading)：** Go 不會新增此功能。
* **參數化方法 (Parameterized Methods)：** 函式不能在方法層級定義額外的型別參數。這意味著不能實現類似 `mySlice.Map().Reduce()` 的方法鏈式呼叫。
* **可變參數型別參數 (Variadic Type Parameters)：** 不支援。
* **其他進階特性：** 包含專門化 (Specialization)、柯里化 (Currying) 和元編程 (Metaprogramming)。

以下是這三個特性及其在 Go 語言中被排除的原因的詳細解釋：

##### 1. 特化（Specialization）

**定義：**
特化（Specialization）指的是在程式設計中，允許為泛型函數或方法（generic function or method）針對特定的類型參數提供一個或多個「型別特定版本」（type-specific versions），以覆寫或補充泛型預設的行為。

**運作模式：**
當您使用泛型定義了一個適用於多種型別的函數時，特化允許您為某些型別（例如 `int` 或 `string`）編寫一個單獨的、非泛型的實作。這個特化版本可以在該特定型別下提供更優化的效能或獨特的邏輯。

**在 Go 語言中被排除的原因：**
Go 語言的設計決定是**不支援函數或方法重載（overloading）**。缺乏重載是為了讓程式碼更清晰、更易於理解，避免在追蹤程式執行時產生歧義。由於特化本質上就是一種基於型別參數的重載機制，如果將其加入泛型，將會破壞 Go 語言不重載的核心設計原則。

##### 2. 柯里化（Currying）

**定義：**
柯里化（Currying）允許根據另一個泛型函數或型別，透過指定部分型別參數來「部分實例化」（partially instantiate）該函數或型別。

**運作模式：**
柯里化是源自函數式編程（Functional Programming）的概念，它將一個接受多個參數的函數轉換為一系列接受單個參數的函數。在泛型的背景下，它允許您固定泛型結構或函數中的一部分型別參數，從而衍生出一個新的、更具體的泛型結構或函數。

**在 Go 語言中的背景：**
Go 雖然不是一個純粹的函數式語言，但它支援高階函數（higher-order functions），即函數可以作為參數傳入或作為回傳值傳出。然而，Go 選擇保持其泛型實作的簡潔性，避免引入像柯里化這樣複雜的泛型型別參數操作機制。

##### 3. 元編程（Metaprogramming）

**定義：**
元編程（Metaprogramming）是一種技術，允許您指定在「編譯時期」（compile time）執行的程式碼，而這些程式碼的目標是產生將在「執行時期」（runtime）運行的程式碼。

**運作模式：**
元編程使得程式能夠操作自身的程式碼，例如，在編譯前根據某些規範（如資料庫 schema 或介面定義）自動生成大量的重複程式碼。這在一些追求極致抽象和簡潔的語言中很常見。

**在 Go 語言中的替代方案與排除原因：**
Go 語言雖然沒有原生支援在語言層面進行編譯時元編程，但它提供了多種工具來處理類似的需求，這些工具通常比內建的元編程機制更具體且易於控制：

* **`go generate` 工具：** Go 開發人員通常使用 `go generate` 這個標準工具，它根據原始碼中的特殊註釋來執行外部程式（如 `stringer` 或 `protoc`），從而產生新的 Go 程式碼。這是在編譯前產生執行時程式碼的標準方式。
* **反射（Reflection）：** Go 的 `reflect` 套件允許程式在執行時檢查、修改甚至建立型別和值。雖然這是在**執行時**操作程式碼，但它滿足了許多元編程中動態處理型別的需求，例如資料格式的編組（marshaling）和解組（unmarshaling），如 `encoding/json` 套件即大量使用反射。然而，Go 鼓勵在程式邊界處（與外部資料交換時）使用反射，並警告在一般業務邏輯中應謹慎使用，因為反射較慢且更脆弱。

總體來說，Go 語言的泛型實作強調實用性，優先考慮減少重複程式碼（DRY原則）和提高型別安全，同時避免引入會使語言核心變得龐大或模糊程式碼執行流程的複雜特性。

### 11. 慣用語 Go 與泛型 (Idiomatic Go and Generics)

泛型改變了 Go 的慣用寫法：

* **替代品替換：** 應使用 `any` 替換 `interface{}`；應使用泛型函式處理多種數值型別，而不是單獨依賴 `float64`。
* **性能權衡：** 將原本使用**介面參數**的函式改寫為**泛型型別參數**的函式，在某些情況下，可能會由於 Go 運行時需要額外的查詢 (runtime lookups) 來區分型別而**導致性能下降**（例如，在 Go 1.20 中，簡單函式呼叫慢了約 30%）。
* **建議：** 撰寫 Go 程式時，應優先選擇可維護性。如果一個函式只需要依賴介面的抽象，則保留介面參數；不要為了追求假定的性能提升而盲目泛型化。

## 12. 標準函式庫中的泛型 (Adding Generics to the Standard Library)

在 Go 1.21 版本之後，標準函式庫開始廣泛利用泛型，以提供更安全、更簡潔的 API：

* **`slices` 和 `maps` 套件：** 提供了 `Equal`、`Insert`、`Delete` 和 `Clone` 等函式，用於安全且高效地操作切片和映射，取代了開發者以往需要自行編寫的複雜邏輯。
* **`sync` 套件：** 提供了 `sync.OnceValue` 和 `sync.OnceValues` 等泛型函式，用於確保程式碼只執行一次並快取結果。

## 13. 泛型的未來展望 (Future Features Unlocked)

泛型為 Go 未來的語言特性發展奠定了基礎，其中一個備受關注的可能性是**Sum Types (聯合型別)**。透過在介面中使用型別元素（`|`），可以明確指定一個變數可能屬於的有限型別集合。這將有助於解決 JSON 編碼的靈活性問題，並強化 Go 語言中較弱的列舉（Enum）功能。
