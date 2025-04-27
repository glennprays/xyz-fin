package model

type Consumer struct {
	NIK            string `json:"nik"`
	PhoneNumber    string `json:"phone_number"`
	PasswordHash   string `json:"-"`
	FullName       string `json:"full_name"`
	LegalName      string `json:"legal_name"`
	TempatLahir    string `json:"tempat_lahir"`
	TanggalLahir   string `json:"tanggal_lahir"`
	Gaji           string `json:"gaji"`
	FotoKTPPath    string `json:"foto_ktp_path"`
	FotoSelfiePath string `json:"foto_selfie_path"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type ConsumerResponse struct {
	NIK          string `json:"nik"`
	PhoneNumber  string `json:"phone_number"`
	FullName     string `json:"full_name"`
	LegalName    string `json:"legal_name"`
	TempatLahir  string `json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	Gaji         string `json:"gaji"`
	CreatedAt    string `json:"created_at"`
}
