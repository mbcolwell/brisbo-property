package brisboproperty

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var apiKey string

func parseApiKey() {
	if apiKey != "" {
		return
	}
	apiKeyBytes, err := os.ReadFile(Config.apiKeyPath)
	if err != nil {
		fmt.Print(err)
	}
	apiKey = string(apiKeyBytes)[:len(apiKeyBytes)-1]
}

func PullSales() {
	extract_url := "https://www.domain.com.au/sold-listings/brisbane-region-qld/house/?excludepricewithheld=1&page=%d"
	houses := []scrapedHouse{}
	wg := new(sync.WaitGroup)
	for i := 1; i <= 50; i++ {
		wg.Add(1)
		go extractHouses(i, wg, &houses, extract_url)
	}
	wg.Wait()

	parseApiKey()

	housePageSize := 10
	for i := len(houses); i > 0; i = i - housePageSize {
		getNewHouses(houses[max(0, i-housePageSize):i])
	}
}

func listingsUrlBuilder(suburb string, bedrooms int) string {
	return fmt.Sprintf("www.domain.com.au/sale/%s/?bedrooms=%d&price=1000000-1600000&excludeunderoffer=1&ssubs=0", suburb, bedrooms)
}
func PullListings() {
	parseApiKey()

	file, err := os.OpenFile(Config.listingCsvPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("Error while opening the file", err)
	}
	defer file.Close()

	// SUBURBS := []string{"kenmore-qld-4069"}

	listingsUrlBuilder("kenmore-qld-4069", 3)
}
