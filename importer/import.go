package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/danryan/env"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/microcosm-cc/bluemonday"
)

type Config struct {
	DbName      string `env:"key=BLOG_DATABASE_NAME default=blog_development"`
	DbUser      string `env:"key=BLOG_DATABASE_USER default=jonkgrimes"`
	DbPassword  string `env:"key=BLOG_DATABASE_PASSWORD"`
	Port        string `env:"key=BLOG_PORT default=:8080"`
	Environment string `env:"key=ENVIRONMENT default=development"`
	NewRelicKey string `env:"key=NEW_RELIC_KEY"`
}

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
	config := &Config{}
	if err := env.Process(config); err != nil {
		fmt.Println(err)
	}

	const layout = "Mon, 2 Jan 2006 15:04:05 -0700"
	p := bluemonday.StrictPolicy()

	db := InitDb(config)

	file, err := os.Open("wordpress.xml")
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

type Post struct {
	Id        int64
	Title     string `sql:"size:255"`
	Body      string `sql:"text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func InitDb(c *Config) gorm.DB {
	tmpl, err := template.New("connection").Parse("user={{.DbUser}}{{if .DbPassword}} password={{.DbPassword}}{{end}} dbname={{.DbName}} sslmode=disable")
	var b bytes.Buffer
	err = tmpl.Execute(&b, c)
	connString := b.String()

	db, err := gorm.Open("postgres", connString)

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
