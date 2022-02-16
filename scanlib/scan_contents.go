package scanUri

import (
	"bytes"
	"context"
	"log"
	urlp "net/url"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
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
	u, _ := urlp.Parse(url)

	var extensions []string = strings.Split(u.Host, ".")
	var ext string = extensions[len(extensions)-1]
	// Alternative solution by regex: https://play.golang.org/p/BibG75p2Jz4

	countries := map[string]string{
		"tr": "Turkey", "gov": "USA", "edu": "USA", "mil": "USA", "ac": "Ascention Island", "ad": "Andorra", "ae": "United Arab Emirates", "af": "Afghanistan", "ag": "Antigua and Barbuda", "ai": "Anguillia", "al": "Albania", "am": "Armenia", "ao": "Angola", "aq": "Antarctique", "ar": "Argentina", "as": "American Samoa", "at": "Austria", "au": "Australia", "aw": "Aruba", "ax": "Åland", "az": "Azerbaijan", "ba": "Bosnia and Herzegovina", "bb": "Barbados", "bd": "Bangladesh", "be": "Belgium", "bf": "Burkina Faso", "bg": "Bulgaria", "bh": "Bahrain", "bi": "Burundi", "bj": "Benin", "bm": "Bermuda", "bn": "Brunei", "bo": "Bolivia", "bq": "Caribbean Netherlands", "br": "Brazil", "bs": "Bahamas", "bt": "Bhutan", "bw": "Botswana", "by": "Belarus", "bz": "Belize", "ca": "Canada", "cc": "Cocos (Keeling) Islands", "cd": "Democratic Republic of the Congo", "cf": "Central African Republic", "cg": "Republic of the Congo", "ch": "Switzerland", "ci": "Ivory Coast", "ck": "Cook Islands", "cl": "Chile", "cm": "Cameroon", "cn": "People's Republic of China", "co": "Colombia", "cr": "Costa Rica", "cu": "Cuba", "cv": "Cape Verde", "cw": "Curaçao", "cx": "Christmas Island", "cy": "Cyprus", "cz": "Czech Republic", "de": "Germany", "dj": "Djibouti", "dk": "Denmark", "dm": "Dominica", "do": "Dominican Republic", "dz": "Algeria", "ec": "Ecuador", "ee": "Estonia", "eg": "Egypt", "eh": "Western Sahara", "er": "Eritrea", "es": "Spain", "et": "Ethiopia", "eu": "European Union", "fi": "Finland", "fj": "Fiji", "fk": "Falkland Islands", "fm": "Federated States of Micronesia", "fo": "Faroe Islands", "fr": "France", "ga": "Gabon", "gd": "Grenada", "ge": "Georgia", "gf": "French Guiana", "gg": "Guernsey", "gh": "Ghana", "gi": "Gibraltar", "gl": "Greenland", "gm": "The Gambia", "gn": "Guinea", "gp": "Guadeloupe", "gq": "Equatorial Guinea", "gr": "Greece", "gs": "South Georgia and the South Sandwich Islands", "gt": "Guatemala", "gu": "Guam", "gw": "Guinea-Bissau", "gy": "Guyana", "hk": "Hong Kong", "hm": "Heard Island and McDonald Islands", "hn": "Honduras", "hr": "Croatia", "ht": "Haiti", "hu": "Hungary", "id": "Indonesia", "ie": "Ireland", "il": "Israel", "im": "Isle of Man", "in": "India", "io": "British Indian Ocean Territory", "iq": "Iraq", "ir": "Iran", "is": "Iceland", "it": "Italy", "je": "Jersey", "jm": "Jamaica", "jo": "Jordan", "jp": "Japan", "ke": "Kenya", "kg": "Kyrgyzstan", "kh": "Cambodia", "ki": "Kiribati", "km": "Comoros", "kn": "Saint Kitts and Nevis", "kp": "North Korea", "kr": "South Korea", "kw": "Kuwait", "ky": "Cayman Islands", "kz": "Kazakhstan", "la": "Laos", "lb": "Lebanon", "lc": "Saint Lucia", "li": "Liechtenstein", "lk": "Sri Lanka", "lr": "Liberia", "ls": "Lesotho", "lt": "Lithuania", "lu": "Luxembourg", "lv": "Latvia", "ly": "Libya", "ma": "Morocco", "mc": "Monaco", "md": "Moldova", "me": "Montenegro", "mg": "Madagascar", "mh": "Marshall Islands", "mk": "North Macedonia", "ml": "Mali", "mm": "Myanmar", "mn": "Mongolia", "mo": "Macau", "mp": "Northern Mariana Islands", "mq": "Martinique", "mr": "Mauritania", "ms": "Montserrat", "mt": "Malta", "mu": "Mauritius", "mv": "Maldives", "mw": "Malawi", "mx": "Mexico", "my": "Malaysia", "mz": "Mozambique", "na": "Namibia", "nc": "New Caledonia", "ne": "Niger", "nf": "Norfolk Island", "ng": "Nigeria", "ni": "Nicaragua", "nl": "Netherlands", "no": "Norway", "np": "Nepal", "nr": "Nauru", "nu": "Niue", "nz": "New Zealand", "om": "Oman", "pa": "Panama", "pe": "Peru", "pf": "French Polynesia", "pg": "Papua New Guinea", "ph": "Philippines", "pk": "Pakistan", "pl": "Poland", "pm": "Saint-Pierre and Miquelon", "pn": "Pitcairn Islands", "pr": "Puerto Rico", "ps": "Palestine", "pt": "Portugal", "pw": "Palau", "py": "Paraguay", "qa": "Qatar", "re": "Réunion", "ro": "Romania", "rs": "Serbia", "ru": "Russia", "rw": "Rwanda", "sa": "Saudi Arabia", "sb": "Solomon Islands", "sc": "Seychelles", "sd": "Sudan", "se": "Sweden", "sg": "Singapore", "sh": "Saint Helena", "si": "Slovenia", "sk": "Slovakia", "sl": "Sierra Leone", "sm": "San Marino", "sn": "Senegal", "so": "Somalia", "sr": "Suriname", "ss": "South Sudan", "st": "São Tomé and Príncipe", "su": "Soviet Union", "sv": "El Salvador", "sx": "Sint Maarten", "sy": "Syria", "sz": "Eswatini", "tc": "Turks and Caicos Islands", "td": "Chad", "tf": "French Southern and Antarctic Lands", "tg": "Togo", "th": "Thailand", "tj": "Tajikistan", "tk": "Tokelau", "tl": "East Timor", "tm": "Turkmenistan", "tn": "Tunisia", "to": "Tonga", "tt": "Trinidad and Tobago", "tv": "Tuvalu", "tw": "Taiwan", "tz": "Tanzania", "ua": "Ukraine", "ug": "Uganda", "uk": "United Kingdom", "us": "USA", "uy": "Uruguay", "uz": "Uzbekistan", "va": "Vatican City", "vc": "Saint Vincent and the Grenadines", "ve": "Venezuela", "vg": "UK Virgin Islands", "vi": "US Virgin Islands", "vn": "Vietnam", "vu": "Vanuatu", "wf": "Wallis and Futuna", "ws": "Samoa", "ye": "Yemen", "yt": "Mayotte", "za": "South Africa", "zm": "Zambia", "zw": "Zimbabwe",
	}

	value, ok := countries[ext]
	if ok {
		return value
	}
	return "International"
}
