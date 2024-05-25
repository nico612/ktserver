package request

type Empty struct{}

type PageInfo struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Keyword  string `json:"keyword" form:"keyword"` //关键字
}

type Pagination struct {
	Page     int   `json:"page" form:"page"`         //页码
	PageSize int   `json:"pageSize" form:"pageSize"` //每页数量
	Offset   int   `json:"-"`
	Limit    int   `json:"-"`
	Total    int64 `json:"total"` //总数 作为参数不用填写
}

func (p *Pagination) Check() {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.PageSize == 0 {
		p.PageSize = 10
	}

	//if p.PageSize > 100 {
	//	p.PageSize = 100
	//}

	p.Offset = p.PageSize * (p.Page - 1)
	p.Limit = p.PageSize
}
