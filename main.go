package main

import (
	"github.com/magodo/aztft/aztft"
	"log"
)

func main() {
	if err := aztft.Run(); err != nil {
		log.Fatal(err)
	}
}
