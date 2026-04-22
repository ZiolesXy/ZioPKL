package dto

type FileURLResponse struct {
	FileURL string `json:"file_url"`
}

type HealthResponse struct {
	Service string `json:"service"`
}
