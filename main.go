package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	urlp "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	scanUri "github.com/NumanABi/Defalyzer/scanlib"
	"github.com/chromedp/chromedp"
	"github.com/fatih/color"
)

var (
	extensions        []string                              // Array for storing extensions to filter
	urlList           []string                              // Array for storing URLs to scan
	hacked_keywords   []string                              // Array for storing defacement keywords
	nickname_keywords []string                              // Array for storing defacer nicknames
	scanFromFile      bool             = false              // Enable scanning custom websites
	extLookup         bool             = false              // Enable filtering by domain extensions
	infiniteLoop      bool             = false              // Enable infinite loop
	zoneHLookup       bool             = false              // Enable ZoneH Tracking at default
	zoneHAllExtract   bool             = false              // Track all websites in ZoneH Special Archive; not only for websites in websites.txt
	allInterval       bool             = false              // Apply interval time to all websites, not only for defaced ones
	printDefacedOnly  bool             = false              // Print only defaced websites
	countLogs         int              = 0                  // If is there anything to print, increase the count
	setInterval       int64            = 600                // If website is defaced, do not scan it again earlier than X seconds
	intervalMap       map[string]int64 = map[string]int64{} // Note: It is not a queue or sleeping job; so timing may change according to lenght of websites list
)

func main() {

	// Set environment variable for Google OCR
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "google.json")

	// Import websites list to scan
	websites_file, err := os.ReadFile("websites.txt")
	if err != nil {
		log.Fatal("Could not read websites.txt file.")
	}

	// Import defacement keywords to scan
	hacked_k_file, err := os.ReadFile("hacked_wordlist.txt")
	if err != nil {
		log.Fatal("Could not read hacked_wordlist.txt file.")
	}

	// Import defacer nicknames to scan
	nick_k_file, err := os.ReadFile("nickname_wordlist.txt")
	if err != nil {
		log.Fatal("Could not read nickname_wordlist.txt file.")
	}

	// Append websites to urlList array
	scanner := bufio.NewScanner(bytes.NewReader(websites_file))
	for scanner.Scan() {
		urlList = append(urlList, scanner.Text())
	}

	// Append defacement keywords to to hacked_keywords array
	scanner = bufio.NewScanner(bytes.NewReader(hacked_k_file))
	for scanner.Scan() {
		hacked_keywords = append(hacked_keywords, scanner.Text())
	}

	// Append nicknames to nickname_keywords array
	scanner = bufio.NewScanner(bytes.NewReader(nick_k_file))
	for scanner.Scan() {
		nickname_keywords = append(nickname_keywords, scanner.Text())
	}

	// Check for CLI arguments
	if len(os.Args) > 1 {
		extArgv := flag.String("ext", "*", "Filter output by domain name/extension. For example, enter \".gov,.gov.br\" as value to eliminate other extensions.")
		zoneHArgv := flag.String("zoneh", "file", "Enable Zone-H tracking. Enter 'all' as value to track all URLs on Zone-H. Enter 'file' as value to track custom URLs (in websites.txt) on Zone-H.")
		internalArgv := flag.Int("intv", 600, "Re-scan timing as second for hacked websites only. The default value is 600, it means do not scan the previously hacked website earlier than 600 seconds.")
		loopArgv := flag.Bool("loop", false, "Enable infinite loop. When this parameter was not used, the scanning will happen for once.")
		allIntvArg := flag.Bool("allintv", false, "Enable re-scan timing for all websites, not hacked ones only.")
		fileArgv := flag.Bool("file", false, "Enable website tracking from the websites.txt file.")
		defArgv := flag.Bool("defonly", false, "Print defaced websites only.")
		disableColorizedOutput := flag.Bool("no-color", false, "Disable colorized output.")
		flag.Parse()

		// If extensions defined by --ext argument, parse extensions by comma and append to extensions array
		if *extArgv != "*" {
			for _, v := range strings.Split(*extArgv, ",") {
				extensions = append(extensions, v)
				extLookup = true
			}
		}

		// Set interval time by --intv argument
		setInterval = int64(*internalArgv)

		if *fileArgv {
			scanFromFile = true
		}

		if *allIntvArg {
			allInterval = true
		}

		if *loopArgv {
			infiniteLoop = true
		}

		if *defArgv {
			printDefacedOnly = true
		}

		if *disableColorizedOutput {
			color.NoColor = true
		}

		if *zoneHArgv == "all" {
			zoneHLookup = true
			zoneHAllExtract = true
		} else if *zoneHArgv == "file" {
			zoneHLookup = true
		}

	}

	ascii := color.CyanString(fmt.Sprintf(` ____        __       _                    
|  _ \  ___ / _| __ _| |_   _ _______ _ __ 
| | | |/ _ \ |_ / _%c | | | | |_  / _ \ '__|
| |_| |  __/  _| (_| | | |_| |/ /  __/ |   
|____/ \___|_|  \__,_|_|\__, /___\___|_|   v1
			|___/             

Defacement Analyzer and Mirror Tracker
Follow Special Defacements on @defalyzer
Developed by @IKU1337 Cybersecurity Society
`, '`'))
	fmt.Println(ascii)

	if len(os.Args) == 1 {
		fmt.Println(color.WhiteString("Use --help parameter to see commands.\n"))
	}

	var wg sync.WaitGroup

	if zoneHLookup {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for {
				defer wg.Done()

				var zoneHOutput [][]string
				if zoneHAllExtract {
					// If "all" parameter passed to --zoneh argument, track all URLs on the Zone-H Special Archive
					zoneHOutput = scanUri.ExtractZoneH(urlList, "match_all_urls")
				} else {
					// If "file" parameter (default) passed to --zoneh argument, only track URLs in websites.txt on the Zone-H Special Archive
					zoneHOutput = scanUri.ExtractZoneH(urlList, "match_from_file")
				}
				for _, zoneHMirror := range zoneHOutput {

					var urlZoneH string = zoneHMirror[0]
					var zoneHMirrorLink string = zoneHMirror[1]
					var zoneHAttacker string = zoneHMirror[2]

					var continueScan bool = true

					if extLookup {
						continueScan = false
						parsedURL, _ := urlp.Parse(urlZoneH)
						for _, v := range extensions {
							// If extensions defined by user matches with domain's extension, continue to scan
							if parsedURL.Host[len(parsedURL.Host)-len(v):] == v {
								continueScan = true
								break
							}
						}
					}

					if continueScan {
						countLogs++
						var country string = scanUri.DetectCountry(urlZoneH) // Detect country of domain by extension
						var output string = fmt.Sprintf("| Address: %v | Zone-H: %v | Country: %v | Attacker: %v", urlZoneH, zoneHMirrorLink, country, zoneHAttacker)

						file, _ := os.OpenFile("defaced_logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
						defer file.Close()
						log.New(file, "", log.LstdFlags).Println(output) // Print output into file
						log.Println(color.MagentaString(output))         // Print output

					}
				}

				// If infinite loop is activated, sleep for 5 minutes and repeat
				if infiniteLoop {
					time.Sleep(time.Minute * 5)
				} else {
					break
				}
			}
		}(&wg)
	}

	if scanFromFile {

		// Run Chrome browser in Headless and No-Sandbox modes
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(findExecPath()),
			chromedp.Flag("headless", true),
			chromedp.Flag("no-sandbox", true),
		)

		// Create browser context
		parentCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		for {

			var numJobs int = len(urlList)
			jobs := make(chan int, numJobs)
			results := make(chan int, numJobs)

			for worker := 1; worker <= 20; worker++ {
				go func(id int, jobs <-chan int, results chan<- int) {
					for j := range jobs {
						loadURI(urlList[j], &parentCtx)
						results <- j
					}
				}(worker, jobs, results)
			}

			for j := 0; j < numJobs; j++ {
				jobs <- j
			}
			close(jobs)

			for a := 1; a <= numJobs; a++ {
				<-results
			}

			if infiniteLoop {
				// Wait for a while for next loop.
				if allInterval {
					time.Sleep(time.Second * time.Duration(setInterval))
				} else {
					//	time.Sleep(time.Minute * 1)
				}
			} else {
				break
			}

		}
	}

	wg.Wait()

	if countLogs == 0 && len(os.Args) > 1 {
		log.Println(color.YellowString("No match found."))
	}

}

func loadURI(url string, parentCtx *context.Context) {
	var continueScan bool = true

	if extLookup {
		continueScan = false
		parsedURL, _ := urlp.Parse(url)
		for _, v := range extensions {
			// If extensions defined by user matches with domain's extension, continue to scan
			if parsedURL.Host[len(parsedURL.Host)-len(v):] == v {
				continueScan = true
				break
			}
		}
	}

	// If interval time is up, continue to scan
	expTime, exist := intervalMap[url]
	if exist {
		if expTime > time.Now().Unix() {
			continueScan = false
		} else {
			delete(intervalMap, url)
		}
	}

	if continueScan {

		// Get screenshot of the page, source code and visible text
		imageData, htmlresponse, html_text := scanUri.GetSources(*parentCtx, url)

		if len(htmlresponse) > 0 {
			// Extract text from screenshot of the page, by using Google OCR
			var screenshotToText string = scanUri.GoogleOCR(url, imageData)

			var country string = scanUri.DetectCountry(url)
			// Compare screenshot text, source code and visible text with defacement keywords and nicknames
			isDefaced, hackedKeyword, attackerKeyword := scanUri.SearchText(screenshotToText, htmlresponse, html_text, hacked_keywords, nickname_keywords)
			var output string = fmt.Sprintf("| Address: %v | Country: %v | Hacked: %v | Attacker: %v | Keyword: %v", url, country, isDefaced, attackerKeyword, hackedKeyword)

			if isDefaced {
				file, _ := os.OpenFile("defaced_logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
				defer file.Close()
				log.New(file, "", log.LstdFlags).Println(output) // Print output into file
				// Apply interval time to domain and append with expiration time to intervalMap map
				if !allInterval {
					intervalMap[url] = time.Now().Unix() + setInterval
				}
				log.Println(color.RedString(output)) // Print output
				countLogs++
			} else if !printDefacedOnly {
				countLogs++
				log.Println(color.GreenString(output))
			}

		}

	}
}

func findExecPath() string {
	// This function will return path of the Chrome browser
	// Already defined in allocate.go, as a private function.
	// https://github.com/chromedp/chromedp/blob/master/allocate.go
	var locations []string
	switch runtime.GOOS {
	case "darwin":
		locations = []string{
			// Mac
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		}
	case "windows":
		locations = []string{
			// Windows
			"chrome",
			"chrome.exe", // in case PATHEXT is misconfigured
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\Google\Chrome\Application\chrome.exe`),
			filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\Chromium\Application\chrome.exe`),
		}
	default:
		locations = []string{
			// Unix-like
			"headless_shell",
			"headless-shell",
			"chromium",
			"chromium-browser",
			"google-chrome",
			"google-chrome-stable",
			"google-chrome-beta",
			"google-chrome-unstable",
			"/usr/bin/google-chrome",
			"/usr/local/bin/chrome",
			"/snap/bin/chromium",
			"chrome",
		}
	}

	for _, path := range locations {
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	log.Println(color.YellowString("Chrome path could not found."))
	return "google-chrome"
}
