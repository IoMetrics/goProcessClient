package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	configpkg "goProcessClient/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

var remoteDB *sql.DB

func InitRemoteDB(cfg configpkg.RemoteDBConfig) error {
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

	remoteDB = db
	return nil
}

func GetRemoteDB() *sql.DB {
	return remoteDB
}

func CloseRemoteDB() {
	if remoteDB != nil {
		_ = remoteDB.Close()
		remoteDB = nil
	}
}
