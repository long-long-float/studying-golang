package main

import (
	"./lisp"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage :", os.Args[0], " FILE")
	}

	dat, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	src := []rune(string(dat))

	exprs, err := lisp.Parse(src)
	if err != nil {
		log.Fatal(err)
	}

	for _, expr := range exprs {
		fmt.Println(expr.Pretty())
	}

	if err := lisp.Evaluate(exprs); err != nil {
		log.Fatal(err)
	}
}
