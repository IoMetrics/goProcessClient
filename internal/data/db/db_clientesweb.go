package db

import (
	"database/sql"
	"log"
	"time"

	configpkg "goProcessClient/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

var ClientesWebDB *sql.DB

// InitClientesWebDB abre e configura a conexão com o banco clientesweb
func InitClientesWebDB() error {
	db, err := sql.Open("mysql", configpkg.ClientesWebDSN)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}

	ClientesWebDB = db
	log.Println("Conectado ao banco clientesweb com sucesso")
	return nil
}

// CloseClientesWebDB fecha a conexão com clientesweb (chamado no main)
func CloseClientesWebDB() {
	if ClientesWebDB != nil {
		if err := ClientesWebDB.Close(); err != nil {
			log.Printf("Erro ao fechar conexão clientesweb: %v", err)
		}
	}
}
