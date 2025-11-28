package main

import (
	"context"
	"errors"
	configpkg "goProcessClient/internal/config"
	dbpkg "goProcessClient/internal/data/db"
	authhdl "goProcessClient/internal/http/handlers/auth"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("erro na aplicação: %v", err)
	}
}

func run() error {
	// Inicializa conexão com clientesweb
	defer dbpkg.CloseClientesWebDB()

	if err := dbpkg.InitClientesWebDB(); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/login", authhdl.LoginHandler)
	mux.HandleFunc("/health", authhdl.HealthHandler)

	// Exemplo de rota protegida (quando você criar):
	// mux.Handle("/alguma-rota-protegida", AuthMiddleware(http.HandlerFunc(SuaHandlerProtegida)))

	server := &http.Server{
		Addr:    configpkg.ServerAddr,
		Handler: mux,
	}

	// Sobe o servidor em uma goroutine
	go func() {
		log.Printf("Servidor escutando em %s\n", configpkg.ServerAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Erro ao iniciar servidor HTTP: %v", err)
		}
	}()

	// Espera sinal de interrupção (Ctrl+C, kill, etc.)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Sinal recebido, desligando servidor...")

	// Contexto com timeout para shutdown gracioso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Servidor finalizado com sucesso.")
	return nil
}
