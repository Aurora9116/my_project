package model

import "github.com/gin-gonic/gin"

type Page struct {
	Page     int64 `json:"page" from:"page"`
	PageSize int64 `json:"page_size" from:"page_size"'`
}

func (p *Page) Bind(c *gin.Context) {
	_ = c.ShouldBind(&p)
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = 10
	}

}
