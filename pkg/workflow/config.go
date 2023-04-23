package workflow

import "go.uber.org/config"

type ManagerType string

const (
	MemoryManagerType   ManagerType = "memory"
	TemporalManagerType ManagerType = "temporal"
)

type Config struct {
	ManagerType ManagerType `yaml:"manager"`
}

func NewDefaultConfig() Config {
	return Config{
		ManagerType: MemoryManagerType,
	}
}

func NewConfig(provider config.Provider) (Config, error) {
	var c Config
	err := provider.Get("workflow").Populate(&c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
