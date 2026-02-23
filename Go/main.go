package main

import (
	"fmt"
	"os"
)

func main() {
	// RunDevTests()

	if len(os.Args) != 2 {
		fmt.Println("Uso: go run . <porta>")
		return
	}

	port := os.Args[1]

	ts := NewTupleSpace()
	StartServer(port, ts)
}
