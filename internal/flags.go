package brisboproperty

import "flag"

type cfg struct {
	apiKeyPath     string
	soldCsvPath    string
	listingCsvPath string
	PullSales      bool
	PullListings   bool
}

var Config cfg

func init() {
	flag.StringVar(&Config.apiKeyPath, "api_key_path", "apiKey.txt", "Location of API key")
	flag.StringVar(&Config.soldCsvPath, "sold_csv_path", "sold.csv", "Location of sold csv")
	flag.StringVar(&Config.listingCsvPath, "listing_csv_path", "listing.csv", "Location of listing csv")
	flag.BoolVar(&Config.PullSales, "pull_sales", false, "Whether to pull sales information")
	flag.BoolVar(&Config.PullListings, "pull_listings", false, "Whether to pull listing information")
}
