package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/magodo/aztft/aztft"
)

func usage() {
	fmt.Fprint(os.Stderr, `Usage: aztft [option] <ID>

Options:

  -api: (Default: false) Allow to use Azure API to disambiguate matching results (e.g. whether a VM is a Linux VM or Windows VM)
`)
}

func main() {
	flagAPI := flag.Bool("api", false, "")
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) != 1 {
		usage()
		os.Exit(1)
	}
	rts, err := aztft.Query(flag.Args()[0], *flagAPI)
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
