# Chapter 7. Types, Methods, and Interfaces

重於 Go 語言在型別系統中的獨特設計，強調如何透過方法、組合和隱式介面來建構可測試且易於維護的程式碼。

## Types in Go (Go 語言中的型別)

Go 是一種靜態型別語言，支援內建型別（predeclared types）和使用者定義型別（user-defined types）。

* **使用者定義型別（User-Defined Types）：** 使用 `type` 關鍵字宣告。除了基於結構體（struct）字面值之外，也可以基於任何基本型別或複合型別字面值來定義具體型別（concrete type）。

    ```go
    type Score int // 基於基本型別 int
    type Converter func(string)Score // 基於函式型別
    type TeamScores map[string]Score // 基於 map 型別
    ```

* **型別作用域：** 型別可以在任何區塊層級宣告，但只能在其作用域內存取。
* **術語區分：** **抽象型別**（Abstract type，如介面）定義了型別應具備的功能，而**具體型別**定義了如何儲存資料並實現這些功能。

## Methods (方法)

Go 支援在使用者定義型別上附加方法（methods）。

* **宣告方式：** 方法宣告與函式（function）宣告相似，但包含一個「接收器規範」（receiver specification），位於 `func` 關鍵字與方法名稱之間。
* **接收器慣例：** 慣用上，接收器名稱應為型別名稱的簡短縮寫（通常是第一個字母），不使用 `this` 或 `self`。

    ```go
    type Person struct {
        FirstName string
        // ...
    }
    // 接收器 p (值接收器)
    func (p Person) String() string { 
        return fmt.Sprintf("%s %s, age %d", p.FirstName, p.LastName, p.Age)
    }
    ```

* **限制：** 方法不能重載（overloaded），且必須在與其關聯型別相同的套件（package）內宣告。Go 鼓勵直接存取結構體欄位，避免撰寫不必要的 Getter 和 Setter 方法，方法應保留給業務邏輯。

> 方法只能在套件區塊層級定義，而函式則可以在任何區塊內定義。

### Pointer Receivers and Value Receivers (指標接收器和值接收器)

選擇接收器型別至關重要，它決定了方法是否可以修改接收器實例。

| 接收器類型 | 宣告範例 | 使用情境 | 備註 |
| :--- | :--- | :--- | :--- |
| **指標接收器(pointer receivers)** | `func (c *Counter) Increment()` | 必須用於修改接收器狀態或需要處理 `nil` 實例。 | Go 是 Call-by-Value (值傳遞) 語言，傳遞指標副本允許修改原始數據。 |
| **值接收器(value receivers)** | `func (c Counter) String()` | 不修改接收器狀態時可選用。 | 如果型別包含任何指標接收器方法，慣用上會建議所有方法都使用指標接收器以保持一致性。 |

Go 支援語法糖：無論呼叫端變數是指標型別還是值型別，Go 編譯器都會自動轉換接收器，以匹配方法宣告的類型。

### Code Your Methods for nil Instances (為 nil 實例編寫方法)

Go 語言的特殊之處在於，可以在 `nil` 實例上呼叫方法，前提是該方法使用**指標接收器**。如果使用值接收器在 `nil` 指標實例上呼叫，將會導致執行時 Panic。

```go
 type IntTree struct {
    val         int
    left, right *IntTree
} 
 
func (it *IntTree) Insert(val int) *IntTree {
    if it == nil {
        return &IntTree{val: val}
    }
    if val < it.val {
        it.left = it.left.Insert(val)
    } else if val > it.val {
        it.right = it.right.Insert(val)
    }
    return it
} 
 
func (it *IntTree) Contains(val int) bool {
    switch {
    case it == nil:
        return false
    case val < it.val:
        return it.left.Contains(val)
    case val > it.val:
        return it.right.Contains(val)
    default:
        return true
    }
}

 func main() {
    var it *IntTree
    it = it.Insert(5)
    it = it.Insert(3)
    it = it.Insert(10)
    it = it.Insert(2)
    fmt.Println(it.Contains(2))  // true
    fmt.Println(it.Contains(12)) // false
}
```

指標接收者(Pointer receivers)的工作方式就像指標函式參數；方法中傳入的是指標的拷貝。

### Methods Are Functions Too

每當有函數類型的變數或參數時，你可以將方法作為函數的替代品來使用。

```go
type Adder struct {
 start int
} 
func (a Adder) AddTo(val int) int {
    return a.start + val
}
// method value
f1 := myAdder.AddTo
fmt.Println(f1(10)) 
// method expression
// 第一個參數是該方法的接收者；函數簽名為 func(Adder, int) int
f2 := Adder.AddTo
fmt.Println(f2(myAdder, 15))
```

函式（Functions）與方法（Methods）的劃分不僅是語法上的差異，更是 Go 程式設計哲學——追求清晰度、可維護性以及對狀態（State）管理的一種體現。

##### 核心區分原則：對資料的依賴性與狀態管理

選擇使用函式還是方法，其核心決策點在於**該邏輯是否依賴於需要被修改的資料或結構化的狀態**。

1. **使用方法（Methods）**：
    當的邏輯依賴於**儲存在結構（Struct）中的值**，特別是當這些值在程式執行期間會被配置或修改時，則應當將該邏輯實作成方法。
    * **封裝狀態**：Go 語言建議將需要變動的數值狀態（Mutable State）儲存在 `struct` 之中，以避免使用變動的套件層級變數（Package-level variables），因為後者會使資料流分析複雜化，不利於追蹤變動。
    * **業務邏輯實踐**：方法應該用於實現業務邏輯，而不是簡單的資料存取（例如 getter 和 setter），除非是為了滿足介面（Interface）要求或進行非直接賦值的複雜更新操作。

2. **使用函式（Functions）**：
   當邏輯**僅依賴於輸入參數**，並且不涉及或不需修改任何持久化的結構狀態時，則應使用函式。

##### 總結與設計考量

| 特性 | 函式 (Function) | 方法 (Method) | 依循的設計原則 |
| :--- | :--- | :--- | :--- |
| **主要功能** | 處理僅依賴輸入參數的通用邏輯。 | 處理與結構狀態相關的業務邏輯。 | 清晰的資料流。 |
| **狀態依賴** | 獨立於結構實例的狀態，通常是「純淨」的邏輯。 | 依賴於結構實例的欄位狀態。 | 狀態應結構化管理。 |
| **變動性標示** | 參數的變動性透過指標類型傳遞（`*T`）或複合類型行為體現。 | 透過**接收者類型**（值 vs. 指標）明確標示是否允許變動。 | 依賴值傳遞，指標標註變動。 |
| **可互用性** | 函式即值（Values），可賦予變數或作為參數傳遞。 | 方法可以轉化為函式值（Method Value）或函式表達式（Method Expression），具備函式的特性。| 函式是一級公民（First-Class Citizens）。 |

在進行程式設計時，我們應當遵循此原則：**當邏輯需要變動或存取其所屬結構的內部狀態時，務必使用方法；當邏輯是獨立且通用的，應使用函式**。這不僅提高了程式碼的可讀性，也使得介面（Interfaces）的設計更加自然和解耦。

### Type Declarations Aren’t Inheritance (型別宣告不是繼承)

雖然可以基於另一個使用者定義型別宣告新類型（例如 `type Employee Person`），但這並不構成繼承。

```go
 type HighScore Score
 type Employee Person
```

* **型別獨立性：** 兩個型別（即使底層型別相同）之間沒有層次結構。不能將 `HighScore` 實例賦值給 `Score` 變數，反之亦然，除非進行明確的型別轉換。
* **方法集獨立性：** 定義在原類型上的方法不會自動定義在新類型上。
* **型別作為文件：** 宣告自定義型別的主要價值在於作為程式碼文件，提供概念名稱，使程式碼意圖更清晰。

## iota Is for Enumerations—Sometimes (iota 用於列舉—有時適用)

Go 不具備傳統的列舉（enumeration）型別。它使用 `iota` 常數生成器來為一組常數賦予遞增值。

* **慣用寫法：** 通常先定義一個基於 `int` 的型別，然後在 `const` 區塊中使用 `iota`。

    ```go
    type MailCategory int
    const (
        Uncategorized MailCategory = iota // 0
        Personal                          // 1
        Spam                              // 2
    )
    ```

上述範例，每一個後續的行既沒有指定類型，也沒有賦值。當 Go 編譯器看到這種情況時，它會將類型和賦值重複給區塊中的所有後續常數，也就是 `iota`。`iota` 的值會在每個 const 區塊中定義的常數上遞增，從 0 開始。這表示第一個常數（Uncategorized）被賦值為 0，第二個常數（Personal）被賦值為 1，依此類推。當創建新的 const 區塊時，`iota` 會重新設置為 0。

##### `iota` 的核心機制與語法（Core Mechanism and Syntax）

`iota` 只能在 `const` 區塊內使用，它提供了一種自動遞增的機制，主要用於定義一組相關的常數值：

1. **基本原理**：`iota` 的值在每個 `const` 區塊開始時被重設為 0。它會針對該區塊中的每一行常數宣告遞增 1。
2. **型別定義**：慣例上，在使用 `iota` 建立列舉值之前，最佳實踐是先定義一個以 `int` 為基礎的新型別來代表這組有效的數值。

    ```go
    type MailCategory int
    const (
        Uncategorized MailCategory = iota // iota = 0
        Personal                          // iota = 1
        Spam                              // iota = 2
        // ...
    )
    ```

3. **隱式重複**：在 `const` 區塊中，一旦為第一個常數指定了型別和賦值（例如：`MailCategory = iota`），後續未明確指定型別或值的常數，都會自動重複上一個常數的型別和賦值表達式，並將當前的 `iota` 值應用上去。
4. **按行遞增**：`iota` 是根據常數宣告在 `const` 區塊中的**行數**來遞增的。這意味著即使某行沒有使用 `iota`，它也會影響後續行的 `iota` 值。

##### 設計考量與使用限制（Design Considerations and Limitations）

`iota` 雖然提供了方便的列舉模擬，但其本質上的「基於位置」的特性引入了維護上的脆弱性。

1. **專注於名稱而非數值**：
    * `iota`-based 的列舉僅在關心區分一組值，但**不特別關心其底層數值**時才適用。
    * 設計建議是將 `iota` 用於**「內部」目的**，即常數透過名稱而非數值來引用。
    * **避免用於規格定義**：如果常數的值必須符合外部規格或資料庫定義（即數值是固定的），則不應使用 `iota`，而應明確指定常數值。

2. **維護上的脆弱性（Fragility）**：
    * 如果在 `iota` 定義的常數列表的中間插入新的識別符號，這會導致所有後續常數的數值被重新編號。如果這些常數值被外部系統或資料庫依賴，這將導致應用程式以微妙的方式崩潰或產生錯誤。

3. **零值處理（Zero Value Handling）**：
    * 由於 `iota` 從 0 開始編號，該列舉型別的零值（Zero Value）預設會是第一個常數（數值 0）。
    * 如果的列舉中 0 不代表一個合理的預設狀態，一種常見模式是將第一個 `iota` 值（即 0）指定給空白識別符號 `_`，或者指定給一個表示「無效」或「未初始化」的常數。

4. **型別和方法**：
    * 藉由在 `int` 基礎型別上建立使用者定義型別，這些列舉值可以與其他型別（如 `int32` 或 `int64`）保持區分（`int32` 亦可透過 `rune` 型別別名實現）。
    * 這些使用者定義的列舉型別也可以附加方法（Methods），提供與該組值相關的業務邏輯。

`iota` 是 Go 語言用來模擬列舉的慣用機制，它允許常數值自動遞增，避免了重複的數字編號工作。然而，作為一位資深演算法工程師，我們必須意識到其**基於位置的遞增邏輯**使其具有固有的脆弱性。因此，`iota` 應當被保留給那些僅在程式內部使用，且其具體數值在未來變動不會影響外部系統的常數集合。

## Use Embedding for Composition (使用嵌入實現組合)

Go 鼓勵透過**組合**（Composition）而非繼承來重用程式碼，並內建支援**嵌入**（Embedding）和**提升**（Promotion）機制。

* **嵌入欄位：** 在結構體中定義一個欄位但**不給予名稱**，這使得該欄位成為一個嵌入欄位（Embedded Field）。
* **欄位與方法提升：** 嵌入型別的所有欄位和方法都會被「提升」到包含結構體上，可以直接透過包含結構體實例來呼叫。

    ```go
    type Employee struct {
        Name  string
        ID  string
    } 
    func (e Employee) Description() string {
        return fmt.Sprintf("%s (%s)", e.Name, e.ID)
    }
    type Manager struct {
        Employee // 嵌入欄位，無名稱
        Reports []Employee
    }
    m := Manager{
        Employee: Employee{
            Name: "Bob Bobson",
            ID:   "12345",
        },
        Reports: []Employee{},
    }
    m.ID = "12345" // 存取被提升的 Employee 欄位
    ```

## Embedding Is Not Inheritance (嵌入不是繼承)

儘管嵌入提供了代碼重用，但它不是物件導向意義上的繼承。

* **型別賦值：** 上面範例不能將包含結構體（如 `Manager`）實例直接賦值給嵌入型別（如 `Employee`）的變數。
* **非動態分派：** 如果包含結構體和嵌入型別都有同名方法，當嵌入型別的方法在其內部呼叫另一個方法時，**永遠調用嵌入型別自己的方法**，不會被包含結構體的方法覆蓋（無動態分派）。
* **介面實現：** 嵌入型別的方法計入包含結構體的方法集，使其可以實現介面。

## A Quick Lesson on Interfaces (介面速成課)

介面（Interface）是 Go 語言中唯一的抽象型別，也是其設計的核心。

* **介面定義：** 定義了一組必須由具體型別實現的方法集。
* **方法集規則（再次強調）：**
  * **指標實例：** 方法集包含值接收器和指標接收器定義的方法。
  * **值實例：** 方法集只包含值接收器定義的方法。
  * 範例：如果介面方法使用指標接收器定義（如 `Incrementer` 介面），則只有指標實例（`*Counter`）才能實現它；值實例（`Counter`）無法實現該介面。

## Interfaces Are Type-Safe Duck Typing (介面是型別安全的鴨子型別)

Go 介面的關鍵特點是**隱式實作**（Implicit Implementation）。

* **無需宣告：** 具體型別不需要明確宣告它實現了某個介面。只要它的方法集滿足介面定義，它就隱式地實現了該介面。
* **設計哲學：** 這種機制結合了靜態型別的安全性和動態語言的靈活性與解耦性。
* **慣用原則：** 介面應該由**呼叫端（Client Code）**定義，以表達其對所需功能的確切要求（"Program to an interface, not an implementation"）。
* **介面共享：** 標準庫中的介面（如 `io.Reader`）鼓勵使用**裝飾器模式**（Decorator Pattern）來鏈接功能。

```go
type LogicProvider struct {} 
 
func (lp LogicProvider) Process(data string) string {
    // business logic
} 
 
type Logic interface {
    Process(data string) string
} 
 
type Client struct{
    L Logic
} 
 
func(c Client) Program() {
    // get data from somewhere
    c.L.Process(data)
} 
 
main() {
    c := Client{
        L: LogicProvider{},
    }
    c.Program()
}
```

這段 Go 程式碼提供了一個介面 (interface)，但只有呼叫方（Client）知道它；LogicProvider（邏輯提供者）本身並沒有明確聲明它符合(meets)這個介面。（僅憑）這樣（的隱性實作）就足以在未來允許（替換）新的邏輯提供者，同時也提供了一種「可執行的文件」(executable documentation)，以確保任何傳入 Client 的型別（type）都能符合 Client 的需求。

> 面向接口編程，而不是實現編程
> 嵌入不僅僅用於結構體。你也可以在介面中嵌入另一個介面

## Accept Interfaces, Return Structs (接受介面，返回結構體)

這是一條重要的 Go 慣用設計原則。

* **接受介面：** 函式應接受介面作為參數，以確保程式碼的靈活性和低耦合性。
* **返回結構體：** 函式應返回具體型別（struct），因為這使得未來可以在不破壞向後相容性（backward compatibility）的情況下，向該結構體添加新的方法或欄位。
* **性能權衡：** 傳遞介面會增加堆記憶體分配（heap allocation）的開銷，但通常抽象化帶來的可讀性和可維護性更重要。應在性能分析確認是瓶頸後，再考慮重構以減少介面使用。

## Interfaces and nil (介面與 nil)

介面實例的零值是 `nil`。

* **內部實現：** 介面在 Go runtime 內部由包含兩個指標的結構體實現：一個 Type 指標，一個 Value 指標。
* **nil 判斷：** 只有當 Type 和 Value 兩個指標都為 `nil` 時，介面才被視為 `nil`。
* **潛在陷阱：** 如果將一個具體型別的 `nil` **指標**賦值給介面變數，介面變數本身**不是** `nil`，因為 Type 指標指向了具體型別（例如 `*Counter`）。

## Interfaces Are Comparable

##### 介面比較的基礎原則 (Fundamental Principle of Interface Comparison)

Go 語言的設計允許介面實例之間使用相等運算符 `==` 進行比較。然而，介面的可比較性（Comparability）是基於其內部結構的：一個介面實例只有在其潛在型別（Type）和潛在值（Value）都符合特定條件時才被視為相等。

**核心比較規則：**
兩個介面實例 `i1` 和 `i2` 只有在以下兩個條件同時滿足時才被視為相等（`i1 == i2`）：

1. **型別相等 (Type Equality)**：它們所儲存的潛在型別（Concrete Type）必須相等。
2. **值相等 (Value Equality)**：它們所儲存的潛在值（Value）必須相等。

##### 關鍵的運行時風險：不可比較的潛在型別

儘管介面本身是可比較的，但如果介面儲存的**潛在型別**是 Go 語言中不可比較的類型（例如 `slice`、`map`、`channel`、或包含這些欄位的結構體），則嘗試對這兩個介面變數進行比較將會導致運行時恐慌（panic）。

##### 編譯時通過，運行時恐慌 (Compile-time Success, Runtime Panic)

這是使用介面進行比較時最需要警惕的陷阱：

* Go 編譯器允許將不可比較的具體型別（如 `DoubleIntSlice`，一個切片類型）賦予給一個介面變數（如 `Doubler`）。
* 當兩個儲存了該不可比較型別的介面變數進行 `==` 比較時，程式碼可以成功編譯，但會在執行時觸發 `panic: runtime error: comparing uncomparable type...`。
* 在泛型（Generics）的環境中，即使使用了 `comparable` 約束，如果將不可比較的具體型別（如 `slice`）賦值給泛型介面變數，仍然會導致運行時恐慌。

##### 介面作為 Map 鍵的限制 (Interfaces as Map Keys)

* `map` 的鍵（key）必須是可比較的型別。
* 介面型別可以作為 `map` 的鍵。
* 然而，如果向此 `map` 中添加一個鍵值對，而該鍵的潛在型別是不可比較的，則會立即觸發運行時 `panic`。

##### 設計與安全建議 (Design and Safety Recommendations)

由於 Go 語言的介面實作是隱式的，無法透過語法規定一個介面只能由可比較的型別來實作。為了編寫健壯的程式碼，資深工程師應當採取以下預防措施：

1. **避免依賴介面比較**：
    * 應避免在不知道潛在型別為何種具體型別的情況下，使用 `==` 或 `!=` 來比較介面。
    * 在處理介面時，如果不確定其潛在型別是否可比較，可能需要透過反射（Reflection）來檢查其安全性。

2. **型別斷言和型別轉換**：
    * 當需要對介面進行比較或處理時，優先考慮使用型別斷言（Type Assertion）或型別切換（Type Switch）來提取其潛在的具體型別。
    * 一旦提取出具體型別，若該型別是不可比較的（如 `struct` 中包含 `slice`），則應實作一個自定義的比較函式（而非依賴 `==`）來進行邏輯判斷。

3. **遵循慣例**：
    * 在大多數情況下，Go 程式碼中的比較操作應保持清晰與安全。只有當能夠絕對保證介面下的所有潛在型別都是可比較的，或者已準備好處理運行時恐慌時，才應使用介面比較。

## 11. Type Assertions and Type Switches (型別斷言與型別切換)

當處理介面型別的變數時，需要機制來檢查或提取底層的具體型別。

* **型別斷言（Type Assertion）：** 語法為 `i.(T)`。用於檢查介面變數 `i` 是否包含型別 `T` 或實現了另一個介面。
  * **安全操作：** 必須使用 `comma ok` 慣用法來處理型別不匹配的情況，避免執行時 Panic。

    ```go
    i2, ok := i.(int)
    if !ok {
        return fmt.Errorf("unexpected type for %v",i)
    }
    fmt.Println(i2 + 1)
    ```

* **型別切換（Type Switch）：** 語法為 `switch j := i.(type)`。用於處理介面變數可能是多種型別之一的情況。
  * **慣用寫法：** 為了可讀性，慣用上會將被切換的變數賦值給同名變數（這是一種被接受的 Shadowing 慣例）。

    ```go
    func doThings(i any) {
        switch j := i.(type) {
        case nil:
            // i is nil, type of j is any
        case int:
            // j is of type int
        case MyInt:
            // j is of type MyInt
        case io.Reader:
            // j is of type io.Reader
        case string:
        // j is a string
        case 
        bool, rune:
        // i is either a bool or rune, so j is of type any
        default:
        // no idea what i is, so j is of type any
        }
    }
    ```

### 謹慎使用原則 (Use Type Assertions and Type Switches Sparingly)

應盡量避免過度使用型別斷言和切換，因為這會使函式的 API 宣告不準確。它們的主要用途包括：

1. **檢查可選介面（Optional Interfaces）：** 查看底層具體型別是否實現了額外的、非必需的介面，以便進行性能優化（如 `io.Copy` 檢查 `io.WriterTo`）。
2. **API 演進：** 用於檢查 API 隨著時間推移新增的功能（如 `database/sql/driver` 中基於 Context 的新方法）。

## 12. Function Types Are a Bridge to Interfaces (函式型別是通往介面的橋樑)

這是一個極為實用且體現 Go 設計哲學的慣用語。此機制允許簡單的函式（Functions）或閉包（Closures）能夠實作複雜的介面（Interfaces）抽象，從而提高程式碼的靈活性和模組化程度。

* **實現介面：** 這使得函式能夠實現介面。最經典的例子是 `net/http` 套件中的 `http.Handler` 介面，可以透過將一般函式轉換為 `http.HandlerFunc` 來實現。
* **Function/Interface 抉擇：** 如果一個單一函式極可能依賴於許多其他函式或需要配置狀態，則應使用介面參數，並用函式型別來橋接。

##### 核心概念：在函式上定義方法

Go 語言的類型系統設計允許在任何使用者定義的類型上定義方法（Methods），這項規則也適用於基於函式簽名（Function Signature）定義的新類型。

當使用 `type` 關鍵字定義一個新的函式類型時，您可以為這個新類型附加方法。這個新的**函式類型**就充當了將**普通函式**轉換為**介面實例**的「橋樑」：

1. **定義新的函式類型**：例如，定義一個用於 HTTP 處理的函式類型：

    ```go
    type HandlerFunc func(http.ResponseWriter, *http.Request)
    ```

2. **為該類型實作介面方法**：在這個新類型上實作目標介面（例如 `http.Handler`）所需的方法。在實作過程中，該方法會調用其所持有的底層函式值來執行實際的邏輯。

##### 經典應用：HTTP 服務處理（`net/http` Package）

Go 標準函式庫中的 `net/http` 套件正是利用此模式實現解耦和彈性設計的典範：

* **目標介面**：`http.Handler` 介面定義了處理 HTTP 請求的核心方法：

    ```go
    type Handler interface {
        ServeHTTP(http.ResponseWriter, *http.Request)
    }
    ```

* **橋樑類型**：`http.HandlerFunc` 函式類型正是用來連接普通函式和 `http.Handler` 介面的「橋樑」。它實作了 `ServeHTTP` 方法，並在該方法中調用其持有的函式值。

    ```go
    // http.HandlerFunc 實作 ServeHTTP，並呼叫底層函式 f
    func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        f(w, r)
    }
    ```

* **實際效益**：透過類型轉換（Type Conversion），任何符合 `ServeHTTP` 簽名的普通函式（或閉包），都可以被轉換成 `http.HandlerFunc`，從而滿足 `http.Handler` 介面的要求。這使得開發者可以直接使用輕量的閉包作為 Handler，而無需為每個 Handler 邏輯都定義一個完整的 `struct` 並在該 `struct` 上定義方法。

##### 設計決策與慣例（Idiomatic Design）

理解此模式有助於設計清晰的 API：

* **何時應使用函式類型作為參數？**
    當一個函式參數是為了接受一個簡單、無狀態的單一操作時（例如，用於排序的比較函式，如 `sort.Slice`），直接使用函式類型作為參數是清晰且足夠的。
* **何時應使用介面參數並搭配函式類型橋接？**
    當需要一個抽象層，且預期傳入的實作可能涉及狀態、多個相依性，或者作為更複雜的處理鏈（例如 HTTP 處理程序）的入口點時，應當指定**介面參數**。透過定義一個對應的**函式類型**來實作此介面，這使得簡單的函式也能滿足這個抽象需求，提供了最大的實作彈性。

總而言之，函式類型作為介面橋樑的設計，是 Go 語言在保持靜態類型安全（Type Safety）的同時，實現高度解耦和彈性化程式設計的關鍵機制。

## Implicit Interfaces Make Dependency Injection Easier (隱式介面使依賴注入更容易)

Go 的隱式介面是實現解耦（decoupling）和**依賴注入**（Dependency Injection, DI）的理想工具。

* **DI 核心：** 透過將功能需求抽象為介面，並將這些依賴作為參數傳入結構體（通常在工廠函式中），而不是讓具體邏輯直接依賴於具體實作。
* **範例結構：** 工廠函式接受介面，返回具體結構體。

  ```go
  func NewSimpleLogic(l Logger, ds DataStore) SimpleLogic {
    return SimpleLogic{
        l:    l,
        ds: ds,
    }
  }
  ```

* **結果：** 只有 `main` 函式知道所有的具體型別實作，極大地限制了程式碼中耦合的範圍，使程式碼更易於維護和測試。
