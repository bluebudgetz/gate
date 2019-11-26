package util

type Paging struct {
	Page     uint `query:"_page"`
	PageSize uint `query:"_pageSize"`
}
