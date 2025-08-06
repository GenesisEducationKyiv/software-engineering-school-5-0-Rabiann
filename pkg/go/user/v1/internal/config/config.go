package config

type Configuration struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string
	PostgresHost     string
	Prod             string
	ProdDbUrl        string
}
