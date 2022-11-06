package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func main() {
	rangeWarp(CountUrl("https://market.kz/astana/elektronika/computery/noutbuki/")) //Сюда url
}

type Cookie struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Url   string `json:"url"`
}

var ParsCookie []Cookie

func Pars(u string) {
	c := colly.NewCollector()

	c.OnHTML(".a-card__content", func(e *colly.HTMLElement) {
		ParsCookie = append(ParsCookie, Cookie{
			Name:  e.ChildText(".a-card__link"),
			Price: e.ChildText(".a-card__price"),
			Url:   e.ChildAttr("a", "href")})
	})

	err := c.Visit(u)
	if err != nil {
		fmt.Println("Нету ссылки в scrape")
	}
}

func Save() {
	ex := excelize.NewFile()

	for i, e := range ParsCookie {
		_ = ex.SetCellValue("Sheet1", fmt.Sprintf("A%v", i), e.Name)
		_ = ex.SetCellValue("Sheet1", fmt.Sprintf("B%v", i), e.Price)
		_ = ex.SetCellValue("Sheet1", fmt.Sprintf("C%v", i), e.Url)
	}

	if err := ex.SaveAs("Cookie.xlsx"); err != nil {
		fmt.Println("Ошибка в сохранение")
	}
}

func CountUrl(url string) (int, string) {
	var strCount string

	v := colly.NewCollector()

	v.OnHTML("div.pagination > ul > li", func(h *colly.HTMLElement) {
		strCount = h.Text
	})

	err := v.Visit(url)
	if err != nil {
		fmt.Println("Нету ссылки в поиске количество url")

	}

	last, _ := strconv.Atoi(strCount)

	return last, url
}

func rangeWarp(page int, url string) {
	fmt.Println("Start")
	for i := 1; i <= page; i++ {
		Pars(url + fmt.Sprintf("?page=%v", i))
		fmt.Printf("Parsing pages № %v\n", i)
	}
	Save()
	fmt.Println("Excel Save")
}
