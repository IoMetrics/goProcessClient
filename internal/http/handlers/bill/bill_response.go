package bill

import (
	"encoding/json"
	respdm "goProcessClient/internal/domain"
	"net/http"
)

// writeJSON escreve qualquer payload em JSON
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// writeError facilita responder erros padronizados
func writeError(w http.ResponseWriter, status int, msg string) {
	resp := respdm.ErrorResponse{
		Success: false,
		Error:   msg,
	}
	writeJSON(w, status, resp)
}
