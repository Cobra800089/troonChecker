package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gtuk/discordwebhook"
)

var (
	cdnUrl               = "https://cdn5.editmysite.com/app/store/api/v28/editor/users/131270493/sites/827516815791883917/products"
	username             = "TroonBot"
	discordErrorUser     = ""
	discordWebhookURL    = ""
	discordListingRoleId = ""
	discordSaleRoleId    = ""
)

type troonData struct {
	Data []struct {
		Name             string    `json:"name"`
		AbsoluteSiteLink string    `json:"absolute_site_link"`
		UpdatedDate      time.Time `json:"updated_date"`
		OnSale           bool      `json:"on_sale"`
	} `json:"data"`
}

func sendDiscord(content string) {
	msg := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}
	if err := discordwebhook.SendMessage(discordWebhookURL, msg); err != nil {
		log.Printf("Error sending message: %v\n", err)
	}
}

func main() {
	var previousBeers []string
	var previousBeersURL []string
	var beerUrl = ""
	var startup = 1

	beerClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 10 seconds
	}

	req, err := http.NewRequest(http.MethodGet, cdnUrl, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("User-Agent", "troonChecker")

	sendDiscord("Troon bot is online")

	// loop forever while doing a check every minute
	for {
		timer1 := time.NewTimer(time.Second * 60)

		// recover unexpected panics per-iteration
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("panic during request loop:", r)
					sendDiscord("<@" + discordErrorUser + "> Error: panic during request loop: " + fmt.Sprint(r))
				}
			}()

			res, getErr := beerClient.Do(req)
			if getErr != nil {
				log.Println(getErr)
				sendDiscord("<@" + discordErrorUser + "> Error: " + getErr.Error())
				return
			}

			if res.Body != nil {
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						log.Println(err)
					}
				}(res.Body)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.Println(err)
				return
			}

			beerList := troonData{}
			jsonErr := json.Unmarshal(body, &beerList)
			if jsonErr != nil {
				log.Println(jsonErr)
				sendDiscord("<@" + discordErrorUser + "> Error: " + jsonErr.Error() + "\\n\\nJSON payload contents:\\n" + string(body[:]))
			}
			//check to see if there is a beer for sale
			if len(beerList.Data) > 0 {
				//loop through all beers in the list
				for i := 0; i < len(beerList.Data); i++ {
					//check to make sure we aren't alerting for the same beer
					if !slices.Contains(previousBeers, beerList.Data[i].Name) {
						beerUrl = beerList.Data[i].AbsoluteSiteLink
						// Don't send a discord alert if the bot is starting up.
						if startup == 0 {
							sendDiscord("<@&" + discordListingRoleId + "> " + beerList.Data[i].Name + " was just listed. (For sale probably later today.)")
						}
						previousBeers = append(previousBeers, beerList.Data[i].Name)
						previousBeersURL = append(previousBeersURL, beerList.Data[i].AbsoluteSiteLink)

					} else if (strings.Contains(previousBeersURL[slices.Index(previousBeers, beerList.Data[i].Name)], "filler")) && (!strings.Contains(beerList.Data[i].AbsoluteSiteLink, "filler")) {
						beerUrl = beerList.Data[i].AbsoluteSiteLink
						previousBeersURL[slices.Index(previousBeers, beerList.Data[i].Name)] = beerUrl
						sendDiscord("<@&" + discordSaleRoleId + "> " + beerList.Data[i].Name + " is now for sale! " + beerUrl)
					}
				}
			}
		}()
		//start the timer and turn startup off
		<-timer1.C
		startup = 0
	}
}
