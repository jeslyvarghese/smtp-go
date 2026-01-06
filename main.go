package main

import (
	"flag"
	"fmt"
)

func main() {
	wordPtr := flag.String("word", "foo", "a word")

	numbrPtr := flag.Int("number", 42, "A number")
	forkPtr := flag.Bool("flag", false, "a bool")

	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()

	fmt.Println("word:", *wordPtr)
	fmt.Println("number:", *numbrPtr)
	fmt.Println("bool:", *forkPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())

	fmt.Printf("Hello world")
}
