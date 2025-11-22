package main

import "fmt"

type Printable interface {
	fmt.Stringer
	~int | ~float64
}

type PrintInteger int

// Implementing fmt.Stringer for PrintInt
func (pi PrintInteger) String() string {
	return fmt.Sprintf("PrintInt value: %d", pi)
}

type PrintFloat float64

// Implementing fmt.Stringer for PrintFloat
func (pf PrintFloat) String() string {
	return fmt.Sprintf("PrintFloat value: %.2f", pf)
}

func PrintNumber[T Printable](p T) {
	fmt.Println(p.String())
}

func main() {
	var i PrintInteger = 20
	PrintNumber(i)

	var f PrintFloat = 10.23
	PrintNumber(f)
}
