package baidu_pc_service_test

import (
	"fmt"
	"html_parse_api/services/baidu_pc_service"
	"io/ioutil"
	"strings"
	"testing"
)

const dataFileBaiduPc2 = "./test_html/baidu_pc_action.html"

func TestParseBaiduPCSearchResultHtml(t *testing.T) {
	contents, err := ioutil.ReadFile(dataFileBaiduPc2)
	if err != nil {
		t.Fatal("读取文件错误")
	}
	searchHTML := strings.Replace(string(contents), "\n", "", 1)
	br, err := baidu_pc_service.ParseBaiduPCSearchResultHtml(searchHTML)
	if err != nil {
		panic(err)
	}

	fmt.Println(br)
}

func TestParseBaiduPCSearchAdResultHtml(t *testing.T) {
	contents, err := ioutil.ReadFile(dataFileBaiduPc2)
	if err != nil {
		t.Fatal("读取文件错误")
	}
	searchHTML := strings.Replace(string(contents), "\n", "", 1)
	br, err := baidu_pc_service.ParseBaiduPCSearchAdResultHtml(searchHTML)
	if err != nil {
		panic(err)
	}

	fmt.Println(br)
}
