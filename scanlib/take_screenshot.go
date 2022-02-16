package scanUri

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func fullScreenshot(urlstr string, quality int, res *[]byte, response *string, source_text *string) chromedp.Tasks {
	return chromedp.Tasks{
		emulation.SetUserAgentOverride("Mozilla/5.0 (Defalyzer) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"),
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
		chromedp.OuterHTML("html", response),
		chromedp.Text(`body`, source_text, chromedp.NodeVisible, chromedp.ByQuery),
	}
}

func GetSources(parentCtx context.Context, url string) ([]byte, string, string) {

	var buf []byte
	var htmlresponse string
	var html_text string

	if isURL(url) {

		ctx, cancel := chromedp.NewContext(parentCtx)
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		if err := chromedp.Run(ctx, fullScreenshot(url, 100, &buf, &htmlresponse, &html_text)); err != nil {
			log.Printf("Could not access this website: " + url)
			// fmt.Println(err)
		}

		return buf, htmlresponse, html_text

	}

	return buf, "", ""

}
