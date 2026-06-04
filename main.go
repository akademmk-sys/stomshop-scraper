package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Category struct {
	Name string
	Link string
	sub  []SubCategory
}

func (p *Category) HeaderWriter(c, l string, s []SubCategory, slice []Category) []Category {

	p.Name = c
	p.Link = l
	slice = append(slice, Category{Name: c, Link: l, sub: s})
	return slice
}

type SubCategory struct {
	SubName string
	SubLink string
}

func (s *SubCategory) subHeaderWriter(c, l string, slice []SubCategory) []SubCategory {

	s.SubName = c
	s.SubLink = l
	slice = append(slice, SubCategory{SubName: c, SubLink: l})
	return slice
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("stomshop.pro"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	})

	slice := make([]Category, 0)

	c.OnHTML("div.mycategory.product-thumb", func(r *colly.HTMLElement) {
		subSlice := make([]SubCategory, 0)
		var h Category
		var d SubCategory
		r.ForEach("div.child ul.list-unstyled li a", func(_ int, a *colly.HTMLElement) {
			if a.Text == "показать все ...." {
				return
			}
			subSlice = d.subHeaderWriter(a.Text, a.Attr("href"), subSlice)

		})
		r.ForEach("div.child h4 a", func(_ int, a *colly.HTMLElement) {

			slice = h.HeaderWriter(a.Text, a.Attr("href"), subSlice, slice)
		})

	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Ошибка:", r.StatusCode, err)
	})

	c.Visit("https://stomshop.pro/")
	fmt.Println(slice[0])
	fmt.Println(slice[1])

}
