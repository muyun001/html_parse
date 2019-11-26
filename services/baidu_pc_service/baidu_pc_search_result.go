package baidu_pc_service

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strconv"
	"strings"
)

var protocolPrefixes = []string{
	"http://",
	"https://",
	"ftp://",
	"mailto://",
}

func GetDomain(href string) (string, error) {
	href = formatUrl(href)
	domainUrl, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	return domainUrl.Host, nil
}

func formatUrl(href string) string {
	hasProtocol := false
	for _, protocolPrefix := range protocolPrefixes {
		if strings.HasPrefix(href, protocolPrefix) {
			hasProtocol = true
			break
		}
	}
	if !hasProtocol {
		href = "http://" + href
	}
	return href
}

type SearchResult struct {
	Port                       string   `json:"port"`
	Rank                       int      `json:"rank"`
	BaiduURL                   string   `json:"baidu_url"`
	Title                      string   `json:"title"`
	RealUrl                    string   `json:"real_url"`
	DisplayUrl                 string   `json:"display_url"`
	SiteName                   string   `json:"site_name"`
	Type                       string   `json:"type"` //vid_pocket 视频，
	TitleMatchWords            []string `json:"title_match_words"`
	BaiduDescriptionMatchWords []string `json:"baidu_description_match_words"` //百度显示的description的飘红字
	BaiduDescription           string   `json:"baidu_description"`             // 百度显示的description
	CacheUrl                   string   `json:"cache_url"`
	IsEnterpriseCertificate    bool     `json:"is_enterprise_certificate"` // 是否企业实名认证
	IsAd                       bool     `json:"is_ad"`                     // 是否是广告
}

func (sr *SearchResult) GetPCRealUrl() error {
	if sr.RealUrl == "" {
		// 如果displayUrl可以作为real则不用发送请求了
		if sr.SiteName == "" && sr.DisplayUrl != "" && !strings.Contains(sr.DisplayUrl, "...") {
			if !strings.Contains(sr.DisplayUrl, "http") {
				sr.RealUrl = "http://" + sr.DisplayUrl
			} else {
				sr.RealUrl = sr.DisplayUrl
			}
		} else {
			sr.RealUrl = DecodeBaiduEncURL(sr.BaiduURL)
		}
	}
	return nil
}

func (sr *SearchResult) IsHomePage() bool {
	domain := ""
	if sr.DisplayUrl != "" && !strings.Contains(sr.DisplayUrl, "...") {
		if strings.Index(sr.DisplayUrl, "http") >= 0 {
			domain, _ = GetDomain(sr.DisplayUrl)
		} else {
			domain, _ = GetDomain("http://" + sr.DisplayUrl)
		}
	} else {
		if strings.Contains(sr.DisplayUrl, "...") {
			if strings.Contains(sr.DisplayUrl, "?") {
				return false
			}
			td := strings.Replace(sr.DisplayUrl, "//", "", -1)
			if strings.Contains(td, "/") {
				return false
			}
		}
		_ = sr.GetPCRealUrl()
		if sr.RealUrl != "" {
			domain, _ = GetDomain(sr.RealUrl)
		}
	}
	if sr.RealUrl != "" {
		return strings.HasSuffix(strings.Replace(sr.RealUrl, "/", "", -1), domain)
	} else if sr.DisplayUrl != "" {
		return strings.HasSuffix(strings.Replace(sr.DisplayUrl, "/", "", -1), domain)
	}
	return false

}

const PcPort = "PC"

func ParseBaiduPCSearchResultHtml(html string) (*[]SearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	var results []SearchResult
	doc.Find("div.c-container").Each(func(index int, searchResultElement *goquery.Selection) {
		resItem := SearchResult{Port: PcPort}
		if rank := getBaiduPcSearchResultRank(searchResultElement); rank == 0 {
			return
		} else {
			resItem.Rank = rank
		}

		// title相关
		//fmt.Println(index)
		var titleElement *goquery.Selection
		searchResultElement.Find("h3.t>a").Each(func(_ int, te *goquery.Selection) {
			if href, ok := te.Attr("href"); !ok || href == "" {
				return
			}
			titleElement = te
		})
		if titleElement == nil {
			return
		}
		baiduUrl, ok := titleElement.Attr("href")
		if !ok {
			return
		} else {
			resItem.Title = titleElement.Text()
			titleElement.Find("em").Each(func(_ int, redElement *goquery.Selection) {
				if redElement.Text() != "..." {
					resItem.TitleMatchWords = append(resItem.TitleMatchWords, redElement.Text())
				}
			})
			resItem.BaiduURL = baiduUrl
		}

		// description相关
		abstractElement := searchResultElement.Find(".c-abstract")
		//fmt.Println(goquery.OuterHtml(searchResultElement))
		resItem.BaiduDescription = abstractElement.Text()
		if resItem.BaiduDescription == "" {
			abstractElement = searchResultElement.Find(".op-vmp-zxenterprise-contianer")
			resItem.BaiduDescription = abstractElement.Text()
		}
		abstractElement.Find("em").Each(func(_ int, redElement *goquery.Selection) {
			if redElement.Text() != "..." {
				resItem.BaiduDescriptionMatchWords = append(resItem.BaiduDescriptionMatchWords, redElement.Text())
			}
		})

		// 底部url相关
		displayUrlEle := searchResultElement.Find(".c-showurl")
		if displayUrlEle.Find("style").Size() > 0 {
			resItem.SiteName = displayUrlEle.Find("span").Text()
		} else {
			resItem.DisplayUrl = strings.TrimSpace(displayUrlEle.Text())
		}

		// cacheUrl
		cacheUrlElem := searchResultElement.Find("div.f13 a.m")
		if strings.Contains(cacheUrlElem.Text(), "百度快照") {
			href, exist := cacheUrlElem.Attr("href")
			if exist {
				resItem.CacheUrl = href
			}
		}

		// 企业认证
		enterpriseCertificateItem := searchResultElement.Find("span.c-icon-baozhang-new")
		if len(enterpriseCertificateItem.Nodes) > 0 {
			resItem.IsEnterpriseCertificate = true
		}

		results = append(results, resItem)
	})

	return &results, err
}

// 广告解析
func ParseBaiduPCSearchAdResultHtml(html string) (*[]SearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	var results []SearchResult

	doc.Find("div#content_left div[data-click]").Each(func(index int, searchResultElement *goquery.Selection) {
		resItem := SearchResult{Port: PcPort}
		if rank := getBaiduPcSearchResultRank(searchResultElement); rank == 0 {
			return
		} else {
			resItem.Rank = rank
		}

		// 广告
		if idStr := searchResultElement.AttrOr("id", ""); idStr == "" {
			return
		} else {
			if idInt, err := strconv.Atoi(idStr); err != nil {
				return
			} else {

				if idInt > 2000 {
					resItem.IsAd = true
				}

				if idInt <= 2000 && idInt >= 1000 {
					fmt.Println("-------------------------------------------------------------------")
				}
			}
		}

		// 如果不是ad，就不继续解析
		if !resItem.IsAd {
			return
		}

		// title相关
		//fmt.Println(index)
		var titleElement *goquery.Selection
		searchResultElement.Find("h3.t>a").EachWithBreak(func(i int, te *goquery.Selection) bool {
			if href, ok := te.Attr("href"); !ok || href == "" {
				return false
			}
			titleElement = te
			return false
		})
		if titleElement == nil {
			return
		}
		baiduUrl, ok := titleElement.Attr("href")
		if !ok {
			return
		} else {
			resItem.Title = titleElement.Text()
			titleElement.Find("font").Each(func(_ int, redElement *goquery.Selection) {
				if redElement.Text() != "..." {
					resItem.TitleMatchWords = append(resItem.TitleMatchWords, redElement.Text())
				}
			})
			resItem.BaiduURL = baiduUrl
		}

		// description相关
		searchResultElement.Find("a[hidefocus]").EachWithBreak(func(i int, selection *goquery.Selection) bool {
			resItem.BaiduDescription = selection.Text()
			if resItem.BaiduDescription == "" {
				return true
			}

			selection.Find("font").Each(func(i int, fontSelection *goquery.Selection) {
				resItem.BaiduDescriptionMatchWords = append(resItem.BaiduDescriptionMatchWords, fontSelection.Text())
			})
			return false
		})

		searchResultElement.Find("div a[target] span").EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if selection.Text() != "" {
				resItem.DisplayUrl = selection.Text()
				return false
			}
			return false
		})

		// cacheUrl
		cacheUrlElem := searchResultElement.Find("div.f13 a.m")
		if strings.Contains(cacheUrlElem.Text(), "百度快照") {
			href, exist := cacheUrlElem.Attr("href")
			if exist {
				resItem.CacheUrl = href
			}
		}

		// 企业认证
		enterpriseCertificateItem := searchResultElement.Find("span.icons")
		if len(enterpriseCertificateItem.Nodes) > 0 {
			resItem.IsEnterpriseCertificate = true
		}

		results = append(results, resItem)
	})

	return &results, err
}

func getBaiduPcSearchResultRank(searchResultElement *goquery.Selection) int {
	if idStr, ok := searchResultElement.Attr("id"); ok {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return 0
		}
		return id
	} else {
		return 0
	}
}

// 通过 主域匹配 这种方式是包含的关系 算是模糊匹配
func MatchRankByDomain(srs *[]SearchResult, domain string) (rank int) {
	return MatchRank(srs, domain, "", "", "")
}

func MatchRankByDisplayUrl(srs *[]SearchResult, displayUrl string) (rank int) {
	return MatchRank(srs, "", displayUrl, "", "")
}

// 匹配排名，根据真实url
func MatchRankByReal(srs *[]SearchResult, realUrl string) (rank int) {
	//先获取没有协议的真实url
	realUrlWithoutProtocol := strings.Replace(realUrl, "http://", "", 1)
	realUrlWithoutProtocol = strings.Replace(realUrlWithoutProtocol, "https://", "", 1)

	//第一遍先用display匹配一次
	for i, sr := range *srs {
		// 如果有displayUrl 先和 displayUrl进行匹配
		if sr.DisplayUrl != "" && !strings.Contains(sr.DisplayUrl, "...") {
			// 这里可以直接获取real
			_ = (*srs)[i].GetPCRealUrl()
			sr = (*srs)[i]
			//这里因为 RealUrl 还有可能是https的
			if strings.HasSuffix(sr.RealUrl, realUrlWithoutProtocol) || strings.HasSuffix(sr.RealUrl, realUrlWithoutProtocol+"/") {
				// 百度对于一次搜索结果的url应该具有唯一性， 匹配到就返回
				rank = sr.Rank
				return
			}
		}
	}

	// 这一遍用realUrl匹配了
	for i, sr := range *srs {
		if sr.DisplayUrl != "" || (sr.DisplayUrl == "" && sr.SiteName == "") {
			continue
		}

		// 排除百度系
		// 但是这个不合理,万一real是百度就出错了不过为了减少查询次数，还是加上
		if strings.Contains(sr.DisplayUrl, "baidu.com") || strings.Contains(sr.SiteName, "百度") {
			continue
		}

		if sr.BaiduURL != "" {
			_ = (*srs)[i].GetPCRealUrl()
			sr = (*srs)[i]
			if strings.HasSuffix(sr.RealUrl, realUrlWithoutProtocol) || strings.HasSuffix(sr.RealUrl, realUrlWithoutProtocol+"/") {
				// 百度对于一次搜索结果的url应该具有唯一性， 匹配到就返回
				rank = sr.Rank
				return
			}
		}
	}
	return
}

// 匹配排名 根据多重条件
// domain  displayUrl siteName 属于非强制型匹配， 即匹配不上还会进行其它项匹配
// title 属于强制型匹配 匹配不上则 直接判定匹配不上
func MatchRank(srs *[]SearchResult, domain string, displayUrl string, siteName string, title string) (rank int) {
	for _, sr := range *srs {
		matchFlag := false

		//这里是模糊匹配
		if domain != "" && sr.DisplayUrl != "" {
			if strings.Contains(sr.DisplayUrl, domain) {
				matchFlag = true
			}
		}

		// 这里算是精确匹配了
		if sr.DisplayUrl != "" && displayUrl != "" {
			if strings.HasSuffix(sr.DisplayUrl, displayUrl) || strings.HasSuffix(sr.DisplayUrl, displayUrl+"/") || strings.HasSuffix(sr.DisplayUrl, displayUrl+"...") {
				matchFlag = true
			}
		}

		// 这一条是和displayUrl算是互斥，有siteName 就不太会有display
		if siteName != "" && sr.SiteName != "" {
			if sr.SiteName == siteName {
				matchFlag = true
			}
		}

		// title 算是强制匹配了，如果没匹配上则跳过
		if title != "" && sr.Title != "" {
			// 如果title 是组合匹配中的一项 则需要 其它组合匹配项能狗匹配
			if displayUrl != "" || domain != "" || siteName != "" && matchFlag == false {
				continue
			}
			//取前17个字符匹配
			titlePart := strings.Split(title, "")
			matchTitle := strings.Join(titlePart[0:17], "")
			if strings.HasPrefix(sr.Title, matchTitle) {
				matchFlag = true
			} else {
				matchFlag = false
			}
		}

		// 经过多轮匹配后都过关了，则确定排名
		if matchFlag {
			rank = sr.Rank
			break
		}

	}
	return
}

// 获取一个站的首页位置，一般是配合site使用
func GetFirstHomePageRank(srs *[]SearchResult, domain string) (rank int) {
	return MatchRankByReal(srs, "http://"+domain)
}


func getRedWords(element *goquery.Selection) (words []string) {
	words = []string{}
	element.Find("em").Each(func(_ int, redElement *goquery.Selection) {
		if redElement.Text() != "..." {
			words = append(words, redElement.Text())
		}
	})
	return
}
