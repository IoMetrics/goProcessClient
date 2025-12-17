package repository

import (
	"database/sql"
	"fmt"
	dbpkg "goProcessClient/internal/data/db"
	"goProcessClient/internal/domain"

	_ "github.com/go-sql-driver/mysql" // certifique-se que o driver está importado em algum lugar do projeto
)

// BuscarConta busca a conta na tabela vendastemp usando comanda.
func BuscarConta(comanda string) (*domain.BillData, error) {
	db := dbpkg.GetLocalDB()
	if db == nil {
		return nil, fmt.Errorf("base local não inicializada (GetLocalDB=nil)")
	}

	const sqlConta = `
		SELECT
			id,
			produto,
			qte,
			valor,
			desconto_item,
			comissao,
			comissaosrv,
			data,
			hora,
			horafim,
			fechada,
			pessoas,
			vendedor,
			obs,
			obs_item
		FROM vendastemp
		WHERE comanda = ? 
		ORDER BY id
	`

	rows, err := db.Query(sqlConta, comanda)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar vendastemp: %w", err)
	}
	defer rows.Close()

	var (
		bill      domain.BillData
		items     []domain.BillItem
		headerSet bool
	)

	for rows.Next() {
		var (
			id           int
			produto      sql.NullString
			qte          sql.NullFloat64
			valor        sql.NullFloat64
			descontoItem sql.NullFloat64
			comissao     sql.NullFloat64
			comissaoSrv  sql.NullFloat64
			data         sql.NullString
			hora         sql.NullString
			horafim      sql.NullString
			fechada      sql.NullString
			pessoas      sql.NullInt64
			vendedor     sql.NullString
			obs          sql.NullString
			obsItem      sql.NullString
		)

		if err := rows.Scan(
			&id,
			&produto,
			&qte,
			&valor,
			&descontoItem,
			&comissao,
			&comissaoSrv,
			&data,
			&hora,
			&horafim,
			&fechada,
			&pessoas,
			&vendedor,
			&obs,
			&obsItem,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de vendastemp: %w", err)
		}

		// Preenche cabeçalho da conta a partir da primeira linha
		if !headerSet {
			bill.Comanda = comanda
			bill.Date = safeString(data)
			bill.Time = safeString(hora)

			if horafim.Valid {
				s := horafim.String
				bill.HourEnd = &s
			}
			if fechada.Valid {
				s := fechada.String
				bill.ClosedAt = &s
			}
			if pessoas.Valid {
				bill.People = int(pessoas.Int64)
			}
			bill.Waiter = safeString(vendedor)
			bill.Obs = safeString(obs)

			// por enquanto fixo; depois pode vir de tabela de parâmetros
			bill.ServiceTaxPercent = 10.0

			headerSet = true
		}

		item := domain.BillItem{
			ID:            id,
			ProductCode:   safeString(produto),       // produto
			Quantity:      safeFloat64(qte),          // qte
			UnitPrice:     safeFloat64(valor),        // valor
			Discount:      safeFloat64(descontoItem), // desconto_item
			Commission:    safeFloat64(comissao),     // comissao
			CommissionSrv: safeFloat64(comissaoSrv),  // comissaosrv
			ObsItem:       safeString(obsItem),       // obs_item (detalhes do item)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro no cursor da consulta vendastemp: %w", err)
	}

	if !headerSet {
		// nenhuma linha para essa comanda
		return nil, fmt.Errorf("comanda %s não encontrada", comanda)
	}

	bill.Items = items
	return &bill, nil
}

// Helpers para lidar com sql.Null*
func safeString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func safeFloat64(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0
}
