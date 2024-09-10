package config

import "sync"

type Config interface {
	GetAddress() string
}

var (
	mu           sync.RWMutex
	globalConfig *config
)

type config struct {
	address string
}

func (c *config) GetAddress() string {
	return c.address
}

func ObtainFromFlags() Config {
	if globalConfig.isInit() {
		return globalConfig
	}

	conf := &config{}
	conf.initFromFlags()

	return conf
}

func (c *config) initFromFlags() {
	mu.Lock()
	defer mu.Unlock()

	// инициализация структуры конфига
	// из значений, переданных через флаги
	configValuesFromFlags := parseFlags().configSetters()
	c.eval(configValuesFromFlags...)
}

func (c *config) isInit() bool {
	mu.RLock()
	defer mu.RUnlock()

	return globalConfig != nil
}

type configSetter func(cfg *config)

func (c *config) eval(setters ...configSetter) {
	for _, setter := range setters {
		setter(c)
	}
}
