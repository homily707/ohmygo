package crawler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
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
		stars := GetStars(r)
		if stars > 0 {
			fmt.Println(r + ": " + fmt.Sprint(stars))
		}
	}
}

func GetStars(r string) int {
	ss := strings.Split(r, "/")
	if len(ss) < 3 {
		return 0
	}
	c := colly.NewCollector()
	c.OnXML()
	resp, err := http.Get("https://api.github.com/repos/" + ss[0] + "/" + ss[1])
	if err != nil {
		fmt.Println(err)
		return 0
	}
	bs, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	stars := gjson.GetBytes(bs, "stargazers_count").Int()
	return int(stars)
}
