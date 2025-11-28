package domain

// Vendor representa o vendedor/autorizado a usar o app.
type Vendor struct {
	Cod     int    `json:"cod"`
	Nome    string `json:"nome"`
	Usuario string `json:"usuario"`
	Nivel   string `json:"nivel"`
}
