package search

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type CrawlData struct {
	Url          string
	Success      bool
	ResponseCode int
	CrawlData    ParseBody
}

type ParseBody struct {
	CrawlTime       time.Duration
	PageTitle       string
	PageDescription string
	Headings        string
	Links           Links
}

type Links struct {
	Internal []string
	External []string
}

func runCrawler(inputUrl string) CrawlData {
	resp, err := http.Get(inputUrl)
	baseUrl, _ := url.Parse(inputUrl)

	if err != nil || resp == nil {
		fmt.Println("Something WAENT Wrong while fetching body")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: 0, CrawlData: ParseBody{}}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Not Status Code of 200 Found")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParseBody{}}
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {

	} else {
		fmt.Println("Non HTML Response")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParseBody{}}
	}
}

func parseBody(body io.Reader, baseUrl *url.URL) (ParseBody, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return ParseBody{}, err
	}
	start := time.Now()

	// Get Links
	links := getLinks()

	// Get Page Title and description
	title, desc := getPageData(doc)

	// Get H1 Tags
	headings := getPageHeadings(doc)

	// Return The Time and data
	end := time.Now()

	return ParseBody{
		CrawlTime:       end.Sub(start),
		PageTitle:       title,
		PageDescription: desc,
		Headings:        headings,
		Links:           links,
	}, nil
}

func getPageData(node *html.Node) (string, string) {
	if node == nil {
		return "", ""
	}
	title, description := "", ""
	var findMetaAndTitle func(*html.Node)
	findMetaAndTitle = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "title" {
			if node.FirstChild == nil {
				title = " "
			} else {
				title = node.FirstChild.Data
			}
		}else if node.Type==html.ElementNode && node.Data="meta"{
			var name,content string
			for _,attr:=range  node.Attr{
				if attr.Key=="name" {
					name=attr.Val	
				} else if attr.Key=="content"{
					content=attr.Val
				}
			}
			if name=="description" {
				description=content
			}
		}

	}
	for child :=node.FirstChild;child != nil; child=child.NextSibling{
		findMetaAndTitle(child)
	}
	findMetaAndTitle(node)
	return title,description

}





func getPageHeadings(n *html.Node) string {
	if n == nil {
		return ""
	}
	var headings strings.Builder
	var findH1 func(*html.Node)
	findH1 = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			if n.FirstChild != nil {
				headings.WriteString(n.FirstChild.Data)
				headings.WriteString(", ")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findH1(c)
		}
	}

	return strings.TrimSuffix(headings.String(), ",")

}
