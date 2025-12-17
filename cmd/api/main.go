package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	configpkg "goProcessClient/internal/config"
	dbpkg "goProcessClient/internal/data/db"

	authhdl "goProcessClient/internal/http/handlers/auth"
	billhdl "goProcessClient/internal/http/handlers/bill"
	orderhdl "goProcessClient/internal/http/handlers/order"
	// configdm "goProcessClient/internal/domain" // <-- ajuste o path se seu dmLocalDb estiver em outro pacote
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("erro na aplicação: %v", err)
	}
}

func run() error {
	// 1) Carrega INI ao lado do executável
	iniPath, err := configpkg.ResolveINIPath("goProcessClient.ini")
	if err != nil {
		return err
	}

	appCfg, err := configpkg.LoadFromINI(iniPath)
	if err != nil {
		return err
	}

	log.Printf("Config carregado: DB=%s IP=%s Usuario=%s PastaErro=%s\n",
		appCfg.LocalDB.Banco, appCfg.LocalDB.IP, appCfg.LocalDB.Usuario, appCfg.LocalDB.PastaErro,
	)

	// 2) Converte INI -> dmLocalDb e dmRemoteDB (seu domain model)
	dmLocalDb := configpkg.LocalDBConfig{
		Banco:     appCfg.LocalDB.Banco,
		IP:        appCfg.LocalDB.IP,
		Usuario:   appCfg.LocalDB.Usuario,
		Senha:     appCfg.LocalDB.Senha,
		PastaErro: appCfg.LocalDB.PastaErro,
		// Se existir: Senha/Porta etc.
	}

	dmRemoteDb := configpkg.RemoteDBConfig{
		Banco:   appCfg.RemoteDB.Banco,
		IP:      appCfg.RemoteDB.IP,
		Usuario: appCfg.RemoteDB.Usuario,
		Senha:   appCfg.RemoteDB.Senha,
	}
	// Se existir: Senha/Porta etc.

	// 3) Inicializa DB local (Autocom)
	defer dbpkg.CloseLocalDB()
	if err := dbpkg.InitLocalDB(dmLocalDb); err != nil {
		return fmt.Errorf("InitLocalDB: %w", err)
	}

	// 4) Inicializa DB clientesweb (como já era)
	defer dbpkg.CloseRemoteDB()
	if err := dbpkg.InitRemoteDB(dmRemoteDb); err != nil {
		return err
	}

	// 5) Configura handler de order (pasta erro)
	orderhdl.Configure(orderhdl.Options{
		ErrorDir: appCfg.LocalDB.PastaErro,
		UseINI:   true,
	})

	// 6) Rotas
	mux := http.NewServeMux()
	mux.HandleFunc("/login", authhdl.LoginHandler)
	mux.HandleFunc("/health", authhdl.HealthHandler)
	mux.Handle("/bill", billhdl.BillMiddleware(http.HandlerFunc(billhdl.BillHandler)))
	mux.HandleFunc("/order/send", orderhdl.SendOrderHandler) // <-- era order.SendOrderHandler (errado)

	server := &http.Server{
		Addr:    configpkg.ServerAddr,
		Handler: mux,
	}

	// Sobe o servidor
	go func() {
		log.Printf("Servidor escutando em %s\n", configpkg.ServerAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Erro ao iniciar servidor HTTP: %v", err)
		}
	}()

	// Shutdown gracioso
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Sinal recebido, desligando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Servidor finalizado com sucesso.")
	return nil
}
