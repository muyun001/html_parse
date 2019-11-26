package structs

// 请求内容
type ParseRequest struct {
	Url  string `json:"url"`
	Html string `json:"html"`
}
