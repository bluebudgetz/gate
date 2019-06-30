package util

import "time"

type CORS struct {
	Host string
	Port int
}

type JWT struct {
	Key         string
	ExpDuration time.Duration
}

type HTTP struct {
	Port int
	CORS CORS
	JWT  JWT
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
