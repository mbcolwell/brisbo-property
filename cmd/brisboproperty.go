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

	extract_url := "https://www.domain.com.au/sold-listings/brisbane-region-qld/house/?excludepricewithheld=1&page=%d"
	houses := []brisboproperty.ScrapedHouse{}
	wg := new(sync.WaitGroup)
	for i := 1; i <= 50; i++ {
		wg.Add(1)
		go brisboproperty.ExtractHouses(i, wg, &houses, extract_url)
	}
	wg.Wait()

	apiKeyBytes, err := os.ReadFile(apiKeyPath)
	if err != nil {
		fmt.Print(err)
	}
	apiKey := string(apiKeyBytes)[:len(apiKeyBytes)-1]

	pageSize := 10
	paginatedHouses := [][]brisboproperty.ScrapedHouse{}
	for i := len(houses); i > 0; i = i - pageSize {
		paginatedHouses = append(paginatedHouses, houses[max(0, i-pageSize):i])
	}

	for _, housePage := range paginatedHouses {
		brisboproperty.GetNewHouses(housePage, apiKey, dataCsvPath)
	}
}
