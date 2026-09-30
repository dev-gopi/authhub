package response

import "fmt"

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

func NewPagination(page, pageSize, total int) *Pagination {
	fmt.Printf("Creating pagination: page=%d, size=%d, total=%d\n", page, pageSize, total)
	return &Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
}
