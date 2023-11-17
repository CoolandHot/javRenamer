// to initialize the go module, run the following commands:
// ##########################
// go mod init CoolandHot/javRenamer
// go mod tidy
// ##########################

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"sync"

	"net/url"

	"strconv"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
)

// set proxy and choose javbus or avmoo
const (
	SOCKsProxy = "127.0.0.1:1080"
)

// set order
var reNameOrder = [...]string{"actress", "publishDate", "javID", "title"}

// it's not slice, therefore can't append()
var JavSites = [...]string{
	"https://www.javbus.com/search/",
	"https://avmoo.online/cn/search/",
	"http://www.javlibrary.com/cn/vl_searchbyid.php?keyword=",
}
var JavMatchRules = map[string][]string{
	"actress": {
		".star-name",
		"#avatar-waterfall .avatar-box span",
		"#video_cast span.star",
	},
	"title": {
		"div.container > h3",
		"div.container > h3",
		"#video_title > h3 > a",
	},
	"javID": {
		"div.row.movie > div.col-md-3.info > p:nth-child(1) > span:nth-child(2)",
		"div.row.movie > div.col-md-3.info > p:nth-child(1) > span:nth-child(2)",
		"#video_id > table > tbody > tr > td.text",
	},
	"publishDate": {
		"div.row.movie > div.col-md-3.info > p:nth-child(2)",
		"div.row.movie > div.col-md-3.info > p:nth-child(2)",
		"#video_date > table > tbody > tr > td.text",
	},
	"searchResults": {
		".movie-box",
		".movie-box",
		".video > a",
	},
}
var siteX int = 0
var with_proxy bool = true

func clientScrape(link string) *goquery.Document {
	// Create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", SOCKsProxy, nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}
	// Setup HTTP transport and a client
	tr := &http.Transport{Dial: dialer.Dial}
	client := &http.Client{}
	if with_proxy {
		client = &http.Client{Transport: tr}
	}
	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
	// request.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode == http.StatusFound { // 302
		newURL, _ := resp.Location()
		fmt.Println("Redirected to ", newURL.String())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK { // 200
		var doc *goquery.Document
		doc, _ = goquery.NewDocumentFromReader(resp.Body)
		return doc
	} else {
		fmt.Println(request.URL.Host, resp.Status)
		return nil
	}
}

func getDetail(doc *goquery.Document) (string, string, string, string) { // go into the detail pages
	var title, publishDate, heroine, javID string
	javID = doc.Find(JavMatchRules["javID"][siteX]).Text()
	// javID = strings.TrimSpace(javID)
	title = strings.ReplaceAll(doc.Find(JavMatchRules["title"][siteX]).Text(), javID+" ", "") // delete id inside
	title = strings.ReplaceAll(title, "/", " ")                                               // delete illegal strings
	title = strings.ReplaceAll(title, `\`, " ")                                               // delete illegal strings
	regexpMatch, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	publishDate = doc.Find(JavMatchRules["publishDate"][siteX]).Text()
	publishDate = regexpMatch.FindString(publishDate)
	var actresses []string
	doc.Find(JavMatchRules["actress"][siteX]).Each(func(i int, s *goquery.Selection) {
		actress := s.Text()
		actresses = append(actresses, actress)
	})
	if len(actresses) == 0 {
		heroine = "unknown"
	} else {
		heroine = strings.Join(actresses, " ")
	}
	return javID, title, publishDate, heroine
}

func getWebs(javID string) (string, string, string, string) { //get the search result
	var title, publishDate, heroine string
	doc := clientScrape(JavSites[siteX] + javID)
	if doc != nil {
		if doc.Find(JavMatchRules["searchResults"][siteX]).Length() != 0 {
			link, _ := doc.Find(JavMatchRules["searchResults"][siteX]).Eq(0).Attr("href")
			// in case the href is a relative link, convert it to absolute link
			baseUrl, err := url.Parse(JavSites[siteX]) // parse only base url
			if err != nil {
				log.Fatal(err)
			}
			fullUrl, err := baseUrl.Parse(link) // then use it to parse relative URLs
			if err != nil {
				log.Fatal(err)
			}
			doc = clientScrape(fullUrl.String())
		}
		javID, title, publishDate, heroine = getDetail(doc)
	}
	return javID, title, publishDate, heroine
}

func startRename(FullPath string, ch chan string, wg *sync.WaitGroup) {
	var title, publishDate, heroine, javID string
	// split filename and suffix
	fileFullnameMatch, _ := regexp.Compile(`[^\\]+\.[A-Za-z0-9]{3,10}$`)
	fileFullname := fileFullnameMatch.FindString(FullPath)
	basePath := strings.Replace(FullPath, fileFullname, "", 1)
	suffixMatch, _ := regexp.Compile(`\.[A-Za-z0-9]{3,10}$`)
	suffix := suffixMatch.FindString(fileFullname)
	filename := strings.Replace(fileFullname, suffix, "", 1)

	// look for jav-id
	matchRules := []string{
		`(:[^A-Za-z])?[A-Za-z]{1,5}-\d{3,5}(:\D)?`, // MIDE-939, C-2743
		`(:[^A-Za-z])?[A-Za-z]{1,5}\d{3,5}(:\D)?`,  // MIDE939, C2743
		`(:[^A-Za-z])?[A-Za-z]\d{2}-\d{3,5}(:\D)?`, // T28-556
	}
	for _, matchRule := range matchRules {
		javMatch, _ := regexp.Compile(matchRule)
		javIDs := javMatch.FindAllString(filename, -1)
		if len(javIDs) == 1 {
			javID, title, publishDate, heroine = getWebs(javIDs[0])
		}
		if title != "" {
			break
		}
	}
	if title == "" {
		fmt.Println(`fail on searching the javID`, fileFullname)
		return
	}

	// remove illegal characters in title
	illegal_chr, err2 := regexp.Compile(`[\\/:*?<>|]`)
	if err2 != nil {
		log.Fatal(err2)
	}
	title = illegal_chr.ReplaceAllString(title, "")

	// new name according to the order
	nameOrders := map[string]string{
		"actress":     heroine,
		"javID":       javID,
		"title":       title,
		"publishDate": publishDate}
	newFilename := nameOrders[reNameOrder[0]] + "-["
	for i := 1; i < 3; i++ {
		newFilename += nameOrders[reNameOrder[i]] + "]-["
	}
	newFilename += nameOrders[reNameOrder[3]] + "]" + suffix
	// rename file
	err := os.Rename(FullPath, basePath+newFilename)
	if err != nil {
		ch <- `fail on renaming file ` + fileFullname + ` to ` + newFilename
	} else {
		ch <- `successfully rename ` + fileFullname
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan string, 100)
	// cmd line arguments are strings, can't convert to int with int(), use strconv.ParseInt()
	site_num, err := strconv.ParseInt(os.Args[2], 10, 0)
	if err != nil {
		fmt.Println("site_num parameter must be integer")
		os.Exit(1)
	}
	siteX = int(site_num) //site_num is int64
	withProxy, err := strconv.ParseBool(os.Args[4])
	with_proxy = withProxy
	if err != nil {
		fmt.Println("with_proxy parameter error")
		os.Exit(1)
	}
	for _, arg := range strings.Split(os.Args[5], "***") {
		if arg != "" {
			// check if file exists
			if _, err := os.Stat(arg); err == nil {
				wg.Add(1) //increment the WaitGroup counter for each
				go startRename(arg, ch, &wg)
			}
		}
	}
	wg.Wait() //Block until the WaitGroup counter goes back to 0
	close(ch) //close the channel so that range iterator can operate on it
	for i := range ch {
		fmt.Println(i)
	}
}
