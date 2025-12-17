package order

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	repo "goProcessClient/internal/data/repository"
	dmOrder "goProcessClient/internal/domain/order"
)

func SendOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Método não permitido, use POST")
		return
	}

	// proteção básica contra payload gigante
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2MB

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Erro ao ler body")
		return
	}
	fmt.Println("Body recebido (order):", string(bodyBytes))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req dmOrder.SendOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido na requisição")
		return
	}

	// validações mínimas
	if req.Comanda == "" {
		writeError(w, http.StatusBadRequest, "comanda é obrigatória")
		return
	}
	if req.Vendedor == "" {
		writeError(w, http.StatusBadRequest, "vendedor é obrigatório")
		return
	}
	if len(req.Items) == 0 {
		writeError(w, http.StatusBadRequest, "items não pode ser vazio")
		return
	}
	for i, it := range req.Items {
		if it.Product == "" && it.ProductDto.Product == "" {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("items[%d].product é obrigatório", i))
			return
		}
		if it.Quantity <= 0 {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("items[%d].quantity deve ser > 0", i))
			return
		}
		if it.UnitPrice < 0 {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("items[%d].unit_price deve ser >= 0", i))
			return
		}
	}

	// ✅ Tenta gravar no banco (vendastemp)
	if err := repo.SaveOrderToVendastemp(r.Context(), req); err != nil {
		// ❗ fallback: grava na pasta_erro com last_error
		errorFile, _ := repo.SaveOrderErrorFile(opts.ErrorDir, req, err)

		ack := dmOrder.AckResponse{
			Status:     "warning",
			Message:    "pedido recebido, mas falhou ao gravar no banco; salvo na pasta_erro",
			ReceivedAt: time.Now().Format(time.RFC3339),
			File:       errorFile,
		}
		writeJSON(w, http.StatusAccepted, ack) // 202 = recebido, mas não finalizado
		return
	}

	ack := dmOrder.AckResponse{
		Status:     "ok",
		Message:    "pedido recebido e gravado na vendastemp",
		ReceivedAt: time.Now().Format(time.RFC3339),
		File:       "",
	}
	writeJSON(w, http.StatusOK, ack)
}
