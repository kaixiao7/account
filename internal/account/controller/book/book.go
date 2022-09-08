package book

import (
	"kaixiao7/account/internal/account/service"
)

type BookController struct {
	bookSrv service.BookSrv
}

func NewBookContorller() *BookController {
	return &BookController{bookSrv: service.NewBookSrv()}
}
