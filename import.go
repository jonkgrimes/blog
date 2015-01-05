package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/microcosm-cc/bluemonday"
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

func importBlog() {
	const layout = "Mon, 2 Jan 2006 15:04:05 -0700"
	p := bluemonday.StrictPolicy()

	db := InitDb()

	file, err := os.Open("data/wordpress.xml")
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
				Body:      p.Sanitize(item.Content),
				CreatedAt: t,
				UpdatedAt: t,
			}
			db.Create(&post)
		}
	}
}

/*
type Post struct {
	Id        int64
	Title     string `sql:"size:255"`
	Body      string `sql:"text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func InitDb() gorm.DB {
	db, err := gorm.Open("postgres", "user=jonkgrimes dbname=blog_development sslmode=disable")

	checkErr(err, "gorm.Open failed")

	db.AutoMigrate(&Post{})
	db.LogMode(true)

	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
*/
