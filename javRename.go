// go get github.com/PuerkitoBio/goquery
// go get -u golang.org/x/net

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
)

// set proxy and choose javbus or avmoo
const (
	SOCKsProxy = "127.0.0.1:1080"
	JavBus     = true
)

// set order
var reNameOrder = [...]string{"actress", "javID", "title", "publishDate"}

func clientScrape(link string) *goquery.Document {
	// Create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", SOCKsProxy, nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}
	// Setup HTTP transport and a client
	tr := &http.Transport{Dial: dialer.Dial}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
	request.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var doc *goquery.Document
	if resp.StatusCode == 200 {
		doc, _ = goquery.NewDocumentFromReader(resp.Body)
	}
	return doc
}

func getDetail(detailLink string) (string, string, string, string) { // go into the detail pages
	var title, publishDate, heroine, javID string
	doc := clientScrape(detailLink)
	javID = doc.Find("span.header").Eq(0).Next().Text()
	title = strings.ReplaceAll(doc.Find("h3").Eq(0).Text(), javID+" ", "") // delete id inside
	title = strings.ReplaceAll(title, "/", " ")                            // delete illegal strings
	title = strings.ReplaceAll(title, `\`, " ")                            // delete illegal strings
	regexpMatch, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	publishDate = doc.Find("span.header").Eq(1).Parent().Text()
	publishDate = regexpMatch.FindString(publishDate)
	var actresses []string
	if JavBus {
		doc.Find(".star-name").Each(func(i int, s *goquery.Selection) {
			actress := s.Text()
			actresses = append(actresses, actress)
		})
	} else {
		doc.Find(".avatar-box span").Each(func(i int, s *goquery.Selection) {
			actress := s.Text()
			actresses = append(actresses, actress)
		})
	}
	if len(actresses) == 0 {
		heroine = "unknown"
	} else {
		heroine = strings.Join(actresses, " ")
	}
	return javID, title, publishDate, heroine
}

func getWebs(javBus string, javID string) (string, string, string, string) { //get the search result
	var title, publishDate, heroine string
	doc := clientScrape(javBus + javID)
	doc.Find(".movie-box").Each(func(i int, content *goquery.Selection) {
		res1 := strings.ReplaceAll(content.Find("date").Eq(0).Text(), "-", "")
		res2 := strings.ReplaceAll(javID, "-", "") // in case any - remain
		res1, res2 = strings.ToUpper(res1), strings.ToUpper(res2)
		if res1 == res2 {
			link, _ := content.Attr("href")
			javID, title, publishDate, heroine = getDetail(link)
			return
		}
	})
	return javID, title, publishDate, heroine
}

func startRename(FullPath string, ch chan string) {
	var title, publishDate, heroine, avmoo, javID string
	if JavBus {
		avmoo = "https://www.javbus.com/search/"
	} else {
		avmoo = "https://avmoo.cyou/cn/search/"
	}
	// split filename and suffix
	fileFullnameMatch, _ := regexp.Compile(`[^\\]+\.[A-Za-z0-9]{3,10}$`)
	fileFullname := fileFullnameMatch.FindString(FullPath)
	basePath := strings.Replace(FullPath, fileFullname, "", 1)
	suffixMatch, _ := regexp.Compile(`\.[A-Za-z0-9]{3,10}$`)
	suffix := suffixMatch.FindString(fileFullname)
	filename := strings.Replace(fileFullname, suffix, "", 1)

	// look for jav-id
	matchRules := []string{
		`(:[^A-Za-z])?[A-Za-z]{2,5}-\d{3,5}(:\D)?`, // MIDE-939
		`(:[^A-Za-z])?[A-Za-z]{2,5}\d{3,5}(:\D)?`,  // MIDE939
		`(:[^A-Za-z])?[A-Za-z]\d{2}-\d{3,5}(:\D)?`, // T28-556
	}
	for _, matchRule := range matchRules {
		javMatch, _ := regexp.Compile(matchRule)
		javIDs := javMatch.FindAllString(filename, -1)
		if len(javIDs) == 1 {
			javID, title, publishDate, heroine = getWebs(avmoo, javIDs[0])
		}
		if title != "" {
			break
		}
	}
	if title == "" {
		fmt.Println(`fail on searching the javID`, fileFullname)
		return
	}

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
}

func main() {
	ch := make(chan string, 100)
	for _, arg := range strings.Split(os.Args[1], " ") {
		go startRename(arg, ch)
		fmt.Println(<-ch)
	}
}
