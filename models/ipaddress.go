package models

type IpAddressDay struct {
	DayName string `json:"day_name"`
	Total   uint64 `json:"total"`
}

type IpAddressBlockCategoryDay struct {
	CategoryName string `json:"category_name"`
	Total        uint64 `json:"total"`
}
