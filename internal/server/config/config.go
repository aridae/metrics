package config

import "sync"

const (
	defaultAddress = "localhost:8080"
)

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
	// значения конфига по умолчанию
	c.defaults()

	// инициализация структуры конфига
	// из значений, переданных через флаги
	configValuesFromFlags := parseFlags().configSetters()
	c.eval(configValuesFromFlags...)

	// env, если есть, затирает флаги
	configValuesFromEnv := readEnv().configSetters()
	c.eval(configValuesFromEnv...)
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

func (c *config) defaults() {
	c.address = defaultAddress
}
