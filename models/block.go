package models

type Block struct {
	BlockId           string `json:"block_id"`
	Domain            string `json:"domain"`
	BaseDomain        string `json:"base_domain"`
	IpAddress         string `json:"ip_Address"`
	HasSubdomain      uint64 `json:"has_subdomain"`
	BlockCategoryName string `json:"block_category_name"`
	BlockName         string `json:"block_name"`
	LogDatetime       string `json:"log_datetime"`
	CreatedAt         string `json:"created_at"`
}
