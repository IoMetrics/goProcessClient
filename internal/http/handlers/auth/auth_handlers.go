package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	cnfRepo "goProcessClient/internal/data/repository"
	dmAuth "goProcessClient/internal/domain"
	jwt "goProcessClient/internal/http/handlers/jwt"
	"io"
	"net/http"
)

// LoginHandler usa o repository para validar e gerar tokens
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Método não permitido, use POST")
		return
	}

	// 1) Lê o body bruto
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Erro ao ler body")
		return
	}

	// 2) Loga / printa o body como string
	fmt.Println("Body recebido:", string(bodyBytes))

	// 3) Reconstroi o r.Body para poder usar o Decoder depois
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req dmAuth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido na requisição")
		return
	}

	if req.WsChave == "" {
		writeError(w, http.StatusBadRequest, "ws_chave é obrigatório")
		return
	}

	// 2) Usa o banco ACHADO pela chave pra validar o usuário
	vendor, err := cnfRepo.BuscarUsuarioPorLogin(req.Usuario, req.Senha)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Usuário ou senha inválidos")
		return
	}

	// 3) Gera tokens JWT (estrutura antiga que o app já conhece)
	loginResp, err := jwt.GenerateTokens(*vendor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Erro ao gerar tokens")
		return
	}

	// 4) Carrega catálogo de produtos e grupos usando o mesmo banco do cliente
	produtos, err := cnfRepo.BuscarProdutos()
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Erro ao buscar produtos: %v", err),
		)
		return
	}

	grupos, err := cnfRepo.BuscarGruposProdutos()
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Erro ao buscar grupos de produtos: %v", err),
		)
		return
	}

	// 5) Monta a resposta completa (tokens + usuário + catálogo)
	fullResp := dmAuth.FullLoginResponse{
		LoginResponse: loginResp, // embute os tokens na raiz do JSON
		User: dmAuth.UserInfo{
			ID:       vendor.ID, // se Cod for string e UserInfo.ID for int, troque o tipo de ID pra string
			Name:     vendor.Name,
			Username: vendor.Username,
			Level:    vendor.Level,
		},
		Catalog: dmAuth.CatalogResponse{
			Products: produtos,
			Groups:   grupos,
		},
	}

	// 6) Responde no formato que o app espera + extras
	writeJSON(w, http.StatusOK, fullResp)
}

// HealthHandler – simples endpoint de health check: GET /health
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
