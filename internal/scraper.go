package brisboproperty

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type ScrapedHouse struct {
	Id    int
	Price int
}

func ExtractHouses(page int, wg *sync.WaitGroup, houses *[]ScrapedHouse, url string) {
	defer wg.Done()

	resp, err := http.Get(fmt.Sprintf(url, page))
	if err != nil {
		panic(err)
	}

	z := html.NewTokenizer(resp.Body)
	z.Next()

	for z.Next() != html.ErrorToken {
		tagName, _ := z.TagName()
		// First find p tags since this is where price is found
		if string(tagName) != "p" {
			continue
		}

		var price int
		var atrKey, atrVal []byte
		for hasMore := true; hasMore; atrKey, atrVal, hasMore = z.TagAttr() {
			if string(atrKey) == "data-testid" && string(atrVal) == "listing-card-price" {
				z.Next()
				price, _ = strconv.Atoi(strings.ReplaceAll(string(z.Text())[1:], ",", ""))
				break
			}
		}
		// Only continue if price was extracted
		if price == 0 {
			continue
		}

		for {
			z.Next()
			tagName, _ := z.TagName()
			if string(tagName) == "a" {
				for {
					atrKey, atrVal, _ = z.TagAttr()
					if string(atrKey) == "href" {
						break
					}
				}
				atrValSlice := strings.Split(string(atrVal), "-")
				listing_id, _ := strconv.Atoi(atrValSlice[len(atrValSlice)-1])
				*houses = append(*houses, ScrapedHouse{listing_id, price})
				break
			}
		}
	}
}
