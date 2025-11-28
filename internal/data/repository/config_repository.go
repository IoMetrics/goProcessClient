// internal/repository/config_repository.go
package repository

import (
	"database/sql"
	"fmt"
	dbpkg "goProcessClient/internal/data/db"
	configdm "goProcessClient/internal/domain"
	"log"
)

// BuscarConfigPorChave busca o registro em webservice_databases pela ws_chave
func BuscarConfigPorChave(wsChave string) (*configdm.WSConfig, error) {
	const sqlSelect = `
		SELECT
			id,
			ws_chave,
			ip,
			port,
			user,
			password,
			database_autocom,
			database_nfeservice,
			database_financeiro,
			ws_modo
		FROM webservice_databases
		WHERE ws_chave = ?
		LIMIT 1
	`

	log.Printf("BuscarConfigPorChave: ws_chave=%q", wsChave)

	if dbpkg.ClientesWebDB == nil {
		return nil, fmt.Errorf("conexão com clientesweb não inicializada")
	}

	row := dbpkg.ClientesWebDB.QueryRow(sqlSelect, wsChave)

	var cfg configdm.WSConfig
	err := row.Scan(
		&cfg.ID,                 // id
		&cfg.Chave,              // ws_chave
		&cfg.IP,                 // ip
		&cfg.Port,               // port
		&cfg.User,               // user
		&cfg.Password,           // password
		&cfg.DatabaseAutocom,    // database_autocom
		&cfg.DatabaseNFe,        // database_nfeservice
		&cfg.DatabaseFinanceiro, // database_financeiro
		&cfg.WSModo,             // ws_modo
	)

	if err == sql.ErrNoRows {
		return nil, nil // chave não encontrada
	}
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
