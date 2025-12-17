package domain

// ErrorResponse é a resposta padrão de erro em JSON.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// SuccessResponse é uma resposta de sucesso mais "rica",
// usada quando você precisa devolver vendedor + config de ambiente.
type SuccessResponse struct {
	Success    bool       `json:"success"`
	Vendedores []UserInfo `json:"vendedores"`
	Modo       string     `json:"ws_modo,omitempty"`
}
