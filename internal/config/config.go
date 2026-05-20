package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	SigningKey = iota
	TokenTTL
	RedisHost
	RedisPassword
	KafkaBrokers
)

var conf *Config

type Config struct {
	SigningKey    string        `yaml:"signing-key"`
	TokenTTL      time.Duration `yaml:"token-ttl"`
	RedisHost     string        `yaml:"redis-host"`
	RedisPassword string        `yaml:"redis-password"`
	KafkaBrokers  string        `yaml:"kafka-brokers"`
}

func Init() error {
	ENV := os.Getenv("ENV")

	body, err := os.ReadFile(fmt.Sprintf("./configs/values_%s.yaml", strings.ToLower(ENV)))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, &conf)

	conf.RedisPassword = os.Getenv("REDIS_PASSWORD")
	conf.SigningKey = os.Getenv("SIGNING_KEY")

	return err
}

func Get(key int) interface{} {
	switch key {
	case SigningKey:
		return conf.SigningKey
	case TokenTTL:
		return conf.TokenTTL
	case RedisHost:
		return conf.RedisHost
	case RedisPassword:
		return conf.RedisPassword
	case KafkaBrokers:
		return strings.Split(conf.KafkaBrokers, ",")
	default:
		panic(ErrConfigNotFoundByKey(key))
	}
}
