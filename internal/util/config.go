package util

type CORS struct {
	Host string
	Port int
}

type HTTP struct {
	Port int
	CORS CORS
}

type Metrics struct {
	Port int
}

type MongoDB struct {
	URI string
}

type Database struct {
	MongoDB MongoDB
}

type Config struct {
	Environment string
	LogLevel    string
	HTTP        HTTP
	Metrics     Metrics
	Database    Database
}
