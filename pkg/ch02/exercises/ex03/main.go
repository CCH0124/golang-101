package main

import "fmt"

type Employee struct {
	firstName string
	lastName  string
	id        int
}

func main() {
	emp1 := Employee{
		"Itachi",
		"Test",
		12345,
	}
	emp2 := Employee{
		firstName: "Naruto",
		lastName:  "Test",
		id:        12346,
	}

	var emp3 Employee
	emp3.firstName = "Madara"
	emp3.lastName = "Test"
	emp3.id = 12347

	fmt.Println(emp1)
	fmt.Println(emp2)
	fmt.Println(emp3)
}
