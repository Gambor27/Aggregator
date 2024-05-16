package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSResponse struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Pubdate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSResponse, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSResponse{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSResponse{}, err
	}

	rssFeed := RSSResponse{}

	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return RSSResponse{}, err
	}

	return rssFeed, nil
}
