package bill

import (
	"encoding/json"
	"fmt"
	cnfRepo "goProcessClient/internal/data/repository"
	"goProcessClient/internal/domain"
	"net/http"
)

// BillHandler: POST /bill
func BillHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Método não permitido, use POST")
		return
	}

	var req domain.BillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido na requisição")
		return
	}

	// Agora a validação é em cima de IP + banco + comanda
	if req.Comanda == "" {
		writeError(w, http.StatusBadRequest,
			"server_ip, local_db e comanda são obrigatórios")
		return
	}

	// 1) Buscar conta no banco local (vendastemp) usando IP + banco + comanda
	//
	// Você vai implementar essa função no repositório.
	// A ideia é que ela:
	//  - monte o DSN usando req.ServerIP e req.LocalDB
	//  - abra o MySQL local
	//  - leia a tabela vendastemp filtrando por comanda
	//  - monte um domain.BillData e devolva.
	billData, err := cnfRepo.BuscarConta(req.Comanda)
	if err != nil {
		writeError(w, http.StatusInternalServerError,
			fmt.Sprintf("Erro ao buscar conta: %v", err))
		return
	}

	// 2) Responde para o app com a conta real
	writeJSON(w, http.StatusOK, billData)
}
