package rss

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseRSS1(data []byte) (*Feed, error) {
	feed := rss1_0Feed{}
	p := xml.NewDecoder(bytes.NewReader(data))
	p.CharsetReader = charsetReader
	err := p.Decode(&feed)
	if err != nil {
		return nil, err
	}
	if feed.Channel == nil {
		return nil, fmt.Errorf("Error: no channel found in %q.", string(data))
	}

	channel := feed.Channel

	out := new(Feed)
	out.Title = channel.Title
	out.Description = channel.Description
	out.Link = channel.Link
	out.Image = channel.Image.Image()

	if feed.Items == nil {
		return nil, fmt.Errorf("Error: no feeds found in %q.", string(data))
	}

	out.Items = make([]*Item, 0, len(feed.Items))
	out.ItemMap = make(map[string]struct{})

	// Process items.
	for _, item := range feed.Items {

		if item.GUID == "" {
			if item.Link == "" {
				fmt.Printf("Warning: Item %q has no ID or link and will be ignored.\n", item.Title)
				continue
			}
			item.GUID = item.Link
		}

		next := new(Item)
		next.Title = item.Title
		next.Content = item.Content
		next.Link = item.Link
		if item.Date != "" {
			next.Date, err = parseTime(item.Date)
			if err != nil {
				return nil, err
			}
		} else if item.PubDate != "" {
			next.Date, err = parseTime(item.PubDate)
			if err != nil {
				return nil, err
			}
		}
		next.GUID = item.GUID
		next.Read = false
		if item.Media != nil && item.Media[len(item.Media)-1].Url != "" {
			enclosure := Enclosure{}
			enclosure.Url = item.Media[len(item.Media)-1].Url
			next.Enclosure = enclosure
		} else if strings.Contains(item.Content, "<img") {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(item.Content))
			if err != nil {
				fmt.Println(err)
			}
			imgSrc, _ := doc.Find("img").First().Attr("src")
			enclosure := Enclosure{}
			enclosure.Url = imgSrc
			next.Enclosure = enclosure

		}
		if _, ok := out.ItemMap[next.GUID]; ok {
			fmt.Printf("Warning: Item %q has duplicate ID.\n", next.Title)
			continue
		}

		out.Items = append(out.Items, next)
		out.ItemMap[next.GUID] = struct{}{}
		out.Unread++
	}

	return out, nil
}

type rss1_0Feed struct {
	XMLName xml.Name       `xml:"RDF"`
	Channel *rss1_0Channel `xml:"channel"`
	Items   []rss1_0Item   `xml:"item"`
}

type rss1_0Channel struct {
	XMLName     xml.Name    `xml:"channel"`
	Title       string      `xml:"title"`
	Description string      `xml:"description"`
	Link        string      `xml:"link"`
	Image       rss1_0Image `xml:"image"`
	MinsToLive  int         `xml:"ttl"`
	SkipHours   []int       `xml:"skipHours>hour"`
	SkipDays    []string    `xml:"skipDays>day"`
}

type rss1_0Item struct {
	XMLName   xml.Name  `xml:"item"`
	Title     string    `xml:"title"`
	Content   string    `xml:"description"`
	Link      string    `xml:"link"`
	PubDate   string    `xml:"pubDate"`
	Date      string    `xml:"date"`
	GUID      string    `xml:"guid"`
	Enclosure Enclosure `xml:"enclosure"`
	Media     []Media   `xml:"http://search.yahoo.com/mrss/ content"`
}

type rss1_0Image struct {
	XMLName xml.Name `xml:"image"`
	Title   string   `xml:"title"`
	Url     string   `xml:"url"`
	Height  int      `xml:"height"`
	Width   int      `xml:"width"`
}

func (i *rss1_0Image) Image() *Image {
	out := new(Image)
	out.Title = i.Title
	out.Url = i.Url
	out.Height = uint32(i.Height)
	out.Width = uint32(i.Width)
	return out
}
