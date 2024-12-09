package fetchers

import (
	"fmt"
	"context"
	"log"
	"time"

	"davidhampgonsalves/lifedashboard/pkg/event"

	"github.com/chromedp/chromedp"
)

func NsPower() ([]event.Event, error) {
	// ctx, cancel := chromedp.NewContext(
	// 	context.Background(),
	// 	chromedp.WithDebugf(log.Printf),
	// )
	// defer cancel()
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("mute-audio", true),
)
ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
// create context
ctx, cancel := chromedp.NewContext(
		ctx,
		// chromedp.WithDebugf(log.Printf),
)
defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 1130*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var usage *string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://myaccount.nspower.ca/`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`body > footer`),
		chromedp.SendKeys(`//input[@name="emailid"]`, "davidhampgonsalves@gmail.com"),
		chromedp.SendKeys(`input#loginradius-login-password`, "Panther1"),
		chromedp.Click(`input#loginradius-submit-login`),

		chromedp.WaitVisible(`h2#ProfileName`),
		chromedp.Click(`#redirectBidgely`, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(3 * time.Second)

	err = chromedp.Run(ctx,
		chromedp.Navigate(`https://nsp.bidgely.com/dashboard/gb-download`),
		chromedp.WaitVisible(`h1.title`),
		chromedp.Click(`(//input[@type=radio])[2]`),
		chromedp.WaitVisible(`//input[@type=text][disabled=false]`),
		chromedp.SendKeys(`(//input[@type=text])[1]`, "06/17/2024"),
		chromedp.SendKeys(`(//input[@type=text])[2]`, "06/17/2024"),
		chromedp.Click(`(//button)[3]`),
		// chromedp.Text(`(//*//ul[contains(@class, "repo-list")]/li[1]//p)[1]`, usage),
	)

	if err != nil {
		log.Fatal(err)
	}

	surf := event.Event{Text: fmt.Sprintf("%s", usage)}
	return []event.Event{surf}, nil
}