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
	for i := 48; i < 51; i++ {
		wg.Add(1)
		go brisboproperty.ExtractHouses(i, wg, &houses)
	}
	wg.Wait()

	fmt.Println(len(houses))

	apiKeyBytes, err := os.ReadFile("apiKey.txt")
	if err != nil {
		fmt.Print(err)
	}
	apiKey := string(apiKeyBytes)[:len(apiKeyBytes)-1]

	brisboproperty.GetNewHouses(houses, apiKey, "data.csv")
}
