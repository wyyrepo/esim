package container

import (
	"github.com/google/wire"
	"github.com/jukylin/esim/config"
	"github.com/jukylin/esim/prometheus"
)

func provideMockConf() config.Config {
	conf := config.NewMemConfig()
	conf.Set("debug", true)
	return conf
}

func provideMockProme(conf config.Config) *prometheus.Prometheus {
	return prometheus.NewNullProme()
}

var MockSet = wire.NewSet(
	wire.Struct(new(Esim), "*"),
	provideMockConf,
	provideLogger,
	provideMockProme,
)
