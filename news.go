package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type News struct {
	Date     time.Time `json:"date"`
	Category string    `json:"category"`
	Text     string    `json:"text"`
	Link     string    `json:"link"`
}

func getNewsURL(year, month int) (*url.URL, error) {
	yearStr := strconv.Itoa(year)
	monthStr := fmt.Sprintf("%02d", month)
	urlStr := endpointHost + endpointPathNews + "?" + endpointQueryNews +
		"&dy=" + yearStr + monthStr
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func parseNewsResponse(resp *http.Response) ([]News, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	newsEntries := []News{}
	doc.Find(".box-news ul li").Each(func(_ int, s *goquery.Selection) {
		news := News{}
		news.Date = getNewsDate(s)
		news.Category = getNewsCategory(s)
		news.Text = getNewsText(s)
		news.Link = getNewsLink(s)
		newsEntries = append(newsEntries, news)
	})
	return newsEntries, nil
}

func getNewsDate(s *goquery.Selection) time.Time {
	dateElm := s.Find(".date")
	dateText := strings.TrimSpace(dateElm.Text())
	dateTextSplit := strings.Split(dateText, ".")
	year, _ := strconv.Atoi(dateTextSplit[0])
	month, _ := strconv.Atoi(dateTextSplit[1])
	day, _ := strconv.Atoi(dateTextSplit[2])
	loc, _ := time.LoadLocation("Asia/Tokyo")
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
}

func getNewsCategory(s *goquery.Selection) string {
	categoryElm := s.Find(".category")
	return strings.TrimSpace(categoryElm.Text())
}

func getNewsText(s *goquery.Selection) string {
	textElm := s.Find(".text")
	anchor := textElm.Children().First()
	return anchor.Text()
}

func getNewsLink(s *goquery.Selection) string {
	textElm := s.Find(".text")
	anchor := textElm.Children().First()
	href, exists := anchor.Attr("href")
	if exists {
		return endpointHost + href
	}
	return ""
}
