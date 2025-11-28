package repository

import (
	"database/sql"
	"errors"
	"fmt"

	dbpkg "goProcessClient/internal/data/db"
	configdm "goProcessClient/internal/domain"
)

// BuscarUsuarios busca os vendedores + nível no banco AUTOCOM
func BuscarUsuarios(cfg *configdm.WSConfig) ([]configdm.Vendor, error) {
	dbAutocom, err := dbpkg.OpenAutocomDB(cfg)
	if err != nil {
		return nil, err
	}
	defer dbAutocom.Close()

	const sqlVendedores = `
		SELECT a.cod, a.nome, a.usuario, b.nivel
		FROM vendedor a
		JOIN senhas b ON a.usuario = b.usuario
	`

	rows, err := dbAutocom.Query(sqlVendedores)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar select de vendedores: %w", err)
	}
	defer rows.Close()

	var lista []configdm.Vendor
	for rows.Next() {
		var v configdm.Vendor
		if err := rows.Scan(&v.Cod, &v.Nome, &v.Usuario, &v.Nivel); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de vendedores: %w", err)
		}
		lista = append(lista, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro no cursor de vendedores: %w", err)
	}

	return lista, nil
}

// BuscarUsuarioPorLogin valida usuario + senha no Autocom
func BuscarUsuarioPorLogin(cfg *configdm.WSConfig, usuario, senha string) (*configdm.Vendor, error) {
	dbAutocom, err := dbpkg.OpenAutocomDB(cfg)
	if err != nil {
		return nil, err
	}
	defer dbAutocom.Close()

	const sqlUsuario = `
        SELECT a.cod, a.nome, a.usuario, b.nivel
        FROM vendedor a
        JOIN senhas b ON a.usuario = b.usuario
        WHERE b.usuario = ? AND b.senha = ?
    `

	var v configdm.Vendor
	err = dbAutocom.QueryRow(sqlUsuario, usuario, senha).Scan(
		&v.Cod, &v.Nome, &v.Usuario, &v.Nivel,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("usuario ou senha inválidos")
		}
		return nil, fmt.Errorf("erro ao buscar usuario: %w", err)
	}

	return &v, nil
}

// BuscarUsuarioCod busca o usuário apenas pelo código
func BuscarUsuarioCod(cfg *configdm.WSConfig, cod string) (*configdm.Vendor, error) {
	dbAutocom, err := dbpkg.OpenAutocomDB(cfg)
	if err != nil {
		return nil, err
	}
	defer dbAutocom.Close()

	const sqlUsuario = `
        SELECT a.cod, a.nome, a.usuario, b.nivel
        FROM vendedor a
        JOIN senhas b ON a.usuario = b.usuario
        WHERE a.cod = ?
    `

	var v configdm.Vendor
	err = dbAutocom.QueryRow(sqlUsuario, cod).Scan(
		&v.Cod, &v.Nome, &v.Usuario, &v.Nivel,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("usuario inválido")
		}
		return nil, fmt.Errorf("erro ao buscar usuario: %w", err)
	}

	return &v, nil
}
