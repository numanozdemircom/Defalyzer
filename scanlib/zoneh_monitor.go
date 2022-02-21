package scanUri

import (
	"encoding/xml"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Rss struct {
	Channel struct {
		Item []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func ExtractZoneH(urlList []string, matchFromFlag string) [][]string {

	//var matchedDefacements map[string]string = map[string]string{}
	var matchedDefacements [][]string
	r, _ := regexp.Compile("^(.*? ){3}")

	resp, err := http.Get("https://zone-h.org/rss/specialdefacements")
	if err != nil {
		log.Println("Could not access Zone-H website")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not access Zone-H website")
	}

	var websitesList Rss
	xml.Unmarshal(body, &websitesList)

	if matchFromFlag == "match_all_urls" {
		for _, v := range websitesList.Channel.Item {

			var attackerNickname string = html.UnescapeString(html.UnescapeString(r.ReplaceAllString(v.Description, "")))
			matchedDefacements = append(matchedDefacements, []string{v.Title, v.Link, html.UnescapeString(attackerNickname)})
		}
	} else {
		for _, url := range urlList {
			for _, v := range websitesList.Channel.Item {

				if strings.Contains(v.Title, getHostName(url)) {

					var attackerNickname string = html.UnescapeString(html.UnescapeString(r.ReplaceAllString(v.Description, "")))
					matchedDefacements = append(matchedDefacements, []string{v.Title, v.Link, attackerNickname})

					break
				}
			}
		}
	}

	return matchedDefacements
}

func getHostName(url string) string {
	r, _ := regexp.Compile(`:\/\/(.[^/]+)`)
	var hostName string = r.FindString(url)
	if hostName == "" {
		return "Host name could not detect."
	}
	return r.FindString(url)
}
