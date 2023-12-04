package main

import (
	"fmt"
)

// user-defined data here
type Employee struct {
	Name string
	Age  int
}

func newEmployee(name string, age int) Employee {
	return Employee{name, age}
}

func main() {
	table  := Table{inmem: make(map[int]*Entry)}
	albert := BasicEntry(newEmployee("albert einstein", 144))

    table.offsets = CalculateOffsets(Employee{})

	table.Add(albert)
	table.Add(albert.Copy(nil))
	table.Add(albert.Copy(nil))

	table.Remove(albert, true)

	query, _ := table.GetWhere("Name", "albert einstein")

	for _, v := range query {
		fmt.Println(*v)
	}
}
