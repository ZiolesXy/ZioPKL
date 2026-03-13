package respons

type User struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Email string `json:"email,omitempty"`
	Password string `json:"-"`
	Role Role `json:"role,omitempty"`
}

type Profile struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Email string `json:"email,omitempty"`
	Password string `json:"-"`
}

type UserWithBarang struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Role Role `json:"role,omitempty"`
	Email string `json:"email"`
	Password string `json:"-"`
	Barangs []Barang `json:"barangs"`
}

type OwnerPost struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Barangs Barang `json:"barangs"`
}