package so_service_test

import (
	"fmt"
	"html_parse_api/services/so_service"
	"io/ioutil"
	"strings"
	"testing"
)

const dataFileSo = "./test_html/so_pc.html"

func TestParseSoSearchInfoFromHtml(t *testing.T) {
	contents, err := ioutil.ReadFile(dataFileSo)
	if err != nil {
		t.Fatal("读取文件错误")
	}

	searchHTML := strings.Replace(string(contents), "\n", "", 1)
	url := "https://www.so.com/s?q=%E8%8B%8F%E5%B7%9E%E5%BA%9F%E5%93%81%E5%9B%9E%E6%94%B6&pn=4&psid=2539e98f7538f264ce81143c4c7d9f97&src=srp_paging&fr=none"
	soi, err := so_service.ParseSoSearchInfoFromHtml(searchHTML, url)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(soi)
}
