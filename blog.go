package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Blog struct of blog
type Blog struct {
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Author       string    `json:"author"`
	PostedDate   time.Time `json:"postedDate"`
	ThumbnailURL string    `json:"thumbnailURL"`
}

func getBlogURL(memberID, page int) (*url.URL, error) {
	rawurl := endpointHost + endpointPathBlog + "?" + endpointQueryBlog

	if _, ok := members[memberID]; ok {
		rawurl += "&ct=" + getMemberID(memberID)
	}

	if page >= 0 {
		rawurl += "&page=" + strconv.Itoa(page)
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func parseBlogResponse(resp *http.Response) ([]Blog, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	blogs := []Blog{}
	doc.Find("article").Each(func(_ int, s *goquery.Selection) {
		blog := Blog{}
		articleChildren := s.Children()

		blog.Title = getTitle(articleChildren)
		blog.URL = getLink(articleChildren)
		blog.Author = getAuthor(articleChildren)
		blog.PostedDate = getPostedDate(articleChildren)
		blog.ThumbnailURL = getThumbnailURL(articleChildren)

		blogs = append(blogs, blog)
	})
	return blogs, nil
}

func getTitle(s *goquery.Selection) string {
	title := s.Find("h3").Text()
	return strings.TrimSpace(title)
}

func getLink(s *goquery.Selection) string {
	href, exists := s.Find("h3 a").Attr("href")
	if exists {
		return endpointHost + href
	}
	return ""
}

func getAuthor(s *goquery.Selection) string {
	author := s.Find("p.name").Text()
	return strings.TrimSpace(author)
}

func getPostedDate(s *goquery.Selection) time.Time {
	dateElm := s.Find("div.box-date").Children()
	time1 := dateElm.First().Text()
	time2 := dateElm.Last().Text()
	year, _ := strconv.Atoi(strings.Split(time1, ".")[0])
	month, _ := strconv.Atoi(strings.Split(time1, ".")[1])
	day, _ := strconv.Atoi(time2)
	loc, _ := time.LoadLocation("Asia/Tokyo")
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
}

func getThumbnailURL(s *goquery.Selection) string {
	src, exists := s.Find("img").Attr("src")
	if exists {
		return src
	}
	return ""
}
