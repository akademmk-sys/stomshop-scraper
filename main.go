package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/xuri/excelize/v2"
)

type Product struct {
	PName       string
	PLink       string
	ImgLink     string
	Price       string
	Number      string
	Country     string
	Description string
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

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       3 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	})

	AllData := make([]Category, 0)

	c.OnHTML("div.mycategory.product-thumb", func(r *colly.HTMLElement) {

		SubHaders := make([]SubCategory, 0)

		r.ForEach("div.child ul.list-unstyled li a", func(_ int, a *colly.HTMLElement) {
			var d SubCategory
			if a.Text == "показать все ...." {
				return
			}
			SubHaders = d.subHeaderWriter(a.Text, a.Attr("href"), SubHaders)
		})

		r.ForEach("div.child h4 a", func(_ int, a *colly.HTMLElement) {
			var h Category
			AllData = h.HeaderWriter(a.Text, a.Attr("href"), SubHaders, AllData)
		})
	})

	c.OnHTML("div.row.products div.product-layout div.product-thumb", func(r *colly.HTMLElement) {
		currentURL := r.Request.URL
		baseSupCategryURl := currentURL.Scheme + "://" + currentURL.Host + currentURL.Path
		price := r.ChildText("span.price-new b")
		if price == "" {
			price = r.ChildText("span.price.text-muted")
		}
		unit := Product{
			PName:   r.ChildText("div.caption a"),
			PLink:   r.ChildAttr("div.caption a", "href"),
			ImgLink: r.ChildAttr("div.image a img", "src"),
			Price:   price,
			Number:  r.ChildText("span.code span"),
			Country: r.ChildText("div.manufacturer a"),
		}
		for i := range AllData {
			for j := range AllData[i].SubCat {
				if AllData[i].SubCat[j].SubLink == baseSupCategryURl {
					AllData[i].SubCat[j].Products = append(AllData[i].SubCat[j].Products, unit)
				}

			}
		}
	})
	c.OnHTML("#tab-description > div", func(r *colly.HTMLElement) {
		fmt.Println("ОПИСАНИЕ СРАБОТАЛО:", r.Text)
		currentURL := r.Request.URL
		ProductURl := currentURL.Scheme + "://" + currentURL.Host + currentURL.Path
		html := r.DOM.Text()
		unit := Product{
			Description: html,
		}
		for i := range AllData {
			for j := range AllData[i].SubCat {
				for l := range AllData[i].SubCat[j].Products {
					if AllData[i].SubCat[j].Products[l].PLink == ProductURl {
						AllData[i].SubCat[j].Products[l].Description = unit.Description
					}
				}
			}
		}

	})

	c.OnHTML("ul.pagination li.active + li a", func(h *colly.HTMLElement) {
		c.Visit(h.Attr("href"))
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
	// c.Visit(AllData[0].SubCat[0].Products[0].PLink)

	for i := range AllData {
		for j := range AllData[i].SubCat {
			for l := range AllData[i].SubCat[j].Products {
				c.Visit(AllData[i].SubCat[j].Products[l].PLink)
			}

		}
	}

	// for _, elem := range AllData {
	// 	fmt.Println("============================================================", elem.Name, "============================================================")
	// 	fmt.Println("_____ ", elem.Link)
	// 	for _, subElem := range elem.SubCat {
	// 		fmt.Println("+++ ", subElem.SubName, " +++")
	// 		fmt.Println("____", subElem.SubLink)
	// 		for _, l := range subElem.Products {
	// 			fmt.Println(l.PName, "|", l.PLink, "|", l.Price, "|", l.Number, "|", l.Country, "|", l.Description)
	// 		}
	// 	}
	// }
	t := excelize.NewFile()
	defer func() {
		if err := t.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	t.SetCellValue("Sheet1", "A1", "Категория")
	t.SetCellValue("Sheet1", "B1", "Ссылка на категрию")
	t.SetCellValue("Sheet1", "C1", "Подкатегория")
	t.SetCellValue("Sheet1", "D1", "Ссылка на подкатегрию")
	t.SetCellValue("Sheet1", "E1", "Товар")
	t.SetCellValue("Sheet1", "F1", "Ссылка на товар")
	t.SetCellValue("Sheet1", "G1", "Ссылка на картинку")
	t.SetCellValue("Sheet1", "H1", "Цена")
	t.SetCellValue("Sheet1", "I1", "Артикул")
	t.SetCellValue("Sheet1", "J1", "Производитель")
	t.SetCellValue("Sheet1", "K1", "Описание")
	row := 2
	for i := range AllData {
		for j := range AllData[i].SubCat {
			for l := range AllData[i].SubCat[j].Products {
				t.SetCellValue("Sheet1", "A"+strconv.Itoa(row), AllData[i].Name)
				t.SetCellValue("Sheet1", "B"+strconv.Itoa(row), AllData[i].Link)
				t.SetCellValue("Sheet1", "C"+strconv.Itoa(row), AllData[i].SubCat[j].SubName)
				t.SetCellValue("Sheet1", "D"+strconv.Itoa(row), AllData[i].SubCat[j].SubLink)
				t.SetCellValue("Sheet1", "E"+strconv.Itoa(row), AllData[i].SubCat[j].Products[l].PName)
				t.SetCellValue("Sheet1", "F"+strconv.Itoa(row), AllData[i].SubCat[j].Products[l].PLink)
				t.SetCellValue("Sheet1", "G"+strconv.Itoa(row), AllData[i].SubCat[j].Products[l].ImgLink)
				t.SetCellValue("Sheet1", "H"+strconv.Itoa(row), AllData[i].SubCat[j].Products[l].Price)
				t.SetCellValue("Sheet1", "I"+strconv.Itoa(row), AllData[i].SubCat[j].Products[l].Number)
				t.SetCellValue("Sheet1", "J"+strconv.Itoa(row), AllData[i].SubCat[j].Products[l].Country)
				t.SetCellValue("Sheet1", "K"+strconv.Itoa(row), AllData[i].SubCat[j].Products[l].Description)
				row++
			}
		}
	}
	t.SaveAs("testcatalog.xlsx")
}
