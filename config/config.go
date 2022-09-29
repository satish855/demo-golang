package config

import "github.com/byteintellect/go_commons/config"

type CacheConfig struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password" envconfig:"CACHE_PASSWORD"`
}

type UserSvcConfig struct {
	config.BaseConfig `yaml:"base_config" json:"base_config"`
	CacheConfig       CacheConfig `yaml:"cache_config" json:"cache_config"`
}
