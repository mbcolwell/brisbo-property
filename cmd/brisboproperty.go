package main

import (
	"flag"
	"fmt"

	brisboproperty "github.com/mbcolwell/brisbo-property/internal"
)

func main() {
	var n, m int

	flag.IntVar(&n, "n", 0, "First value")
	flag.IntVar(&m, "m", 0, "First value")
	flag.Parse()

	r := brisboproperty.Add(n, m)
	fmt.Println(r)
}
