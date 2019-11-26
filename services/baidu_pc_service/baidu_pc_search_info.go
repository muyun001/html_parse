package baidu_pc_service

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

// 返回内容
type BaiduSearchInfo struct {
	Port              string          `json:"port"`
	BaiduMatchCount   int             `json:"baidu_match_count"`
	MainPageCount     int             `json:"main_page_count"`
	IsEscape          bool            `json:"is_escape"`
	EscapeWord        string          `json:"escape_word"`
	SearchResultCount int             `json:"search_result_count"`
	SearchResults     *[]SearchResult `json:"search_results"`
	SearchAdResults   *[]SearchResult `json:"search_ad_results"`
}

// 解析百度pc
func ParseBaiduPcSearchInfo(html string) (bsi *BaiduSearchInfo, err error) {
	bsi = &BaiduSearchInfo{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	t := doc.Find("div.nums>span.nums_text").Text()
	if t == "" {
		bsi.BaiduMatchCount = -1
	}

	t = strings.Replace(t, "百度为您找到相关结果约", "", -1)
	t = strings.Replace(t, "个", "", -1)
	t = strings.Replace(t, ",", "", -1)
	bsi.BaiduMatchCount, err = strconv.Atoi(t)

	seTip := doc.Find("#super_se_tip").Text()
	if strings.Contains(seTip, "已显示") {
		doc.Find("#super_se_tip strong").Each(func(_ int, strong *goquery.Selection) {
			_, ok := strong.Attr("class")
			if !ok {
				bsi.IsEscape = true
				bsi.EscapeWord = strong.Text()
			}
		})
	}

	srs, err := ParseBaiduPCSearchResultHtml(html)
	if err != nil {
		return
	}
	bsi.SearchResults = srs

	searchAdResults, err := ParseBaiduPCSearchAdResultHtml(html)
	if err != nil {
		return
	}
	bsi.SearchAdResults = searchAdResults

	for _, sr := range *bsi.SearchResults {
		if sr.IsHomePage() {
			bsi.MainPageCount++
		}
	}

	// 解析SearchResultCount:搜索结果数
	doc.Find("div#page").Each(func(i int, selection *goquery.Selection) {
		nowPage, sumPage := 0, 0
		var haveNextPage bool
		if nowPageStr := selection.Find("strong span.pc").Text(); nowPageStr != "" {
			if nowPage, err = strconv.Atoi(nowPageStr); err != nil {
				return
			}
		} else {
			bsi.SearchResultCount = 10 // 只有一页
		}

		selection.Find("a").Each(func(j int, subSelection *goquery.Selection) {
			if pageStr := subSelection.Text(); pageStr != "" {
				if page, err := strconv.Atoi(pageStr); err != nil {
					if strings.Contains(pageStr, "下一页") {
						haveNextPage = true
					} else {
						return
					}
				} else {
					sumPage = page
				}
			}
		})

		maxPage := GetMax(nowPage, sumPage)
		if maxPage == 0 {
			return
		}

		if nowPage >= 10 && haveNextPage {
			bsi.SearchResultCount = 101
			return
		}

		if maxPage >= 10 && haveNextPage {
			if haveNextPage {
				bsi.SearchResultCount = 101
			}
		} else {
			bsi.SearchResultCount = maxPage * 10
		}
	})

	return bsi, nil
}

func GetMax(i, j int) int {
	if i >= j {
		return i
	}
	return j
}
