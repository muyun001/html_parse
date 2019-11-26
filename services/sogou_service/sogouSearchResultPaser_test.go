package sogou_service_test

import (
	"fmt"
	"github.com/kevin-zx/baidu-seo-tool/search/sogou"
	"html_parse_api/services/sogou_service"
	"io/ioutil"
	"strings"
	"testing"
)

const dataFileSogouPc = "./test_html/sogou_pc.html"

func TestParseSogouSearchResultHtml(t *testing.T) {
	contents, err := ioutil.ReadFile(dataFileSogouPc)
	if err != nil {
		t.Fatal("读取文件错误")
	}

	searchHTML := strings.Replace(string(contents), "\n", "", 1)
	sogouParseResult, err := sogou_service.ParseSogouSearchResultHtml(searchHTML, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println(sogouParseResult)
}

func TestParseSogouSearchAdResultHtml(t *testing.T) {
	contents, err := ioutil.ReadFile(dataFileSogouPc)
	if err != nil {
		t.Fatal("读取文件错误")
	}

	searchHTML := strings.Replace(string(contents), "\n", "", 1)
	sogouParseAdResult, err := sogou_service.ParseSogouSearchAdResultHtml(searchHTML)
	if err != nil {
		panic(err)
	}

	fmt.Println(sogouParseAdResult)
}
