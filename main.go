package main

import (
	"fmt"

	"github.com/minetest-go/mapparser"
)

func main() {
	mapparser.Parse([]byte{}, 0, nil)
	fmt.Println("Starting")
}
