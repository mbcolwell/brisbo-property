package main

import (
	"flag"

	brisboproperty "github.com/mbcolwell/brisbo-property/internal"
)

func main() {
	flag.Parse()

	if brisboproperty.Config.PullSales {
		brisboproperty.PullSales()
	}
	if brisboproperty.Config.PullListings {
		brisboproperty.PullListings()
	}
}
