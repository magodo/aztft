package main

import (
	"fmt"
	"log"
	"os"

	"github.com/magodo/aztft/aztft"
)

const usage = `Usage: aztft <ID>`

func main() {
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}
	rts, err := aztft.Resolve(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	switch len(rts) {
	case 0:
		fmt.Println("No match")
	case 1:
		fmt.Println(rts[0])
	default:
		fmt.Println(rts)
	}
}
