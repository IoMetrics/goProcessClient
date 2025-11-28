package domain

// ErrorResponse é a resposta padrão de erro em JSON.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// SuccessResponse é uma resposta de sucesso mais "rica",
// usada quando você precisa devolver vendedor + config de ambiente.
type SuccessResponse struct {
	Success    bool      `json:"success"`
	Vendedores []Vendor  `json:"vendedores"`
	Modo       string    `json:"ws_modo,omitempty"`
	Empresa    string    `json:"ws_empresa,omitempty"`
	IP         string    `json:"ip,omitempty"`
	Port       int       `json:"port,omitempty"`
	Databases  DBNames   `json:"databases,omitempty"`
	Pix        PixConfig `json:"pix,omitempty"`
	FTP        FTPConfig `json:"ftp,omitempty"`
}
