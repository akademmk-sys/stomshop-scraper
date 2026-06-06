package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Product struct {
	PName   string
	PLink   string
	Price   string
	Number  string
	Country string
}

type Category struct {
	Name   string
	Link   string
	SubCat []SubCategory
}

func (p *Category) HeaderWriter(c, l string, s []SubCategory, slice []Category) []Category {

	p.Name, p.Link = c, l
	slice = append(slice, Category{Name: c, Link: l, SubCat: s})
	return slice
}

type SubCategory struct {
	SubName  string
	SubLink  string
	Products []Product
}

func (s *SubCategory) subHeaderWriter(c, l string, slice []SubCategory) []SubCategory {

	s.SubName, s.SubLink = c, l
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

	AllData := make([]Category, 0)

	c.OnHTML("div.mycategory.product-thumb", func(r *colly.HTMLElement) {
		var h Category
		var d SubCategory
		SubHaders := make([]SubCategory, 0)

		r.ForEach("div.child ul.list-unstyled li a", func(_ int, a *colly.HTMLElement) {
			if a.Text == "показать все ...." {
				return
			}
			SubHaders = d.subHeaderWriter(a.Text, a.Attr("href"), SubHaders)

		})

		r.ForEach("div.child h4 a", func(_ int, a *colly.HTMLElement) {

			AllData = h.HeaderWriter(a.Text, a.Attr("href"), SubHaders, AllData)
		})

	})

	c.OnHTML("div.product-layout div.product-thumb", func(r *colly.HTMLElement) {
		currentUrl := r.Request.URL.String()
		unit := Product{
			PName:   r.ChildText("div.caption a"),
			PLink:   r.ChildAttr("div.caption a", "href"),
			Price:   r.ChildText("span.price-new b"),
			Number:  r.ChildText("span.code span"),
			Country: r.ChildText("div.manufacturer a"),
		}
		for i := range AllData {
			for j := range AllData[i].SubCat {
				if AllData[i].SubCat[j].SubLink == currentUrl {
					AllData[i].SubCat[j].Products = append(AllData[i].SubCat[j].Products, unit)
				}

			}
		}

	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Ошибка:", r.StatusCode, err)
	})

	c.Visit("https://stomshop.pro/")
	for i := range AllData {
		for j := range AllData[i].SubCat {
			c.Visit(AllData[i].SubCat[j].SubLink)
		}
	}
	for _, elem := range AllData {
		fmt.Println("============================================================", elem.Name, "============================================================")
		fmt.Println("_____ ", elem.Link)
		for _, subElem := range elem.SubCat {
			fmt.Println("+++ ", subElem.SubName, " +++")
			fmt.Println("____", subElem.SubLink)
			for _, l := range subElem.Products {
				fmt.Println(l.PName, "|", l.PLink, "|", l.Price, "|", l.Number, "|", l.Country)
			}
		}
	}
}
