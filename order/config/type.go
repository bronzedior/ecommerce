package config

type Config struct {
	App AppConfig `yaml:"app" validate:"required"`
}

type AppConfig struct {
	Port string `yaml:"port" validate:"required"`
}
