package so_service_test

import (
	"fmt"
	"html_parse_api/services/so_service"
	"io/ioutil"
	"strings"
	"testing"
)

const dataFileSoPc2 = "./test_html/so_pc.html"

func TestParseSoSearchResultHtml(t *testing.T) {
	contents, err := ioutil.ReadFile(dataFileSoPc2)
	if err != nil {
		t.Fatal("读取文件错误")
	}

	searhHTML := strings.Replace(string(contents), "\n", "", 1)
	sogouPaserResult, err := so_service.ParseSoSearchResultHtml(searhHTML, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println(sogouPaserResult)
}

func TestParseSoSearchAdResultHtml(t *testing.T) {
	contents, err := ioutil.ReadFile(dataFileSoPc2)
	if err != nil {
		t.Fatal("读取文件错误")
	}

	searhHTML := strings.Replace(string(contents), "\n", "", 1)
	sogouPaserAdResult, err := so_service.ParseSoSearchAdResultHtml(searhHTML)
	if err != nil {
		panic(err)
	}

	fmt.Println(sogouPaserAdResult)
}
