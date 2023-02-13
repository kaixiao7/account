package model

type Page struct {
	// 总页数
	Total    int `json:"total"`
	PageSize int `json:"pageSize"`
	PageNum  int `json:"pageNum"`
	Data     any `json:"data"`
}
