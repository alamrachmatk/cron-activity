package models

type DnsDay struct {
	DayName string `json:"day_name"`
	Total   uint64 `json:"total"`
}

type DnsBlockCategoryDay struct {
	CategoryName string `json:"category_name"`
	Total        uint64 `json:"total"`
}
