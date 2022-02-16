# Project

Defalyzer is a cross-platform software that focuses on defacement analyzing and mirror tracking. You can track your (or popular) websites against defacements. We also provide defacement announcements in our [Twitter](https://twitter.com/defalyzer) account and [Telegram](https://t.me/vullnerability) channel.

## How it works?

It simply:
1. collects websites you would like to scan, from the websites.txt file
1. scans the website and gets source code, visible text and full-size screenshot,
2. sends screenshot to Google OCR Servers to get an accurate plaintext output,
2. checks all the collected data if those include some "hacked" keywords or popular defacer nicknames,
1. checks Zone-H mirror database if the website is recently noticed, 
4. prints scan results and logs details into the file (defaced_logs.txt)

> Defalyzer has two features: it scans websites directly, and it checks them on the Zone-H as default. You can also use only Zone-H tracking option instead of scanning all URLs and you can simply filter domains by their names/extensions.

## Installation

**NOTE! Defalyzer requires a Chrome browser to run. Please be sure that you have installed the browser first.**

After downloading the content, you can simply compile the script by this command:

`go build main.go`

Now you can see the commands you will use :)

`./main --help`

> Tested on MacOS Monterey / Windows 10 / Ubuntu 20.

**NOTE! If would you like to scan your custom websites by analyzing screenshots, you need to provide google.json file in the working directory.**

> You can find your JSON service account file by following: "IAM & Admin > Service Accounts" tabs in Google Cloud Console. Download it to your work directory and rename it as google.json for enabling OCR scanning. You will find an example of google.json file in the repo. [Tutorial video to download JSON service file.](https://www.youtube.com/watch?v=rWcLDax-VmM)

> If you still have an "OCR Error" warning, try to set an environment variable (GOOGLE_APPLICATION_CREDENTIALS) manually by reviewing [this document.](https://cloud.google.com/vision/docs/ocr) Look at the "Set up your GCP project and authentication" title.

## Usage and screenshots

Defalyzer has some flags (parameters) to customize or make your queries faster.

| Parameter     |  Description |
| ----------- | ----------- |
| --ext      |  Filter output by domain name/extension. For example, enter ".gov,.gov.br" as value to eliminate other extensions. |
| --zoneh   | Enable Zone-H tracking. Enter 'all' as value to track all URLs on Zone-H. Enter 'file' as value to track custom URLs only (in websites.txt) on Zone-H. |
| --intv | Re-scan timing as second for hacked websites only. The default value is 600, it means do not scan the previously hacked website earlier than 600 seconds. |
| --allintv | Enable re-scan timing for all websites, not hacked ones only. |
| --loop | Enable infinite loop. When this parameter was not used, the scanning will happen once. |
| --file | Enable website tracking from the websites.txt file. |
| --defonly | Print defaced websites only. |
| --no-color | Disable colorized output. |

- Let us use --zoneh and --file parameters together. So, it will scan URLs from website.txt and check those URLs on Zone-H:
![](https://i.hizliresim.com/sg6dq44.png)

- Now, it will only enable Zone-H tracking and will not scan websites in the websites.txt file:
![](https://i.hizliresim.com/glco0fx.png)

- We can also filter output by using --ext and --defonly parameters together. This command will show "defaced .go.id and .gov.br domains" only:
![](https://i.hizliresim.com/7o72uw2.png)


## TO-DO List
- Improve defacement analyzing conditions by creating a defacement dataset and return an hacked-score between 0.0-1.0.
- Integrate other mirror databases.
- Enlarge popular defacers wordlist. (nickname_wordlist.txt)

## Thanks ❤
- Berat SULAR
- Rıza SABUNCU
- Elif ÖNEY
- IKU1337 Cybersecurity Society

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)