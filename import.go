package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	PublishedAt string   `xml:"pubDate"`
	Content     string   `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	Status      string   `xml:"status"`
	ItemType    string   `xml:"post_type"`
}

func runImport() {
	const layout = "Mon, 2 Jan 2006 15:04:05 -0700"

	db := InitDb()

	file, err := os.Open("tmp/wordpress.xml")
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}
	defer file.Close()

	var r Rss
	err = xml.NewDecoder(file).Decode(&r)
	if err != nil {
		fmt.Println("Error parsing XML: ", err)
	}

	for _, item := range r.Channel.Items {
		if item.Status == "publish" && item.ItemType == "post" {
			t, _ := time.Parse(layout, item.PublishedAt)
			post := Post{
				Title:     item.Title,
				Body:      item.Content,
				CreatedAt: t,
				UpdatedAt: t,
			}
			db.Create(&post)
		}
	}
}
