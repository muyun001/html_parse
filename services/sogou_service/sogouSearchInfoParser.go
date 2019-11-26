package sogou_service

import (
	"github.com/PuerkitoBio/goquery"
	"html_parse_api/utils"
	"strconv"
	"strings"
)

type SogouSearchInfo struct {
	Port                 string
	SogouMatchCount      int
	MainPageCount        int
	IsEscape             bool
	EscapeWord           string
	SogouSearchResults   *[]SogouSearchResult
	SogouSearchAdResults *[]SogouSearchResult
}

// 结果解析
func ParseSogouSearchInfoFromHtml(html string, url string) (ssi *SogouSearchInfo, err error) {
	ssi = &SogouSearchInfo{}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}

	numStr := dom.Find("div.search-info p.num-tips").Text()
	if numStr == "" {
		ssi.SogouMatchCount = -1
	}

	numStr = strings.Replace(numStr, "搜狗已为您找到约", "", -1)
	numStr = strings.Replace(numStr, "条相关结果", "", -1)
	numStr = strings.Replace(numStr, "投资有风险，选择需谨慎。", "", -1)
	numStr = strings.Replace(numStr, ",", "", -1)
	if numStr != "" {
		ssi.SogouMatchCount, err = strconv.Atoi(numStr)
		if err != nil {
			return
		}
	}

	pageNum, err := utils.ParsePageFromUrl(url, "page")
	if err != nil {
		return
	}

	searchResults, err := ParseSogouSearchResultHtml(html, pageNum)
	if err != nil {
		return
	}
	ssi.SogouSearchResults = searchResults

	searchAdResults, err := ParseSogouSearchAdResultHtml(html)
	if err != nil {
		return
	}
	ssi.SogouSearchAdResults = searchAdResults

	return ssi, nil
}
