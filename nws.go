package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/mmcdole/gofeed"
)

type nwsItem struct {
	Title    string
	Href     string
	Content  string
	ImageURL string
}

func grabFeed() (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	return fp.ParseURL("https://www.vrt.be/vrtnws/nl.rss.articles.xml")
}

func serveNewsList(c echo.Context) error {
	feed, err := grabFeed()
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	tmpl := template.Must(template.ParseFiles("./static/nws/list.wml"))

	c.Response().Header().Set("Content-Type", "text/vnd.wap.wml")
	c.Response().Header().Set("Cache-Control", "no-cache, must-revalidate")

	nwsItems := []nwsItem{}
	for _, item := range feed.Items {
		nwsItems = append(nwsItems, nwsItem{
			Title: item.Title,
			Href:  fmt.Sprintf("/nws/item?id=%s", template.URLQueryEscaper(item.GUID)),
		})
	}
	return tmpl.Execute(c.Response().Writer, struct{ Items []nwsItem }{Items: nwsItems})
}

func serveNewsItem(c echo.Context) error {
	feed, err := grabFeed()
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	tmpl := template.Must(template.ParseFiles("./static/nws/item.wml"))

	id := c.QueryParam("id")
	var article *gofeed.Item

	for _, item := range feed.Items {
		if item.GUID == id {
			article = item
		}
	}

	if article == nil {
		log.Println(id)
		return c.String(http.StatusNotFound, "")
	}

	item := nwsItem{
		Title:   article.Title,
		Content: article.Description,
	}

	if article.Image != nil {
		item.ImageURL = fmt.Sprintf("http://photon.innovatete.ch/%s?w=100&f=jpg", strings.Replace(article.Image.URL, "https://", "", -1))
	} else if len(article.Enclosures) > 0 {
		for _, enclosure := range article.Enclosures {
			if strings.HasPrefix(enclosure.Type, "image") {
				item.ImageURL = fmt.Sprintf("http://photon.innovatete.ch/%s?w=100&f=jpg", strings.Replace(enclosure.URL, "https://", "", -1))
				break
			}
		}
	}

	c.Response().Header().Set("Content-Type", "text/vnd.wap.wml")
	err = tmpl.Execute(c.Response().Writer, struct{ Item nwsItem }{Item: item})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
