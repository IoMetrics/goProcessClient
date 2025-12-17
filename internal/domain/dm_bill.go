package domain

// {
//   "server_ip": "192.168.1.50:3306",
//   "local_db":  "AUTOCOM_LOJA01",
//   "comanda":   "0012"
// }
type BillRequest struct {
	Comanda string `json:"comanda"` // número da comanda
}

// BillItem -> deve bater 1:1 com o modelo Kotlin BillItem

type BillItem struct {
	ID            int     `json:"id"`             // vendastemp.id
	ProductCode   string  `json:"product"`        // vendastemp.produto
	Quantity      float64 `json:"quantity"`       // vendastemp.qte
	UnitPrice     float64 `json:"unit_price"`     // vendastemp.valor (unitário ou total, você decide)
	Discount      float64 `json:"discount"`       // vendastemp.desconto_item
	Commission    float64 `json:"commission"`     // vendastemp.comissao
	CommissionSrv float64 `json:"commission_srv"` // vendastemp.comissaosrv
	ObsItem       string  `json:"obs_item"`       // vendastemp.obs_item (detalhes do pedido)
}

// BillData -> deve bater com o data class BillData do app
type BillData struct {
	Comanda           string     `json:"comanda"`
	Date              string     `json:"date"`              // vendastemp.data
	Time              string     `json:"time"`              // vendastemp.hora
	HourEnd           *string    `json:"hour_end"`          // vendastemp.horafim
	ClosedAt          *string    `json:"closed_at"`         // vendastemp.fechada
	People            int        `json:"people"`            // vendastemp.pessoas
	Waiter            string     `json:"waiter"`            // vendastemp.vendedor/garcom
	Obs               string     `json:"obs"`               // vendastemp.obs
	Items             []BillItem `json:"items"`             // itens da vendastemp
	ServiceTaxPercent float64    `json:"serviceTaxPercent"` // taxa de serviço (%)
}
