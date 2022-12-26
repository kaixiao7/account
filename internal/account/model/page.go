package model

type Page struct {
	Total    int `json:"total"`
	PageSize int `json:"pageSize"`
	PageNum  int `json:"pageNum"`
	Data     any `json:"data"`
}
