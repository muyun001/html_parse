package main

import (
	"encoding/json"
	"fmt"
	"github.com/panwenbin/ghttpclient"
	"html_parse_api/services/baidu_pc_service"
	"html_parse_api/services/so_service"
	"html_parse_api/services/sogou_service"
	"io/ioutil"
	"strings"
)

type ParseRequest struct {
	Url  string `json:"url"`
	Html string `json:"html"`
}

type ParseResultBaidu struct {
	Code int                              `json:"code"`
	Msg  string                           `json:"msg"`
	Data baidu_pc_service.BaiduSearchInfo `json:"data"`
}

type ParseResultSo struct {
	Code int                     `json:"code"`
	Msg  string                  `json:"msg"`
	Data so_service.SoSearchInfo `json:"data"`
}

type ParseResultSogou struct {
	Code int                           `json:"code"`
	Msg  string                        `json:"msg"`
	Data sogou_service.SogouSearchInfo `json:"data"`
}

const baiduPcHtmlFile = "./test/test_html/baidu_pc.html"
const soHtmlFile = "./test/test_html/so.html"
const sogouHtmlFile = "./test/test_html/sogou.html"

func main() {
	var url string
	var html string

	// 测试百度pc
	contentsBaiduPc, _ := ioutil.ReadFile(baiduPcHtmlFile)
	url = "https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=1&rsv_idx=1&tn=baidu&wd=%E5%9E%83%E5%9C%BE%E5%A4%84%E7%90%86"
	html = strings.Replace(string(contentsBaiduPc), "\n", "", 1)
	testBaiduPc(url, html)

	// 测试360
	contentsSo, _ := ioutil.ReadFile(soHtmlFile)
	url = "https://www.so.com/s?q=%E8%8B%8F%E5%B7%9E%E5%BA%9F%E5%93%81%E5%9B%9E%E6%94%B6&pn=4&psid=2539e98f7538f264ce81143c4c7d9f97&src=srp_paging&fr=none"
	html = strings.Replace(string(contentsSo), "\n", "", 1)
	testSo(url, html)

	// 测试搜狗
	contentsSogou, _ := ioutil.ReadFile(sogouHtmlFile)
	url = "https://www.sogou.com/web?query=%E5%9E%83%E5%9C%BE%E5%A4%84%E7%90%86&_ast=1574155635&_asf=www.sogou.com&w=01029901&cid=&s_from=result_up"
	html = strings.Replace(string(contentsSogou), "\n", "", 1)
	testSogou(url, html)
}

func testBaiduPc(searchUrl, html string) {
	url := "http://127.0.0.1:9020/html-parse/baidu-pc"
	requestBody := &ParseRequest{
		Url:  searchUrl,
		Html: html,
	}

	jsonBytes, _ := json.Marshal(requestBody)
	parseResult := &ParseResultBaidu{}
	err := ghttpclient.PostJson(url, jsonBytes, nil).ReadJsonClose(parseResult)
	if err != nil {
		panic(err)
	}

	fmt.Print(parseResult)
}

func testSo(searchUrl, html string) {
	url := "http://127.0.0.1:9020/html-parse/so"
	requestBody := &ParseRequest{
		Url:  searchUrl,
		Html: html,
	}

	jsonBytes, _ := json.Marshal(requestBody)
	parseResult := &ParseResultSo{}
	err := ghttpclient.PostJson(url, jsonBytes, nil).ReadJsonClose(parseResult)
	if err != nil {
		panic(err)
	}

	fmt.Println(parseResult)
}

func testSogou(searchUrl, html string) {
	url := "http://127.0.0.1:9020/html-parse/sogou"
	requestBody := &ParseRequest{
		Url:  searchUrl,
		Html: html,
	}

	jsonBytes, _ := json.Marshal(requestBody)
	parseResult := &ParseResultSogou{}
	err := ghttpclient.PostJson(url, jsonBytes, nil).ReadJsonClose(parseResult)
	if err != nil {
		panic(err)
	}

	fmt.Print(parseResult)
}
