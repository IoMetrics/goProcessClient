package order

import (
	configdm "goProcessClient/internal/domain"
)

type SendOrderRequest struct {
	Comanda  string      `json:"comanda"`
	Items    []OrderItem `json:"items"`
	LocalDB  string      `json:"local_db"`
	ServerIP string      `json:"server_ip"`
	Pdv      string      `json:"pdv"`
	Vendedor string      `json:"vendedor"`
}

type OrderItem struct {
	ProductDto configdm.ProductDTO `json:"product_dto"`
	ObsItem    string              `json:"obs_item"`
	Product    string              `json:"product"`
	Quantity   float64             `json:"quantity"`
	UnitPrice  float64             `json:"unit_price"`
}

type AckResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	ReceivedAt string `json:"received_at"`
	File       string `json:"file,omitempty"`
}
