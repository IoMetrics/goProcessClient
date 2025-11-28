package db

import (
	"database/sql"
	"fmt"
	"time"

	confidm "goProcessClient/internal/domain"

	_ "github.com/go-sql-driver/mysql"
)

// OpenAutocomDB abre uma conexão com o banco AUTOCOM usando a configuração
func OpenAutocomDB(cfg *confidm.WSConfig) (*sql.DB, error) {
	dsnAutocom := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=latin1&parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.IP,
		cfg.Port,
		cfg.DatabaseAutocom,
	)

	dbAutocom, err := sql.Open("mysql", dsnAutocom)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão AUTOCOM: %w", err)
	}

	dbAutocom.SetMaxOpenConns(5)
	dbAutocom.SetMaxIdleConns(5)
	dbAutocom.SetConnMaxLifetime(5 * time.Minute)

	// Opcional: testar a conexão
	if err := dbAutocom.Ping(); err != nil {
		dbAutocom.Close()
		return nil, fmt.Errorf("erro ao conectar no AUTOCOM: %w", err)
	}

	return dbAutocom, nil
}
