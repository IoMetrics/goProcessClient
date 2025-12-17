package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	dbpkg "goProcessClient/internal/data/db"
	dmOrder "goProcessClient/internal/domain/order"
)

func SaveOrderToVendastemp(ctx context.Context, req dmOrder.SendOrderRequest) error {
	db := dbpkg.GetLocalDB()
	if db == nil {
		return fmt.Errorf("base local não inicializada (GetLocalDB=nil)")
	}

	cctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	stmt, err := db.PrepareContext(cctx, `
		INSERT INTO vendastemp
			(caixa, venda, item, produto, qte, valor, vendedor, comanda, data, hora, unidade, grupo, descricao, obs_item)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	caixa := atoiDefault(req.Pdv, 0)
	comanda := atoiDefault(req.Comanda, 0)

	// aqui eu usei venda = comanda (bem comum em PDV: caixa/venda/item)
	venda := comanda

	now := time.Now()
	data := now.Format("2006-01-02")
	hora := now.Format("15:04:05")

	for i, it := range req.Items {
		prod := it.ProductDto.Product
		if prod == "" {
			prod = it.Product
		}

		unidade := ""
		if it.ProductDto.Unit != nil {
			unidade = *it.ProductDto.Unit
		}

		_, err := stmt.ExecContext(
			cctx,
			caixa,
			venda,
			i+1,
			prod,
			it.Quantity,
			it.UnitPrice,
			req.Vendedor,
			comanda,
			data,
			hora,
			unidade,
			it.ProductDto.GroupId,
			it.ProductDto.Description,
			it.ObsItem,
		)
		if err != nil {
			log.Printf("Erro ao inserir vendastemp item %d: %v", i+1, err)
			return fmt.Errorf("vendastemp item %d: %w", i+1, err)
		}
	}

	return nil
}

func atoiDefault(s string, def int) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	n := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return def
		}
		n = n*10 + int(ch-'0')
	}
	return n
}

// só pra evitar import "database/sql" não usado em alguns projetos
var _ = sql.ErrNoRows

func SaveIncomingOrder(baseDir string, req dmOrder.SendOrderRequest) (string, error) {
	// Ex: inbox/autocom/2025-12-13/<id>.json
	dateDir := time.Now().Format("2006-01-02")
	dir := filepath.Join(baseDir, req.LocalDB, dateDir)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir inbox: %w", err)
	}

	id := newID()
	filename := filepath.Join(dir, id+".json")

	b, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	if err := os.WriteFile(filename, b, 0o644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return filename, nil
}

func newID() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

type OrderErrorEnvelope struct {
	Order     dmOrder.SendOrderRequest `json:"order"`
	When      time.Time                `json:"when"`
	LastError string                   `json:"last_error"`
}

func SaveOrderErrorFile(errorDir string, req dmOrder.SendOrderRequest, err error) (string, error) {
	if errorDir == "" {
		errorDir = "./erro"
	}

	_ = os.MkdirAll(errorDir, 0o755)

	env := OrderErrorEnvelope{
		Order:     req,
		When:      time.Now(),
		LastError: err.Error(),
	}

	b, jerr := json.MarshalIndent(env, "", "  ")
	if jerr != nil {
		return "", jerr
	}

	name := fmt.Sprintf("order_%s_%d.json", req.Comanda, time.Now().UnixNano())
	full := filepath.Join(errorDir, name)

	if werr := os.WriteFile(full, b, 0o644); werr != nil {
		return "", werr
	}

	return full, nil
}
