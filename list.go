package main

import (
	"fmt"
)

type List interface {
	Car() int
	Cdr() List
}

type Cons struct {
	car int
	cdr List
}

func Car(list Cons) int {
	return list.car
}

func Cdr(list Cons) List {
	return list.cdr
}

type Null struct {
}

func Car(list Null) int {
	return 0
}

func Cdr(list Null) List {
	return nil
}

func makeList(nums ...int) List {
	if len(nums) == 0 {
		return &Null{}
	} else {
		head := &Cons{}
		current := head
		for i, n := range nums {
			current.car = n
			if i < len(nums)-1 {
				current.cdr = &Cons{}
			} else {
				current.cdr = &Null{}
			}
		}
		return head
	}
}

func main() {
	list := makeList(1, 2, 3)
	fmt.Println(list)
}
