package main

import (
	"fmt"
	"log"

	"github.com/yasukotelin/rlpath"
)

func main() {
	scanner := rlpath.Scanner{
		Prompt:  "$ ",
		RootDir: "~/go",
		OnlyDir: false,
	}

	path, err := scanner.Scan()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(path)
}
