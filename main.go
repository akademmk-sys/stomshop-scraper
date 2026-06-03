package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

type ParsedHeaders struct {
	name string
	link string
}

func (p *ParsedHeaders) headerWriter(c, l string) []ParsedHeaders {
	slice := make([]ParsedHeaders, 0)
	p.name = c
	p.link = l
	slice = append(slice, ParsedHeaders{name: c, link: l})
	return slice
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("stomshop.pro"),
	)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	})
	c.OnHTML("#content > div.row.product-layout", func(r *colly.HTMLElement) {

		r.ForEach("div.child h4 a", func(_ int, a *colly.HTMLElement) {
			var h ParsedHeaders
			str := h.headerWriter(a.Text, a.Attr("href"))
			fmt.Println(str)
		})
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Ошибка:", r.StatusCode, err)
	})
	c.Visit("https://stomshop.pro/")
}
