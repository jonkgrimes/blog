package main

import (
	"encoding/xml"
	"fmt"
	"os"
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

func main() {
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
			fmt.Println(item.Title)
			fmt.Println(item.PublishedAt)
			fmt.Println(item.Content, "\n")
		}
	}
}
