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

func Obtain() Config {
	once.Do(func() {
		globalConfig = &config{}
		globalConfig.init()
	})

	return globalConfig
}

func (c *config) init() {
	// инициализация структуры конфига
	// из значений, переданных через флаги
	configValuesFromFlags := parseFlags().configSetters()

	//configValueFromEnv
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
