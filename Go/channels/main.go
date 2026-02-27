package main

import (
	"fmt"
	"linda/core"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Uso: go run . <porta>")
		return
	}

	port := os.Args[1]

	ts := NewTupleSpace()
	core.StartServer(port, ts)
}
