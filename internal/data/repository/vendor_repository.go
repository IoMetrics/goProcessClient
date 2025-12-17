package repository

import (
	"database/sql"
	"errors"
	"fmt"

	dbpkg "goProcessClient/internal/data/db"
	configdm "goProcessClient/internal/domain"
)

// BuscarUsuarios busca os vendedores + nível no banco AUTOCOM
func BuscarUsuarios() ([]configdm.UserInfo, error) {
	dbRemote := dbpkg.GetRemoteDB()
	if dbRemote == nil {
		return nil, fmt.Errorf("base remota não inicializada (GetRemoteDB=nil)")
	}

	const sqlVendedores = `
		SELECT a.cod, a.nome, a.usuario, b.nivel
		FROM vendedor a
		JOIN senhas b ON a.usuario = b.usuario
	`

	rows, err := dbRemote.Query(sqlVendedores)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar select de vendedores: %w", err)
	}
	defer rows.Close()

	var lista []configdm.UserInfo
	for rows.Next() {
		var v configdm.UserInfo
		if err := rows.Scan(&v.ID, &v.Name, &v.Username, &v.Level); err != nil {
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
func BuscarUsuarioPorLogin(usuario, senha string) (*configdm.UserInfo, error) {
	dbRemote := dbpkg.GetRemoteDB()
	if dbRemote == nil {
		return nil, fmt.Errorf("base remota não inicializada (GetRemoteDB=nil)")
	}

	const sqlUsuario = `
        SELECT a.cod, a.nome, a.usuario, b.nivel
        FROM vendedor a
        JOIN senhas b ON a.usuario = b.usuario
        WHERE b.usuario = ? AND b.senha = ?
    `

	var v configdm.UserInfo
	err := dbRemote.QueryRow(sqlUsuario, usuario, senha).Scan(
		&v.ID, &v.Name, &v.Username, &v.Level,
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
func BuscarUsuarioCod(cod string) (*configdm.UserInfo, error) {
	dbRemote := dbpkg.GetRemoteDB()
	if dbRemote == nil {
		return nil, fmt.Errorf("base remota não inicializada (GetRemoteDB=nil)")
	}

	const sqlUsuario = `
        SELECT a.cod, a.nome, a.usuario, b.nivel
        FROM vendedor a
        JOIN senhas b ON a.usuario = b.usuario
        WHERE a.cod = ?
    `
	var v configdm.UserInfo
	err := dbRemote.QueryRow(sqlUsuario, cod).Scan(
		&v.ID, &v.Name, &v.Username, &v.Level,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("usuario inválido")
		}
		return nil, fmt.Errorf("erro ao buscar usuario: %w", err)
	}

	return &v, nil
}

// -------------------------
// Catálogo (produtos/grupos)
// -------------------------

// BuscarProdutos carrega o catálogo de produtos no banco AUTOCOM
func BuscarProdutos() ([]configdm.ProductDTO, error) {
	dbRemote := dbpkg.GetRemoteDB()
	if dbRemote == nil {
		return nil, fmt.Errorf("base remota não inicializada (GetRemoteDB=nil)")
	}

	// Ajuste a query para o schema real do AUTOCOM
	const sqlProdutos = `
		SELECT
			cod,
			descricao,
			valor,
			grupo,
			detalhes,
			unidade
		FROM produtos where isvendamobilesemestoque="1" 
	`

	rows, err := dbRemote.Query(sqlProdutos)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar select de produtos: %w", err)
	}
	defer rows.Close()

	var lista []configdm.ProductDTO
	for rows.Next() {
		var p configdm.ProductDTO
		if err := rows.Scan(
			&p.Product,
			&p.Description,
			&p.Value,
			&p.GroupId,
			&p.Details,
			&p.Unit,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de produtos: %w", err)
		}
		lista = append(lista, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro no cursor de produtos: %w", err)
	}

	return lista, nil
}

// BuscarGruposProdutos carrega os grupos de produtos no banco AUTOCOM
func BuscarGruposProdutos() ([]configdm.ProductGroupDTO, error) {
	dbRemote := dbpkg.GetRemoteDB()
	if dbRemote == nil {
		return nil, fmt.Errorf("base remota não inicializada (GetRemoteDB=nil)")
	}

	// Ajuste a query para a tabela real de grupos no AUTOCOM
	const sqlGrupos = `
		SELECT
			codigo,
			descricao
		FROM grupo
	`

	rows, err := dbRemote.Query(sqlGrupos)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar select de grupos de produtos: %w", err)
	}
	defer rows.Close()

	var lista []configdm.ProductGroupDTO
	for rows.Next() {
		var g configdm.ProductGroupDTO
		if err := rows.Scan(
			&g.Id,
			&g.Description,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de grupos de produtos: %w", err)
		}
		lista = append(lista, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro no cursor de grupos de produtos: %w", err)
	}

	return lista, nil
}
