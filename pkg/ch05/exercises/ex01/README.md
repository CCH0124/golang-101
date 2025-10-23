# Exercise 1

## Question

Create a struct named `Person` with three fields: `FirstName` and `LastName` of type `string` and `Age` of type `int`. Write a function called `MakePerson` that takes in `firstName`, `lastName`, and `age` and returns a `Person`. Write a second function `MakePersonPointer` that takes in `firstName`, `lastName`, and `age` and returns a `*Person`. Call both from `main`. Compile your program with `go build -gcflags="-m"`. This both compiles your code and prints out which values escape to the heap. Are you surprised about what escapes?

## Solution

```bash
go build -gcflags="-m" pkg\ch05\exercises\ex01\main.go
# command-line-arguments
pkg\ch05\exercises\ex01/main.go:11:6: can inline MakePerson
pkg\ch05\exercises\ex01/main.go:19:6: can inline MakePersonPointer
pkg\ch05\exercises\ex01/main.go:28:18: inlining call to MakePerson
pkg\ch05\exercises\ex01/main.go:29:13: inlining call to fmt.Println
pkg\ch05\exercises\ex01/main.go:30:25: inlining call to MakePersonPointer
pkg\ch05\exercises\ex01/main.go:31:13: inlining call to fmt.Println
pkg\ch05\exercises\ex01/main.go:11:17: leaking param: firstName to result ~r0 level=0
pkg\ch05\exercises\ex01/main.go:11:28: leaking param: lastName to result ~r0 level=0
pkg\ch05\exercises\ex01/main.go:19:24: leaking param: firstName
pkg\ch05\exercises\ex01/main.go:19:35: leaking param: lastName
pkg\ch05\exercises\ex01/main.go:20:9: &Person{...} escapes to heap
pkg\ch05\exercises\ex01/main.go:29:13: ... argument does not escape
pkg\ch05\exercises\ex01/main.go:29:14: p1 escapes to heap
pkg\ch05\exercises\ex01/main.go:30:25: &Person{...} escapes to heap
pkg\ch05\exercises\ex01/main.go:31:13: ... argument does not escape
```
