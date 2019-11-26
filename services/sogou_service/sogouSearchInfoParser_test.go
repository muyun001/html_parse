package sogou_service_test

import (
	"fmt"
	"html_parse_api/services/sogou_service"
	"io/ioutil"
	"strings"
	"testing"
)

const dataFileSogou = "./test_html/sogou_pc.html"

func TestParseSogouSearchInfoFromHtml(t *testing.T) {
	contents, err := ioutil.ReadFile(dataFileSogou)
	if err != nil {
		t.Fatal("读取文件错误")
	}

	searchHTML := strings.Replace(string(contents), "\n", "", 1)
	url := "https://www.sogou.com/web?query=%E5%9E%83%E5%9C%BE%E5%A4%84%E7%90%86&cid=&s_from=result_up&sut=7673&sst0=1574154701660&lkt=0%2C0%2C0&sugsuv=1568950571589337&sugtime=1574154701660&page=4&ie=utf8&w=01029901&dr=1"
	soi, err := sogou_service.ParseSogouSearchInfoFromHtml(searchHTML, url)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(soi)
}
