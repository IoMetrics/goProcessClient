package auth

import (
	"encoding/json"
	"fmt"
	cnfRepo "goProcessClient/internal/data/repository"
	dmAuth "goProcessClient/internal/domain"
	jwt "goProcessClient/internal/http/handlers/jwt"
	"net/http"
)

// LoginHandler usa o repository para validar e gerar tokens
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Método não permitido, use POST")
		return
	}

	var req dmAuth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido na requisição")
		return
	}

	if req.WsChave == "" {
		writeError(w, http.StatusBadRequest, "ws_chave é obrigatório")
		return
	}

	// 1) Pega a chave e acha o banco (config do cliente)
	wsConfig, err := cnfRepo.BuscarConfigPorChave(req.WsChave)
	if err != nil {
		writeError(w, http.StatusInternalServerError,
			fmt.Sprintf("Erro ao buscar configuração: %v", err))
		return
	}
	if wsConfig == nil {
		writeError(w, http.StatusUnauthorized, "Chave inválida ou não encontrada")
		return
	}

	// 2) Usa o banco ACHADO pela chave pra validar o usuário
	vendor, err := cnfRepo.BuscarUsuarioPorLogin(wsConfig, req.Usuario, req.Senha)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Usuário ou senha inválidos")
		return
	}

	// 3) Gera tokens JWT
	loginResp, err := jwt.GenerateTokens(*vendor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Erro ao gerar tokens")
		return
	}

	// 4) Responde no formato que o app espera
	writeJSON(w, http.StatusOK, loginResp)
}

// HealthHandler – simples endpoint de health check: GET /health
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
