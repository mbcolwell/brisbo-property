package brisboproperty

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type propertyType int

const (
	sold propertyType = iota
	listing
)

func getStoredIds(filepath string) []int {
	file, err := os.Open(filepath)
	if err != nil {
		log.Panic("Error while reading the file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		fmt.Printf("Error reading records: %v\n", err)
	}

	ids := []int{}
	for r, record := range records[1:] {
		i, err := strconv.Atoi(record[0])
		if err != nil {
			fmt.Printf("Error reading record %d\n", r)
			continue
		}
		ids = append(ids, i)
	}

	return ids
}

func getNewRecords(houses []scrapedHouse, pType propertyType) {

	var filepath string
	switch pType {
	case sold:
		filepath = Config.soldCsvPath
	case listing:
		filepath = Config.listingCsvPath
	}

	storedIds := getStoredIds(filepath)

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("Error while opening the file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	housePageSize := 10
	for i := len(houses); i > 0; i = i - housePageSize { // Paginate into max of 10 concurrent requests at once
		wg := new(sync.WaitGroup)
		for _, house := range houses[max(0, i-housePageSize):i] {
			if slices.Contains(storedIds, house.Id) {
				continue
			}

			wg.Add(1)
			go func(house scrapedHouse) {
				defer wg.Done()

				resp := getHouseInformation(house.Id)

				var row []string
				if house.Price == -1 {
					var photo string
					if len(resp.Media) > 0 {
						photo = resp.Media[0].Url
					} else {
						photo = ""
					}
					row = []string{
						fmt.Sprintf("www.domain.com.au/%d", house.Id),
						resp.AddressParts.DisplayAddress,
						resp.AddressParts.Suburb,
						resp.DateListed,
						strconv.Itoa(resp.Area),
						strconv.Itoa(resp.Beds),
						strconv.Itoa(resp.Baths),
						strings.Join(resp.Features, "|"),
						photo,
						strconv.FormatFloat(resp.Geolocation.Longitude, 'f', -1, 64),
						strconv.FormatFloat(resp.Geolocation.Latitude, 'f', -1, 64),
						strconv.Itoa(resp.Cars),
						resp.DateUpdated,
					}
				} else {
					row = []string{
						strconv.Itoa(house.Id),
						strconv.Itoa(house.Price),
						resp.AddressParts.DisplayAddress,
						resp.AddressParts.Suburb,
						strconv.FormatFloat(resp.Geolocation.Longitude, 'f', -1, 64),
						strconv.FormatFloat(resp.Geolocation.Latitude, 'f', -1, 64),
						strconv.Itoa(resp.Area),
						strconv.Itoa(resp.Beds),
						strconv.Itoa(resp.Baths),
						strconv.Itoa(resp.Cars),
						resp.DateUpdated, // For sales, this is the listed date
						resp.DateListed,  // For sales, this is the sold date
						strings.Join(resp.Features, "|"),
					}
				}

				err = writer.Write(row)
				if err != nil {
					log.Println("Error while writing to file ", err)
				}
			}(house)
		}
		wg.Wait()
	}
	writer.Flush()
}
