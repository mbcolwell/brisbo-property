package brisboproperty

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type domainGeolocation struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type domainAddress struct {
	Suburb         string `json:"suburb"`
	DisplayAddress string `json:"displayAddress"`
}

type domainApiResponse struct {
	Geolocation  domainGeolocation `json:"geoLocation"`
	Area         int               `json:"landAreaSqm"`
	Features     []string          `json:"features"`
	Baths        int               `json:"bathrooms"`
	Beds         int               `json:"bedrooms"`
	Cars         int               `json:"carspaces"`
	DateListed   string            `json:"dateUpdated"`
	DateSold     string            `json:"dateListed"`
	AddressParts domainAddress     `json:"addressParts"`
}

func retryRequest(nRetries *int, err error) {
	if *nRetries > 0 {
		*nRetries--
	} else {
		fmt.Printf("Out of retries, throwing panic\n")
		log.Panic(err)
	}
}

func requestData(url string) []byte {
	nRETRIES := 4
	client := http.Client{
		Timeout: time.Second * 30, // Timeout after 30 seconds
	}
	var body []byte

	for {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			retryRequest(&nRETRIES, err)
			continue
		}
		res, getErr := client.Do(req)
		if getErr != nil {
			retryRequest(&nRETRIES, getErr)
			continue
		}
		if res.Body != nil {
			defer res.Body.Close()
		}
		bodyIo, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			retryRequest(&nRETRIES, readErr)
			continue
		}
		body = bodyIo
		break
	}

	return []byte(body)
}

func getHouseInformation(houseId int, apiKey string) domainApiResponse {
	url := fmt.Sprintf("https://api.domain.com.au/sandbox/v1/listings/%d?api_key=%s", houseId, apiKey)

	i := domainApiResponse{}
	err := json.Unmarshal(requestData(url), &i)

	if err != nil {
		log.Printf("Bad response for house ID %d", houseId)
	}

	return i
}
