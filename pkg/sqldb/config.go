package sqldb

import (
	"github.com/kelseyhightower/envconfig"
	"fmt"
)

// Config struct for app config
type config struct {
	Host     string `default:"Localhost" envconfig:"DB_HOST"`
	User     string `default:"admin" envconfig:"DB_USER"`
	Password string `default:"admin" envconfig:"DB_PASSWORD"`
	DBName   string `default:"bookmark" envconfig:"DB_NAME"`
	Port     string `default:"5433" envconfig:"DB_PORT"`
}

func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", 
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)
}
