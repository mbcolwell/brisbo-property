package brisboproperty

func PullSales() {
	houses := extractSaleIds()

	parseApiKey()

	getNewRecords(houses, sold)
}

func PullListings() {
	listingIds := extractListingIds()

	parseApiKey()

	wrappedIds := []scrapedHouse{}
	for _, i := range listingIds {
		wrappedIds = append(wrappedIds, scrapedHouse{i, -1})
	}

	getNewRecords(wrappedIds, listing)
}
