package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	machinery "github.com/RichardKnop/machinery/v1/config"
)

var (
	config *Config
)

type Config struct {
	API       APIConfig
	Redis     RedisConfig
	GoogleAPI GoogleAPIConfig
	Machinery machinery.Config
}

type APIConfig struct {
	Port int
}

type RedisConfig struct {
	Address       string
	MaxConnection int
}

type GoogleAPIConfig struct {
	APIKey string
}

func Load(path string) error {
	handleErr := func(err error) error {
		fmt.Fprintf(os.Stderr, "[config] Cannot open config file in %s, reason=%+v\n", path, err)
		return err
	}
	file, err := os.Open(path)
	if err != nil {
		return handleErr(err)
	}
	config = &Config{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return handleErr(err)
	}
	fmt.Printf("[config] Loaded config file %s\n", path)
	loadFromEnv(config)
	return nil
}

func mustInt(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return i
}

func ifExist(env string, overrideFn func(value string)) {
	if v := os.Getenv(env); len(v) > 0 {
		fmt.Printf("[config] Found %s=%s, override existing config\n", env, v)
		overrideFn(v)
	}
}

func loadFromEnv(config *Config) {
	ifExist("API_PORT", func(v string) {
		config.API.Port = mustInt(v)
	})
	ifExist("REDIS_ADDRESS", func(v string) {
		config.Redis.Address = v
	})
	ifExist("MACHINERY_BROKER", func(v string) {
		config.Machinery.Broker = v
	})
	ifExist("MACHINERY_RESULT_BACKEND", func(v string) {
		config.Machinery.ResultBackend = v
	})
	ifExist("GOOGLE_API_KEY", func(v string) {
		config.GoogleAPI.APIKey = v
	})
}

func GetConfig() Config {
	return *config
}
