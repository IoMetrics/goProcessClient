package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type LocalDBConfig struct {
	Banco     string
	IP        string
	Usuario   string
	Senha     string
	PastaErro string
}

type RemoteDBConfig struct {
	Banco   string
	IP      string
	Usuario string
	Senha   string
}

type AppConfig struct {
	LocalDB  LocalDBConfig
	RemoteDB RemoteDBConfig
}

// ResolveINIPath procura o INI no diretório do executável
func ResolveINIPath(filename string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exe), filename), nil
}

func LoadFromINI(path string) (AppConfig, error) {
	f, err := ini.Load(path)
	if err != nil {
		return AppConfig{}, fmt.Errorf("não foi possível ler %s: %w", path, err)
	}

	sec := f.Section("Config")
	rmsec := f.Section("Remote_Config")

	cfg := AppConfig{
		LocalDB: LocalDBConfig{
			Banco:     sec.Key("Banco").String(),
			IP:        sec.Key("ip").String(),
			Usuario:   sec.Key("usuario").String(),
			Senha:     sec.Key("senha").String(),
			PastaErro: sec.Key("pasta_erro").String(),
		},

		RemoteDB: RemoteDBConfig{
			Banco:   rmsec.Key("Banco").String(),
			IP:      rmsec.Key("ip").String(),
			Usuario: rmsec.Key("usuario").String(),
			Senha:   rmsec.Key("senha").String(),
		},
	}

	// defaults simples
	if cfg.LocalDB.Banco == "" {
		cfg.LocalDB.Banco = "autocomnb"
	}
	if cfg.LocalDB.IP == "" {
		cfg.LocalDB.IP = "localhost"
	}
	if cfg.LocalDB.Usuario == "" {
		cfg.LocalDB.Usuario = "root"
	}
	if cfg.LocalDB.Senha == "" {
		cfg.LocalDB.Senha = "2525"
	}
	if cfg.LocalDB.PastaErro == "" {
		cfg.LocalDB.PastaErro = filepath.Join(".", "erro")
	}

	_ = os.MkdirAll(cfg.LocalDB.PastaErro, 0o755)

	return cfg, nil
}
