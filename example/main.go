package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rootVIII/cidrcalc"
)

func main() {
	IPCIDR := flag.String("i", "", "IP address")
	flag.Parse()

	var net = &cidrcalc.Subnet{}
	err := net.Calculate(*IPCIDR)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
