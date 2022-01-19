package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var toPhoneNumber = ""
var twilioPhoneNumber = ""
var twilioAuthtoken = ""
var twilioAccountSid = ""

var cdnUrl = "https://cdn5.editmysite.com/app/store/api/v16/editor/users/131270493/sites/827516815791883917/products"

type troonData struct {
	Data []struct {
		ID                    string        `json:"id"`
		OwnerID               string        `json:"owner_id"`
		SiteID                string        `json:"site_id"`
		SiteProductID         string        `json:"site_product_id"`
		Visibility            string        `json:"visibility"`
		VisibilityTl          string        `json:"visibility_tl"`
		VisibilityLocked      bool          `json:"visibility_locked"`
		Name                  string        `json:"name"`
		ShortDescription      interface{}   `json:"short_description"`
		VariationType         string        `json:"variation_type"`
		ProductType           string        `json:"product_type"`
		ProductTypeDetails    []interface{} `json:"product_type_details"`
		Taxable               bool          `json:"taxable"`
		RequiredFeatureToSell interface{}   `json:"required_feature_to_sell"`
		MinPrepTime           interface{}   `json:"min_prep_time"`
		SiteShippingBoxID     interface{}   `json:"site_shipping_box_id"`
		SiteLink              string        `json:"site_link"`
		Permalink             interface{}   `json:"permalink"`
		SeoPageDescription    interface{}   `json:"seo_page_description"`
		SeoPageTitle          interface{}   `json:"seo_page_title"`
		AvgRating             interface{}   `json:"avg_rating"`
		AvgRatingAll          interface{}   `json:"avg_rating_all"`
		Inventory             struct {
			Total                               int  `json:"total"`
			Lowest                              int  `json:"lowest"`
			Enabled                             bool `json:"enabled"`
			MarkedSoldOutAtAllExistingLocations bool `json:"marked_sold_out_at_all_existing_locations"`
			MarkedSoldOutSkusCount              int  `json:"marked_sold_out_skus_count"`
			HasLocationNotTracking              bool `json:"has_location_not_tracking"`
		} `json:"inventory"`
		MeasurementUnitAbbreviation interface{} `json:"measurement_unit_abbreviation"`
		Price                       struct {
			High                              int    `json:"high"`
			HighWithModifiers                 int    `json:"high_with_modifiers"`
			HighFormatted                     string `json:"high_formatted"`
			HighFormattedWithModifiers        string `json:"high_formatted_with_modifiers"`
			HighSubunits                      int    `json:"high_subunits"`
			Low                               int    `json:"low"`
			LowWithModifiers                  int    `json:"low_with_modifiers"`
			LowFormatted                      string `json:"low_formatted"`
			LowFormattedWithModifiers         string `json:"low_formatted_with_modifiers"`
			LowSubunits                       int    `json:"low_subunits"`
			RegularHigh                       int    `json:"regular_high"`
			RegularHighWithModifiers          int    `json:"regular_high_with_modifiers"`
			RegularHighFormatted              string `json:"regular_high_formatted"`
			RegularHighFormattedWithModifiers string `json:"regular_high_formatted_with_modifiers"`
			RegularHighSubunits               int    `json:"regular_high_subunits"`
			RegularLow                        int    `json:"regular_low"`
			RegularLowWithModifiers           int    `json:"regular_low_with_modifiers"`
			RegularLowFormatted               string `json:"regular_low_formatted"`
			RegularLowFormattedWithModifiers  string `json:"regular_low_formatted_with_modifiers"`
			RegularLowSubunits                int    `json:"regular_low_subunits"`
		} `json:"price"`
		OnSale       bool        `json:"on_sale"`
		IsAlcoholic  bool        `json:"is_alcoholic"`
		ImportSource interface{} `json:"import_source"`
		CreatedDate  time.Time   `json:"created_date"`
		UpdatedDate  time.Time   `json:"updated_date"`
		Badges       struct {
			LowStock   bool `json:"low_stock"`
			OutOfStock bool `json:"out_of_stock"`
			OnSale     bool `json:"on_sale"`
		} `json:"badges"`
		SeoProductImageID interface{} `json:"seo_product_image_id"`
		OgTitle           interface{} `json:"og_title"`
		OgDescription     interface{} `json:"og_description"`
		Thumbnail         struct {
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

	var currentBeer = ""

	beerClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, cdnUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "troonChecker")

	// loop forever while doing a check every 5 minutes (300 seconds)
	for {
		timer1 := time.NewTimer(time.Second * 300)

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

		beer1 := troonData{}
		jsonErr := json.Unmarshal(body, &beer1)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		//check to see if there is a beer for sale
		if len(beer1.Data) > 0 {
			//check to make sure we aren't alerting for the same beer
			if beer1.Data[0].Name != currentBeer {
				currentBeer = beer1.Data[0].Name
				client := twilio.NewRestClientWithParams(twilio.RestClientParams{
					Username: twilioAccountSid,
					Password: twilioAuthtoken,
				})

				params := &openapi.CreateMessageParams{}
				params.SetTo(toPhoneNumber)
				params.SetFrom(twilioPhoneNumber)
				params.SetBody("Name: " + beer1.Data[0].Name + "\nPrice: " + beer1.Data[0].Price.LowFormattedWithModifiers)
				//send the sms
				_, err := client.ApiV2010.CreateMessage(params)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("SMS sent successfully!")
				}
			}
		}
		//start the timer
		<-timer1.C
	}
}
