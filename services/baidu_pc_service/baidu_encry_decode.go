package baidu_pc_service

import (
	"github.com/kevin-zx/http-util"
	"strings"
	"time"
)

func DecodeBaiduEncURL(baiduUrl string) string {

	response, err := http_util.SendRequest(baiduUrl, map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36"}, "GET", nil, 10*time.Second)
	if err != nil {
		return ""
	}

	if response.StatusCode == 200 {

		h, _ := http_util.ReadContentFromResponse(response, "")
		if strings.Contains(h, "window.opener&&window.opener.bds&&window.opener.bds.pdc&&window.opener.bds.pdc.sendLinkLog") {

			part1 := strings.Split(h, "window.location.replace(\"")
			if len(part1) < 2 {
				//fmt.Println(baiduUrl)
				return baiduUrl
			} else {
				return strings.Split(part1[1], "\")},timeout")[0]
			}
		} else if strings.Contains(h, `JSON.parse(localStorage.getItem("tc_time_log")`) {
			ps := strings.Split(h, "\n")
			for _, p := range ps {
				if strings.Contains(p, "window.location.replace(") && strings.Contains(p, ")") {
					start := strings.Index(p, `("`)
					end := strings.LastIndex(p, `")`)
					if end > start+1 && (start > 0 && end > 0) {
						//u :=strings.Replace(p[start+1:end],`"`,"",-1)
						//u =
						return p[start+2 : end]
					}
				}
			}
		} else {
			return response.Request.URL.String()
		}
	}

	return baiduUrl
}
