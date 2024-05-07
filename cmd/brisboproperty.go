package main

import (
	"fmt"
	"os"
	"sync"

	brisboproperty "github.com/mbcolwell/brisbo-property/internal"
)

func main() {
	houses := []brisboproperty.ScrapedHouse{}
	wg := new(sync.WaitGroup)
	for i := 1; i < 51; i++ {
		wg.Add(1)
		go brisboproperty.ExtractHouses(i, wg, &houses)
	}
	wg.Wait()

	// nCalls := 0
	fmt.Println(len(houses))

	apiKeyBytes, err := os.ReadFile("apiKey.txt")
	if err != nil {
		fmt.Print(err)
	}
	apiKey := string(apiKeyBytes)[:len(apiKeyBytes)-1]

	fmt.Println(brisboproperty.GetHouse(2019028686, apiKey))
}
