package domain

// DBNames agrupa os nomes dos bancos.
type DBNames struct {
	Autocom    string `json:"database_autocom"`
	NFeService string `json:"database_nfeservice"`
	Financeiro string `json:"database_financeiro"`
}

// PixConfig agrupa os campos de PIX.
type PixConfig struct {
	Origem         string `json:"pix_origem"`
	ClientID       string `json:"pix_clientid"`
	ClientSecret   string `json:"pix_clientsecret"`
	CertificadoCRT string `json:"pix_certificadocrt"`
	CertificadoKEY string `json:"pix_certificadokey"`
	URL            string `json:"pix_url"`
	URLHom         string `json:"pix_urlhom"`
	URLBase        string `json:"pix_urlbase"`
}

// FTPConfig agrupa os campos de FTP.
type FTPConfig struct {
	Host string `json:"ftphost"`
	User string `json:"ftpuser"`
	Pass string `json:"ftpwd"`
}

// WSConfig representa um registro da tabela webservice_databases.
// É usado internamente para saber em qual banco conectar.
type WSConfig struct {
	ID                 int
	Chave              string
	Empresa            string
	IP                 string
	Port               int
	User               string
	Password           string
	DatabaseAutocom    string
	DatabaseNFe        string
	DatabaseFinanceiro string
	WSModo             string
	Pix                PixConfig
	CertificadoCRT     string
	CertificadoKEY     string
	FTP                FTPConfig
}
