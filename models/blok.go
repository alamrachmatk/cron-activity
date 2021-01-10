package models

type Blok struct {
	BlokId           string `json:"blok_id"`
	Domain           string `json:"domain"`
	BaseDomain       string `json:"base_domain"`
	IpAddress        string `json:"ip_Address"`
	HasSubdomain     uint64 `json:"has_subdomain"`
	BlokCategoryName string `json:"blok_category_name"`
	BlokName         string `json:"blok_name"`
	LogDatetime      string `json:"log_datetime"`
	CreatedAt        string `json:"created_at"`
}
