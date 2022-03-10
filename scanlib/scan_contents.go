package scanUri

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	tld "github.com/jpillora/go-tld"
	"gopkg.in/ini.v1"
)

func GoogleOCR(url string, file []byte) string {

	defer func() {
		if err := recover(); err != nil {
			log.Print("OCR Error => ", url)
		}
	}()

	// Set your environment variables if you get credentials error or panic:
	// https://cloud.google.com/vision/docs/quickstart-client-libraries#client-libraries-install-go
	ctx := context.Background()

	client, _ := vision.NewImageAnnotatorClient(ctx)
	defer client.Close()

	image, _ := vision.NewImageFromReader(bytes.NewReader(file))

	labels, _ := client.DetectTexts(ctx, image, nil, 10)

	var ocrOutput string
	for _, label := range labels {
		ocrOutput += label.Description + " "
	}
	return ocrOutput
}

func SearchText(screenshotToText, htmlresponse, html_text string, hacked_keywords, nickname_keywords []string) (bool, string, string) {

	screenshotToText = strings.ToLower(screenshotToText)
	htmlresponse = strings.ToLower(htmlresponse)
	html_text = strings.ToLower(html_text)

	var hacked_matched string = "-"
	var nickname_matched string = "-"

	for _, keyword := range hacked_keywords {
		var lowerKeyword string = strings.ToLower(keyword)
		if strings.Contains(screenshotToText, lowerKeyword) || strings.Contains(htmlresponse, lowerKeyword) || strings.Contains(html_text, lowerKeyword) {
			hacked_matched = keyword
			break
		}
	}

	for _, keyword := range nickname_keywords {
		var lowerKeyword string = strings.ToLower(keyword)
		if strings.Contains(screenshotToText, lowerKeyword) || strings.Contains(htmlresponse, lowerKeyword) || strings.Contains(html_text, lowerKeyword) {
			nickname_matched = keyword
			break
		}
	}

	if hacked_matched != "-" {
		return true, hacked_matched, nickname_matched
	}

	return false, hacked_matched, nickname_matched
}

func DetectCountry(url string) string {
	cfg, err := ini.Load("tlds.ini")
	var DomainExtInfo = []string{"", ""}
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
	}
	u, _ := tld.Parse(url)
	tldSlice := strings.Split(u.TLD, ".")
	if len(tldSlice) == 2 {
		DomainExtInfo[0] = cfg.Section(tldSlice[0]).Key("name").String()
		DomainExtInfo[1] = cfg.Section(tldSlice[1]).Key("name").String()
	} else {
		DomainExtInfo[0] = cfg.Section(tldSlice[0]).Key("name").String()
		DomainExtInfo[1] = ""
	}
	if DomainExtInfo[0] == "" && DomainExtInfo[1] == "" {
		return "International"
	} else if DomainExtInfo[1] == "" {
		return DomainExtInfo[0]
	}
	return DomainExtInfo[1] + " " + DomainExtInfo[0]

}
