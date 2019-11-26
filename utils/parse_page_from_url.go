package utils

import (
	"net/url"
	"strconv"
)

// 从url中解析页码
func ParsePageFromUrl(parseUrl, parameter string) (int, error) {
	urlStr, err := url.Parse(parseUrl)
	if err != nil {
		return 0, err
	}

	pageNum := urlStr.Query().Get(parameter)
	if pageNum == "" {
		return 1, nil
	}

	return strconv.Atoi(pageNum)
}
