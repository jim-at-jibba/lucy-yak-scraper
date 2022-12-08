package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly"
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// https://lucyandyak.com/collections/sleepwear
func Scrape() {

	fmt.Println("We are gonna scrape")
	c := colly.NewCollector(
		colly.AllowedDomains("lucyandyak.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	c.OnHTML("#product-grid", func(e *colly.HTMLElement) {
		var list []string
		var currentPJCount = 6 // as of 07/12/22
		fmt.Println("Product list found")
		e.ForEach("li", func(_ int, h *colly.HTMLElement) {
			list = append(list, h.ChildText(".card__heading"))
		})

		fmt.Printf("Count %d", len(list))
		if len(list) > currentPJCount {
			SendMsg("New PJs available")
		}
	})

	c.Visit("https://lucyandyak.com/collections/sleepwear")
}

func SendMsg(msg string) {
	client := twilio.NewRestClient()

	params := &openapi.CreateMessageParams{}
	params.SetTo(os.Getenv("TO_PHONE_NUMBER"))
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(msg)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("SMS sent successfully!")
	}
}

func runCronJob() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Day().At("10:30").Do(func() {
		Scrape()
	})

	s.StartBlocking()
}

func main() {
	fmt.Println("Scraping")
	SendMsg("Starting the scraper")
	runCronJob()
}
