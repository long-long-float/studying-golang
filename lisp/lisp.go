package main

import (
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

	fmt.Println(string(src))
}
