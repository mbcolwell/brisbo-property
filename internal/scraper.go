package brisboproperty

import (
	"embed"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

//go:embed listing_suburbs.txt
var f embed.FS

type domainId struct {
	Id int
}

type scrapedHouse struct {
	domainId
	Price int
}

func extractSaleIds() []scrapedHouse {
	url := "https://www.domain.com.au/sold-listings/brisbane-region-qld/house/?excludepricewithheld=1&page=%d"
	scrapedHouses := []scrapedHouse{}
	wg := new(sync.WaitGroup)

	for i := 1; i <= 50; i++ {
		wg.Add(1)

		go func(page int, wg *sync.WaitGroup, houses *[]scrapedHouse) {
			defer wg.Done()

			resp, err := http.Get(fmt.Sprintf(url, page))
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

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
						*houses = append(*houses, scrapedHouse{domainId{listing_id}, price})
						break
					}
				}
			}
		}(i, wg, &scrapedHouses)
	}

	wg.Wait()
	return scrapedHouses
}

func readSuburbsFile() string {
	file, err := f.ReadFile("listing_suburbs.txt")
	if err != nil {
		panic(err)
	}

	content := strings.TrimSuffix(string(file), "\n")
	lines := strings.Split(content, "\n")

	return strings.Join(lines, ",")
}

func extractListingIds() []domainId {
	bedroomMin := 3
	bedroomMax := 4
	sizeMin := 400
	sizeMax := 800
	priceMin := 800000
	priceMax := 1600000
	suburbs := readSuburbsFile()

	url := fmt.Sprintf(
		"https://www.domain.com.au/sale/?suburb=%s&ptype=house&bedrooms=%d-%d&landsize=%d-%d&landsizeunit=m2&price=%d-%d&excludeunderoffer=1&ssubs=0&page=%%d",
		suburbs, bedroomMin, bedroomMax, sizeMin, sizeMax, priceMin, priceMax,
	)

	scrapedIds := []domainId{}
	wg := new(sync.WaitGroup)
	for i := 1; i <= 20; i++ {
		wg.Add(1)

		go func(page int, wg *sync.WaitGroup, ids *[]domainId) {
			defer wg.Done()

			resp, err := http.Get(fmt.Sprintf(url, page))
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			z := html.NewTokenizer(resp.Body)
			z.Next()

			for z.Next() != html.ErrorToken {
				tagName, hasAttr := z.TagName()
				if string(tagName) != "li" || !hasAttr {
					continue
				}

				var atr string
				for {
					atrKey, atrVal, hasMore := z.TagAttr()
					atr = string(atrVal)
					if string(atrKey) == "data-testid" && strings.HasPrefix(atr, "listing") {
						listing_id, _ := strconv.Atoi(strings.Split(atr, "-")[1])
						scrapedIds = append(scrapedIds, domainId{listing_id})
						break
					}
					if !hasMore {
						break
					}
				}
			}
		}(i, wg, &scrapedIds)
	}
	wg.Wait()

	return scrapedIds
}
