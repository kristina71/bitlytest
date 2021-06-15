package config

type DbConfig struct {
	Host     string
	User     string
	Password string
	DbName   string
	Sslmode  string
	Port     string
	DbPort   string
}

const (
	host     = "localhost"
	user     = "postgres"
	password = ""
	dbName   = "bitlytest"
	sslmode  = "disable"
	port     = "8000"
	dbPort   = "5432"
)

func New() DbConfig {
	return DbConfig{
		Host:     host,
		User:     user,
		Password: password,
		DbName:   dbName,
		Sslmode:  sslmode,
		Port:     port,
		DbPort:   dbPort,
	}
}
