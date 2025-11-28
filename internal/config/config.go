package config

// Endereço onde o servidor HTTP vai ouvir
const ServerAddr = ":8080"

var JwtSecret = []byte("10m3tr1cs2024!") // segredo para assinar tokens JWT

// DSN de conexão com o banco clientesweb
// Ajuste com seu usuário/senha/host/porta
// formato: usuario:senha@tcp(ip:porta)/database?params
const ClientesWebDSN = "clientesweb:clientesweb@tcp(127.0.0.1:3306)/clientesweb?charset=latin1&parseTime=true"
