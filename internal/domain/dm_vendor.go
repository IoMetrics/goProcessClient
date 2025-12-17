package domain

// Ex.: info básica do usuário que você quer devolver
type UserInfo struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Level    *string `json:"level"`
}

// Ex.: info básica do usuário que você quer devolver
type DbInfo struct {
	Ip_local string `json:"ip_local"`
	Db_local string `json:"db_local"`
}

type ItemDTO struct {
	ID       int        `json:"id"`
	Product  ProductDTO `json:"product"`
	Options  string     `json:"options"`
	Quantity float64    `json:"quantity"`
}

// Produto “achatado” só com o que o app precisa
type ProductDTO struct {
	Product     string  `json:"product"`
	Description string  `json:"description"`
	Value       float64 `json:"value"`
	Printer     *string `json:"printer"`
	GroupId     string  `json:"group_id"`
	Details     *string `json:"details"`
	Unit        *string `json:"unit"`
}

// Grupo de produto
type ProductGroupDTO struct {
	Id          string `json:"id"`
	Description string `json:"Description"`
}

// Catálogo completo enviado no login
type CatalogResponse struct {
	Products []ProductDTO      `json:"products"`
	Groups   []ProductGroupDTO `json:"groups"`
}

// Resposta final do login:
// - Embede dmAuth.LoginResponse (tokens) para manter compatibilidade
// - Adiciona usuário + catálogo
type FullLoginResponse struct {
	LoginResponse                 // campos de token continuam na raiz do JSON
	User          UserInfo        `json:"user"`
	Catalog       CatalogResponse `json:"catalog"`
}
