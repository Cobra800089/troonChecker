package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/gtuk/discordwebhook"
	"slices"
	"strings"
)

var cdnUrl = "https://cdn5.editmysite.com/app/store/api/v28/editor/users/131270493/sites/827516815791883917/products"
var username = "TroonBot"
var discordWebhookURL = ""
var discord_listing_role_id = ""
var discord_sale_role_id = ""

type troonData struct {
	Data []struct {
		ID                          string        `json:"id"`
		OwnerID                     string        `json:"owner_id"`
		MerchantID                  string        `json:"merchant_id"`
		SiteID                      string        `json:"site_id"`
		ChannelID                   string        `json:"channel_id"`
		SiteProductID               string        `json:"site_product_id"`
		SquareID                    string        `json:"square_id"`
		Visibility                  string        `json:"visibility"`
		VisibilityTl                string        `json:"visibility_tl"`
		Name                        string        `json:"name"`
		ShortDescription            string        `json:"short_description"`
		PopularityScore             interface{}   `json:"popularity_score"`
		CrossSellProductIds         []interface{} `json:"cross_sell_product_ids"`
		VariationType               string        `json:"variation_type"`
		ProductType                 string        `json:"product_type"`
		ProductTypeDetails          []interface{} `json:"product_type_details"`
		Taxable                     interface{}   `json:"taxable"`
		RequiredFeatureToSell       interface{}   `json:"required_feature_to_sell"`
		MinPrepTime                 interface{}   `json:"min_prep_time"`
		MinPrepTimeDurationIso8601  interface{}   `json:"min_prep_time_duration_iso8601"`
		AdvancedNoticeDuration      interface{}   `json:"advanced_notice_duration"`
		MeasurementUnitAbbreviation interface{}   `json:"measurement_unit_abbreviation"`
		MeasurementUnitPrecision    interface{}   `json:"measurement_unit_precision"`
		HasMultipleMeasurementUnits bool          `json:"has_multiple_measurement_units"`
		IsAlcoholic                 bool          `json:"is_alcoholic"`
		PerOrderMax                 int           `json:"per_order_max"`
		AllowOrderItemQuantities    bool          `json:"allow_order_item_quantities"`
		Badges                      struct {
			LowStock   bool `json:"low_stock"`
			OutOfStock bool `json:"out_of_stock"`
			OnSale     bool `json:"on_sale"`
		} `json:"badges"`
		OnlySubscribable   bool        `json:"only_subscribable"`
		Fulfillable        bool        `json:"fulfillable"`
		SiteLink           string      `json:"site_link"`
		AbsoluteSiteLink   string      `json:"absolute_site_link"`
		Permalink          interface{} `json:"permalink"`
		SeoProductImageID  interface{} `json:"seo_product_image_id"`
		SeoPageDescription interface{} `json:"seo_page_description"`
		SeoPageTitle       interface{} `json:"seo_page_title"`
		OgTitle            interface{} `json:"og_title"`
		OgDescription      interface{} `json:"og_description"`
		AvgRating          interface{} `json:"avg_rating"`
		AvgRatingAll       interface{} `json:"avg_rating_all"`
		Inventory          struct {
			Total                               int         `json:"total"`
			Lowest                              interface{} `json:"lowest"`
			Enabled                             bool        `json:"enabled"`
			AllVariationsSoldOut                bool        `json:"all_variations_sold_out"`
			SomeVariationsSoldOut               bool        `json:"some_variations_sold_out"`
			MarkedSoldOutAtAllExistingLocations bool        `json:"marked_sold_out_at_all_existing_locations"`
			MarkedSoldOutSkusCount              int         `json:"marked_sold_out_skus_count"`
			AllInventoryTotal                   int         `json:"all_inventory_total"`
			HasLocationNotTracking              bool        `json:"has_location_not_tracking"`
		} `json:"inventory"`
		Price struct {
			High                                  int    `json:"high"`
			HighWithModifiers                     int    `json:"high_with_modifiers"`
			HighWithSubscriptions                 int    `json:"high_with_subscriptions"`
			HighFormatted                         string `json:"high_formatted"`
			HighFormattedWithModifiers            string `json:"high_formatted_with_modifiers"`
			HighFormattedWithSubscriptions        string `json:"high_formatted_with_subscriptions"`
			HighSubunits                          int    `json:"high_subunits"`
			HighWithSubscriptionsSubunits         int    `json:"high_with_subscriptions_subunits"`
			HighWithModifiersSubunits             int    `json:"high_with_modifiers_subunits"`
			Low                                   int    `json:"low"`
			LowWithModifiers                      int    `json:"low_with_modifiers"`
			LowWithSubscriptions                  int    `json:"low_with_subscriptions"`
			LowFormatted                          string `json:"low_formatted"`
			LowFormattedWithModifiers             string `json:"low_formatted_with_modifiers"`
			LowFormattedWithSubscriptions         string `json:"low_formatted_with_subscriptions"`
			LowSubunits                           int    `json:"low_subunits"`
			LowWithSubscriptionsSubunits          int    `json:"low_with_subscriptions_subunits"`
			LowWithModifiersSubunits              int    `json:"low_with_modifiers_subunits"`
			RegularHigh                           int    `json:"regular_high"`
			RegularHighWithModifiers              int    `json:"regular_high_with_modifiers"`
			RegularHighWithSubscriptions          int    `json:"regular_high_with_subscriptions"`
			RegularHighFormatted                  string `json:"regular_high_formatted"`
			RegularHighFormattedWithModifiers     string `json:"regular_high_formatted_with_modifiers"`
			RegularHighFormattedWithSubscriptions string `json:"regular_high_formatted_with_subscriptions"`
			RegularHighSubunits                   int    `json:"regular_high_subunits"`
			RegularHighWithSubscriptionsSubunits  int    `json:"regular_high_with_subscriptions_subunits"`
			RegularHighWithModifiersSubunits      int    `json:"regular_high_with_modifiers_subunits"`
			RegularLow                            int    `json:"regular_low"`
			RegularLowWithModifiers               int    `json:"regular_low_with_modifiers"`
			RegularLowWithSubscriptions           int    `json:"regular_low_with_subscriptions"`
			RegularLowFormatted                   string `json:"regular_low_formatted"`
			RegularLowFormattedWithModifiers      string `json:"regular_low_formatted_with_modifiers"`
			RegularLowFormattedWithSubscriptions  string `json:"regular_low_formatted_with_subscriptions"`
			RegularLowSubunits                    int    `json:"regular_low_subunits"`
			RegularLowWithSubscriptionsSubunits   int    `json:"regular_low_with_subscriptions_subunits"`
			RegularLowWithModifiersSubunits       int    `json:"regular_low_with_modifiers_subunits"`
		} `json:"price"`
		OnSale           bool `json:"on_sale"`
		ComboTypeDetails struct {
			Slots []interface{} `json:"slots"`
		} `json:"combo_type_details"`
		UpdatedDate time.Time `json:"updated_date"`
		CreatedDate time.Time `json:"created_date"`
		IndexedDate time.Time `json:"indexed_date"`
		Preordering struct {
			Pickup   bool `json:"pickup"`
			Delivery bool `json:"delivery"`
			Shipping bool `json:"shipping"`
			DineIn   bool `json:"dine_in"`
			Manual   bool `json:"manual"`
			Download bool `json:"download"`
		} `json:"preordering"`
		FulfillmentAvailability struct {
			Pickup   []interface{} `json:"pickup"`
			Delivery []interface{} `json:"delivery"`
			Shipping []interface{} `json:"shipping"`
			DineIn   []interface{} `json:"dine_in"`
			Manual   []interface{} `json:"manual"`
		} `json:"fulfillment_availability"`
		Thumbnail struct {
			Data interface{} `json:"data"`
		} `json:"thumbnail"`
		PlaceholderImage struct {
			Data struct {
				Placeholder bool   `json:"placeholder"`
				URL         string `json:"url"`
				AbsoluteURL string `json:"absolute_url"`
			} `json:"data"`
		} `json:"placeholder_image"`
	} `json:"data"`
	Meta struct {
		Pagination struct {
			Total       int           `json:"total"`
			Count       int           `json:"count"`
			PerPage     int           `json:"per_page"`
			CurrentPage int           `json:"current_page"`
			TotalPages  int           `json:"total_pages"`
			Links       []interface{} `json:"links"`
		} `json:"pagination"`
	} `json:"meta"`
}

func main() {
	previousBeers := []string{}
	previousBeersURL := []string{}
	var beerUrl = ""
	var content = ""

	beerClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, cdnUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "troonChecker")

	// loop forever while doing a check every minute
	for {
		timer1 := time.NewTimer(time.Second * 60)

		res, getErr := beerClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		beer_list := troonData{}
		jsonErr := json.Unmarshal(body, &beer_list)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		//check to see if there is a beer for sale
		if len(beer_list.Data) > 0 {
			//loop through all beers in the list
			for i := 0; i < len(beer_list.Data); i++ {
				//check to make sure we aren't alerting for the same beer
				if ! slices.Contains(previousBeers, beer_list.Data[i].Name) {
					previousBeers = append(previousBeers, beer_list.Data[i].Name)
					previousBeersURL = append(previousBeersURL, beer_list.Data[i].AbsoluteSiteLink)
					beerUrl = beer_list.Data[i].AbsoluteSiteLink
					content = "<@&" + discord_listing_role_id + "> "+ beer_list.Data[i].Name + " was just listed. (For sale probably later today.)"
					message := discordwebhook.Message{
						Username: &username,
						Content: &content,
					}
				
					err := discordwebhook.SendMessage(discordWebhookURL, message)
					if err != nil {
						log.Fatal(err)
					}
				} else if (strings.Contains(previousBeersURL[slices.Index(previousBeers, beer_list.Data[i].Name)], "filler")) && (! strings.Contains(beer_list.Data[i].AbsoluteSiteLink, "filler")) {
					beerUrl = beer_list.Data[i].AbsoluteSiteLink
					previousBeersURL[slices.Index(previousBeers, beer_list.Data[i].Name)] = beerUrl
					content = "<@&" + discord_sale_role_id + "> "+ beer_list.Data[i].Name + " is now for sale! " + beerUrl
					message := discordwebhook.Message{
						Username: &username,
						Content: &content,
					}
				
					err := discordwebhook.SendMessage(discordWebhookURL, message)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
		//start the timer
		<-timer1.C
	}
}
