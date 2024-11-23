package crawler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func ImportedByRepo(url string) {
	c := colly.NewCollector()

	// c.OnHTML("#main-content article div ul details", func(e *colly.HTMLElement) {
	// 	// 遍历所有 <details> 下的 <summary> 元素
	// 	summary := e.ChildText("summary")
	// 	fmt.Printf("Details: %s\n", summary)
	// })
	repo := []string{}
	c.OnXML(`/html/body/main/article/div/ul/details/summary`, func(e *colly.XMLElement) {
		text := e.Text
		if strings.HasPrefix(text, "github") {
			tail := strings.LastIndex(text, "/")
			repo = append(repo, text[:tail])
		}
	})

	c.OnXML(`/html/body/main/article/div/ul/li/a`, func(e *colly.XMLElement) {
		text := e.Text
		if strings.HasPrefix(text, "github") {
			repo = append(repo, text)
		}
	})

	c.Visit(url)

	for i, r := range repo {
		if i%100 == 0 {
			fmt.Println(i)
		}
		stars, _ := GetStars(r)
		if stars > 0 {
			fmt.Println(r + ": " + fmt.Sprint(stars))
		}
	}
}

func GetStars(r string) (int, error) {
	ss := strings.Split(r, "/")
	if len(ss) < 3 {
		return 0, nil
	}
	star := 0
	c := colly.NewCollector()
	c.OnXML(`//*[@id="repo-stars-counter-star"]`, func(e *colly.XMLElement) {
		star, _ = strconv.Atoi(e.Attr("title"))
	})
	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Print(r.Headers)
	// })
	url := "https://github.com/" + ss[1] + "/" + ss[2]
	err := c.Visit(url)
	if err != nil {
		return 0, err
	}
	c.Wait()
	return star, nil
}
