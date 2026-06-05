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
	SubName  string
	SubLink  string
	Products []Product
}

func (s *SubCategory) subHeaderWriter(c, l string, slice []SubCategory) []SubCategory {

	s.SubName = c
	s.SubLink = l
	slice = append(slice, SubCategory{SubName: c, SubLink: l})
	return slice
}

type Product struct {
	PName   string
	PLink   string
	Price   string
	Number  string
	Country string
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
			req, _ := r.Request.New("GET", a.Attr("href"), nil)
			req.Ctx.Put("sublink", a.Attr("href"))
			req.Do()
		})
		r.ForEach("div.child h4 a", func(_ int, a *colly.HTMLElement) {

			slice = h.HeaderWriter(a.Text, a.Attr("href"), subSlice, slice)
		})

	})

	c.OnHTML("div.product-layout div.product-thumb", func(r *colly.HTMLElement) {
		unit := Product{
			PName:   r.ChildText("div.caption a"),
			PLink:   r.ChildAttr("div.caption a", "href"),
			Price:   r.ChildText("span.price-new b"),
			Number:  r.ChildText("span.code span"),
			Country: r.ChildText("div.manufacturer a"),
		}
		subLink := r.Request.Ctx.Get("sublink")
		fmt.Println("subLink:", subLink)
		if subLink == "" {
			return
		}
		for i := range slice {
			for j := range slice[i].sub {
				if slice[i].sub[j].SubLink == subLink {
					slice[i].sub[j].Products = append(slice[i].sub[j].Products, unit)
				}
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Ошибка:", r.StatusCode, err)
	})

	c.Visit("https://stomshop.pro/")
	for _, elem := range slice {
		fmt.Println("============================================================", elem.Name, "============================================================")
		fmt.Println("_____ ", elem.Link)
		for _, subElem := range elem.sub {
			fmt.Println("+++ ", subElem.SubName, " +++")
			fmt.Println("____", subElem.SubLink)
			for _, l := range subElem.Products {
				fmt.Println(l.PName, "|", l.PLink, "|", l.Price, "|", l.Number, "|", l.Country)
			}
		}
	}
}
