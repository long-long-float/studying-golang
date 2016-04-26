package main

import (
	"fmt"
)

type List struct {
	car int
	cdr *List
}

func cons(car int, cdr *List) *List {
	return &List{car, cdr}
}

func null() *List {
	return &List{}
}

func (list *List) isNull() bool {
	return list.cdr == nil
}

func (list *List) toArray() []int {
	values := []int{}
	current := list
	for !current.isNull() {
		values = append(values, current.car)
		current = current.cdr
	}
	return values
}

func list(nums ...int) *List {
	head := null()
	current := head
	for _, n := range nums {
		current.car = n
		current.cdr = &List{}
		current = current.cdr
	}
	return head
}

func (list *List) collect(f func(int) int) *List {
	if list.isNull() {
		return null()
	} else {
		return cons(f(list.car), list.cdr.collect(f))
	}
}

func main() {
	list := list(1, 2, 3)
	fmt.Println(list.collect(func(i int) int { return i * i }).toArray())
}
