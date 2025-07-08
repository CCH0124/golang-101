# 函式
## 宣告與呼叫函式
Go 語言程式起點是 `main` 函式，它不接收任何參數和回傳值。

對於有回傳和參數的函示如同下面

[sample code](CH5/example1/main.go)
```go
func div(numerator int, denominator int) int {
    if denominator == 0 {
        return 0
    }
    return numerator/denominator
}
```

Go 是有型態的語言，所以要指定參數的回傳型態。如果有定義回傳值則必須使用 `return` 關鍵字。

## 模擬具名（named）與選用（optional）參數
對於 Go 來說沒有具名（named）與選用（optional）參數，如果要模擬可以使用 `struct` 方式。

## 數量不定的輸入參數與 slice
Go 支援*數量不定（variadic）的參數*。數量不定的參數必須是輸入參數的最後一個參數，並在型態前使用 `...`來表示為不定量參數。它最後會是以 `slice` 來建立，操作方式也會是一樣。

[sample code](CH5/example2/main.go)

## 多個回傳值
下面範例使用了 `()` 定義了三個回傳值
```go
func divAndRemainder(numerator int, denominator int) (int, int, error) {
    if denominator == 0 {
        return 0, 0, errors.New("Zero...")
    }
    return numerator/denominator,numerator%denominator, nil
}
```

>如果不需要讀取某個函式的回傳值就使用 `_`

## 函式是值
函式型態包含 `func`、參數的型態與回傳值，這組合是函式的特徵標記。