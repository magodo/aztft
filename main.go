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

  -api		: (Default: false) Allow to use Azure API to disambiguate matching results (e.g. whether a VM is a Linux VM or Windows VM)
  -import	: (Default: false) Print the TF import instruction
`)
}

func main() {
	flagAPI := flag.Bool("api", false, "")
	flagImport := flag.Bool("import", false, "")
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) != 1 {
		usage()
		os.Exit(1)
	}
	var output []string
	if *flagImport {
		rts, specs, err := aztft.QueryImportSpecs(flag.Args()[0], *flagAPI)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(rts); i++ {
			output = append(output, fmt.Sprintf("terraform import %s %s", rts[i], specs[i]))
		}
	} else {
		rts, err := aztft.Query(flag.Args()[0], *flagAPI)
		if err != nil {
			log.Fatal(err)
		}
		output = rts
	}
	switch len(output) {
	case 0:
		fmt.Println("No match")
	case 1:
		fmt.Println(output[0])
	default:
		for _, line := range output {
			fmt.Println(line)
		}
	}
}
