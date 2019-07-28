package main

import (
	"github.com/bluebudgetz/common/pkg/service"
	"github.com/bluebudgetz/gate/internal"
	"github.com/spf13/viper"
)

type HostAndPortConfig struct {
	Host string
	Port int
}

type GatePlugin struct {
	Accounts     HostAndPortConfig
	Transactions HostAndPortConfig
}

func (p *GatePlugin) Key() string {
	return "gate"
}

func (p *GatePlugin) Configure(viper *viper.Viper) error {
	type container struct {
		Accounts     *HostAndPortConfig
		Transactions *HostAndPortConfig
	}
	viper.SetDefault("accounts.host", "")
	viper.SetDefault("accounts.port", 0)
	viper.SetDefault("transactions.host", "")
	viper.SetDefault("transactions.port", 0)
	return viper.Unmarshal(&container{
		Accounts:     &p.Accounts,
		Transactions: &p.Transactions,
	})
}

func (p *GatePlugin) Init(svc *service.MicroService) error {
	handler := internal.NewHandler(svc.Environment())
	router := svc.Plugin(service.HttpKey).(*service.Http).Router()
	router.Handle("/v1/accounts", handler.CreateReverseProxy(p.Accounts.Host, p.Accounts.Port))
	router.Handle("/v1/transactions", handler.CreateReverseProxy(p.Transactions.Host, p.Transactions.Port))
	return nil
}

func main() {
	service.New(
		"gate",
		service.NewMongo(service.MongoConfig{}),
		service.NewMetrics(service.MetricsConfig{}),
		service.NewHttp(service.HttpConfig{}),
		&GatePlugin{},
	).Run()
}
