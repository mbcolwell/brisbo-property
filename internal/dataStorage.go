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

func getStoredIds(dataFilePath string) []int {
	file, err := os.Open(dataFilePath)
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

func GetNewHouses(houses []ScrapedHouse, apiKey string, dataFilePath string) {

	storedIds := getStoredIds(dataFilePath)

	file, err := os.OpenFile(dataFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("Error while opening the file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	wg := new(sync.WaitGroup)
	for _, house := range houses {
		if slices.Contains(storedIds, house.Id) {
			continue
		}

		wg.Add(1)
		go func(house ScrapedHouse) {
			defer wg.Done()

			resp := getHouseInformation(house.Id, apiKey)

			row := []string{
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
				resp.DateListed,
				resp.DateSold,
				strings.Join(resp.Features, "|"),
			}

			err = writer.Write(row)
			if err != nil {
				log.Println("Error while writing to file ", err)
			}
		}(house)
	}
	wg.Wait()
	writer.Flush()
}
