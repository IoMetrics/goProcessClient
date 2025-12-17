package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	configpkg "goProcessClient/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

var localDB *sql.DB

func InitLocalDB(cfg configpkg.LocalDBConfig) error {
	// Aqui eu deixo senha vazia por enquanto (você pode colocar no INI depois)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=latin1&loc=Local",
		cfg.Usuario, cfg.Senha, cfg.IP, 3306, cfg.Banco,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return err
	}

	localDB = db
	return nil
}

func GetLocalDB() *sql.DB {
	return localDB
}

func CloseLocalDB() {
	if localDB != nil {
		_ = localDB.Close()
		localDB = nil
	}
}
