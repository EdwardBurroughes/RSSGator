package feeds

import (
	"encoding/xml"
	"html"
)

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}

func (r *RSSFeed) unescapeFields() {
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Link = html.UnescapeString(r.Channel.Link)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)

	cleanItem := make([]RSSItem, 0)
	for _, item := range r.Channel.Item {
		rssItem := RSSItem{
			Title:       html.UnescapeString(item.Title),
			Link:        html.UnescapeString(item.Link),
			Description: html.UnescapeString(item.Description),
			PubDate:     html.UnescapeString(item.PubDate),
		}
		cleanItem = append(cleanItem, rssItem)
	}
	r.Channel.Item = cleanItem
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Item        []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}
