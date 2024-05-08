package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	brisboproperty "github.com/mbcolwell/brisbo-property/internal"
)

func main() {
	var apiKeyPath, dataCsvPath string
	flag.StringVar(&apiKeyPath, "api_key_path", "apiKey.txt", "Location of API key")
	flag.StringVar(&dataCsvPath, "data_csv_path", "data.csv", "Location of data csv")
	flag.Parse()

	houses := []brisboproperty.ScrapedHouse{}
	wg := new(sync.WaitGroup)
	for i := 20; i < 51; i++ {
		wg.Add(1)
		go brisboproperty.ExtractHouses(i, wg, &houses)
	}
	wg.Wait()

	apiKeyBytes, err := os.ReadFile(apiKeyPath)
	if err != nil {
		fmt.Print(err)
	}
	apiKey := string(apiKeyBytes)[:len(apiKeyBytes)-1]

	nCalls := 0
	brisboproperty.GetNewHouses(houses, apiKey, dataCsvPath, &nCalls)
}
