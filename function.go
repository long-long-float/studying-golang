package main

import (
	"fmt"
)

func fact(n int) int {
	if n <= 1 {
		return 1
	} else {
		return fact(n-1) * n
	}
}

func twiceAndTriple(n int) (int, int) {
	// 複数の値を返せる
	return n * 2, n * 3
}

func swap(x, y *int) {
	// ポインタもおｋ
	*x, *y = *y, *x
}

func main() {
	fmt.Println(fact(5))

	a, b := twiceAndTriple(3)
	fmt.Println(a, b)

	c, _ := twiceAndTriple(2)
	fmt.Println(c)

	x, y := 1, 2
	swap(&x, &y)
	fmt.Println(x, y)
}
