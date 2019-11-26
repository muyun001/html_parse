// 对搜狗搜索结果进行分析
package sogou_service

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

type SogouSearchResult struct {
	Port                    string `json:"port"`
	Rank                    int    `json:"rank"`
	SogouURL                string `json:"sogou_url"`
	Title                   string `json:"title"`
	RealUrl                 string `json:"real_url"`
	DisplayUrl              string `json:"display_url"`
	SiteName                string `json:"site_name"`
	Type                    string `json:"type"` //vid_pocket 视频，
	SogouDescription        string `json:"sogou_description"`
	CacheUrl                string `json:"cache_url"`
	IsEnterpriseCertificate bool   `json:"is_enterprise_certificate"` // 是否企业实名认证
	IsAd                    bool   `json:"is_ad"`                     // 是否是广告
}

// 解析普通数据
func ParseSogouSearchResultHtml(html string, pageNum int) (*[]SogouSearchResult, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var results []SogouSearchResult
	rank := (pageNum - 1) * 10
	dom.Find("div.results > div").Each(func(i int, selection *goquery.Selection) {
		resItem := SogouSearchResult{Port: "Pc"}

		// title和sogouUrl
		titleItem := selection.Find("h3")
		if titleItem.Text() != "" {
			elseResultStr := selection.Find("div.vrTitle span").Text() // 判断是不是“98%的人还搜了”这种无用数据
			if strings.Contains(elseResultStr, "的人还搜了") {
				return
			}
		}
		title := strings.Replace(titleItem.Text(), " ", "", -1)
		title = strings.Replace(title, "\n", "", -1)
		resItem.Title = title

		if hrefStr := titleItem.Find("a[href]").AttrOr("href", ""); hrefStr != "" {
			resItem.SogouURL = GetSogouSourceUrl(hrefStr)
		}

		rank++
		resItem.Rank = rank

		// cacheUrl和realUrl
		fbItem := selection.Find("div.fb")
		cacheUrl := fbItem.Find("a[href]").AttrOr("href", "")
		if cacheUrl != "" {
			resItem.CacheUrl = cacheUrl
			_ = resItem.GetSogouRealUrl()
		}

		// displayUrl和siteName
		displayUrlOrSiteName := fbItem.Find("cite").Text()
		if displayUrlOrSiteName != "" {
			resItem.DisplayUrl, resItem.SiteName = GetDisplayUrlAndSiteName(displayUrlOrSiteName)
		}

		// SogouDescription
		findStrs := []string{"div.ft", "p.str_info"}
		for _, findStr := range findStrs {
			if description := selection.Find(findStr).Text(); description != "" {
				resItem.SogouDescription = description
				break
			}
		}

		results = append(results, resItem)
	})

	return &results, nil
}

// 解析广告数据
func ParseSogouSearchAdResultHtml(html string) (*[]SogouSearchResult, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var results []SogouSearchResult
	dom.Find("div.sponsored div.biz_sponsor >div.biz_rb").Each(func(i int, selection *goquery.Selection) {
		resItem := SogouSearchResult{Port: "Pc"}
		resItem.IsAd = true

		// title和sogouUrl
		titleItem := selection.Find("h3")
		if title := titleItem.Text(); title != "" {
			title = strings.Replace(title, " ", "", -1)
			title = strings.Replace(title, "\n", "", -1)
			resItem.Title = title
		}

		if hrefStr := titleItem.Find("a[href]").AttrOr("href", ""); hrefStr != "" {
			resItem.SogouURL = GetSogouSourceUrl(hrefStr)
		}

		// displayUrl和siteName
		findStrs := []string{"div.biz_fb", "div.biz_fb cite"}
		for _, findStr := range findStrs {
			displayUrlOrSiteName := selection.Find(findStr).Text()
			if displayUrlOrSiteName != "" {
				resItem.DisplayUrl, resItem.SiteName = GetDisplayUrlAndSiteName(displayUrlOrSiteName)
				break
			}
		}

		// SogouDescription
		findStrs = []string{".biz_ft", ".crown_info_box", "table", ".list_text"}
		for _, findStr := range findStrs {
			if description := selection.Find(findStr).Text(); description != "" {
				resItem.SogouDescription = description
				break
			}
		}

		results = append(results, resItem)
	})

	return &results, nil
}

// 获取DisplayUrl和SiteName
func GetDisplayUrlAndSiteName(str string) (string, string) {
	var displayUrl, siteName = "", ""
	if showUrlInfos := strings.SplitN(str, "-", 3); len(showUrlInfos) > 0 {
		switch len(showUrlInfos) {
		case 1:
			displayUrl = showUrlInfos[0]
		case 2:
			displayUrl = showUrlInfos[0]
		case 3:
			siteName = showUrlInfos[0]
			displayUrl = showUrlInfos[1]

			if !strings.Contains(displayUrl, ".") {
				displayUrl = showUrlInfos[0]
				siteName = ""
			}
		}
	}

	displayUrl = strings.Replace(strings.Replace(displayUrl, " ", "", -1), "\n", "", -1)
	siteName = strings.Replace(strings.Replace(siteName, " ", "", -1), "\n", "", -1)
	return displayUrl, siteName
}

// 获取SourceUrl
func GetSogouSourceUrl(href string) string {
	if href[:4] != "http" {
		return "https://www.sogou.com" + href
	} else {
		return href
	}
}

// 获取RealUrl
func (ssi *SogouSearchResult) GetSogouRealUrl() error {
	urlStr, err := url.Parse(ssi.CacheUrl)
	if err != nil {
		return err
	}
	ssi.RealUrl = urlStr.Query().Get("url")
	return nil
}
