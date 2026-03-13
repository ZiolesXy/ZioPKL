package request

type BarangPost struct {
	Name string `json:"name" binding:"required"`
}

type BarangPut struct {
	Name *string `json:"name"`
}
