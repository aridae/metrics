package config

import "sync"

type Config interface {
	GetAddress() string
}

var (
	once         sync.Once
	globalConfig *config
)

type config struct {
	address string
}

func (c *config) GetAddress() string {
	return c.address
}

func ObtainFromFlags() Config {
	once.Do(func() {
		globalConfig = &config{}
		globalConfig.initFromFlags()
	})

	return globalConfig
}

func (c *config) initFromFlags() {
	// инициализация структуры конфига
	// из значений, переданных через флаги
	configValuesFromFlags := parseFlags().configSetters()
	c.eval(configValuesFromFlags...)
}

func (c *config) isInit() bool {
	return globalConfig != nil
}

type configSetter func(cfg *config)

func (c *config) eval(setters ...configSetter) {
	for _, setter := range setters {
		setter(c)
	}
}
