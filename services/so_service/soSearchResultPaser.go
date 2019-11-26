// 对搜狗搜索结果进行分析
package so_service

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

type SoSearchResult struct {
	Port                    string `json:"port"`
	Rank                    int    `json:"rank"`
	SoURL                   string `json:"so_url"`
	Title                   string `json:"title"`
	RealUrl                 string `json:"real_url"`
	DisplayUrl              string `json:"display_url"`
	SiteName                string `json:"site_name"`
	Type                    string `json:"type"` //vid_pocket 视频，
	SoDescription           string `json:"so_description"`
	CacheUrl                string `json:"cache_url"`
	IsEnterpriseCertificate bool   `json:"is_enterprise_certificate"` // 是否企业实名认证
	IsAd                    bool   `json:"is_ad"`                     // 是否是广告
}

// 解析普通数据
func ParseSoSearchResultHtml(html string, pageNum int) (*[]SoSearchResult, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var results []SoSearchResult
	rank := (pageNum - 1) * 10
	dom.Find("ul.result li.res-list").Each(func(i int, selection *goquery.Selection) {
		resItem := SoSearchResult{Port: "Pc"}

		// title和soUrl
		titleItem := selection.Find("h3")
		if title := titleItem.Text(); title != "" {
			title = strings.Replace(title, " ", "", -1)
			title = strings.Replace(title, "\n", "", -1)
			resItem.Title = title
		}

		if href := titleItem.Find("a[href]").AttrOr("href", ""); href != "" {
			resItem.SoURL = href
		}

		rank ++
		resItem.Rank = rank

		// cacheUrl和realUrl
		selection.Find("p.res-linkinfo a").EachWithBreak(func(_ int, subSelection *goquery.Selection) bool {
			if cacheStr := subSelection.Text(); cacheStr == "快照" {
				resItem.CacheUrl = subSelection.AttrOr("href", "")
				_ = resItem.GetSoRealUrl()
				return false
			}
			return true
		})

		// displayUrl和siteName
		displayUrlOrSiteName := selection.Find("cite").Text()
		if displayUrlOrSiteName != "" {
			resItem.DisplayUrl, resItem.SiteName = GetSoDisplayUrlAndSiteName(displayUrlOrSiteName)
		}

		// SoDescription
		if description := selection.Find(".res-desc").Text(); description != "" {
			resItem.SoDescription = strings.Replace(strings.Replace(description, "\n", "", -1), " ", "", -1)
		} else {
			selection.Find(".res-rich div.res-comm-con p").EachWithBreak(func(j int, desSelection *goquery.Selection) bool {
				resItem.SoDescription = strings.Replace(strings.Replace(desSelection.Text(), "\n", "", -1), " ", "", -1)
				return false
			})
		}

		results = append(results, resItem)
	})

	return &results, nil
}

// 解析广告数据
func ParseSoSearchAdResultHtml(html string) (*[]SoSearchResult, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var results []SoSearchResult
	dom.Find("div.spread_test_height li").Each(func(i int, selection *goquery.Selection) {
		resItem := SoSearchResult{Port: "Pc"}
		resItem.IsAd = true

		// title和sogouUrl
		if title := selection.Find("a").Text(); title != "" {
			title = strings.Replace(title, " ", "", -1)
			title = strings.Replace(title, "\n", "", -1)
			resItem.Title = title
		}

		if href := selection.Find("a[href]").AttrOr("href", ""); href != "" {
			resItem.SoURL = href
		}

		// displayUrl和siteName
		displayUrlOrSiteName := selection.Find("div cite a").Text()
		if displayUrlOrSiteName != "" {
			resItem.DisplayUrl, resItem.SiteName = GetSoAdDisplayUrlAndSiteName(displayUrlOrSiteName)
		}

		// SoDescription
		if description := selection.Find("p").Text(); description != "" {
			description = strings.Replace(description, "\n", "", -1)
			description = strings.Replace(description, " ", "", -1)
			resItem.SoDescription = description
		}

		results = append(results, resItem)
	})

	dom.Find("ul#e_idea_pp > li").Each(func(i int, selection *goquery.Selection) {
		resItem := SoSearchResult{Port: "Pc"}
		resItem.IsAd = true

		// title和sogouUrl
		if title := selection.Find("a").Text(); title != "" {
			title = strings.Replace(title, " ", "", -1)
			title = strings.Replace(title, "\n", "", -1)
			resItem.Title = title
		}

		if strings.Contains(resItem.Title, "为您推荐更多优质结果") {
			return
		}

		if href := selection.Find("a[href]").AttrOr("href", ""); href != "" {
			resItem.SoURL = href
		}

		// displayUrl和siteName
		displayUrlOrSiteName := selection.Find("div.e_haoso_fengwu_extend").Text()
		if displayUrlOrSiteName != "" {
			resItem.DisplayUrl, resItem.SiteName = GetSoAdDisplayUrlAndSiteName(displayUrlOrSiteName)
		}

		// SoDescription
		if description := selection.Find(".inner_desc").Text(); description != "" {
			description = strings.Replace(description, "\n", "", -1)
			description = strings.Replace(description, " ", "", -1)
			resItem.SoDescription = description
		}

		results = append(results, resItem)
	})

	return &results, nil
}

// 获取DisplayUrl和SiteName
func GetSoDisplayUrlAndSiteName(str string) (string, string) {
	var displayUrl, siteName = "", ""
	if showUrlInfos := strings.SplitN(str, ">", 2); len(showUrlInfos) > 0 {
		displayUrl = showUrlInfos[0]
	} else {
		displayUrl = str
	}

	displayUrl = strings.Replace(strings.Replace(displayUrl, " ", "", -1), "\n", "", -1)
	siteName = strings.Replace(strings.Replace(siteName, " ", "", -1), "\n", "", -1)
	return displayUrl, siteName
}

// 获取广告信息的DisplayUrl和SiteName
func GetSoAdDisplayUrlAndSiteName(str string) (string, string) {
	var displayUrl, siteName = "", ""
	if showUrlInfos := strings.SplitN(str, "- ", 2); len(showUrlInfos) > 0 {
		if strings.Contains(showUrlInfos[0], ".") {
			displayUrl = showUrlInfos[0]
		} else {
			siteName = showUrlInfos[0]
		}
	} else {
		displayUrl = str
	}

	displayUrl = strings.Replace(strings.Replace(displayUrl, " ", "", -1), "\n", "", -1)
	siteName = strings.Replace(strings.Replace(siteName, " ", "", -1), "\n", "", -1)
	return displayUrl, siteName
}

// 获取RealUrl
func (sosi *SoSearchResult) GetSoRealUrl() error {
	urlStr, err := url.Parse(sosi.CacheUrl)
	if err != nil {
		return err
	}

	sosi.RealUrl = urlStr.Query().Get("u")
	return nil
}
