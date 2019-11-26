package baidu_pc_service_test

import (
	"fmt"
	"html_parse_api/services/baidu_pc_service"
	"io/ioutil"
	"strings"
	"testing"
)

const dataFileBaiduPc = "./test_html/baidu_pc.html"

func TestParseBaiduPcSearchInfo(t *testing.T) {
	contents, err := ioutil.ReadFile(dataFileBaiduPc)
	if err != nil {
		t.Fatal("读取文件错误")
	}

	searchHTML := strings.Replace(string(contents), "\n", "", 1)
	bi, err := baidu_pc_service.ParseBaiduPcSearchInfo(searchHTML)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d, %d\n", len(*bi.SearchResults), len(*bi.SearchAdResults))
	for _, sr := range *bi.SearchAdResults {
		fmt.Printf("%v\n", sr)
	}
	for _, sr := range *bi.SearchResults {
		fmt.Println(sr.IsEnterpriseCertificate)
	}
}
