package config

import (
    "flag"
    "log"
    "os"

    "github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
    Addr string `yaml:"address" env-required:"true"`
}

type Config struct {
    Env         string     `yaml:"env" env:"ENV" env-required:"true"`
    StoragePath string     `yaml:"storage_path" env-required:"true"`
    HttpServer  `yaml:"http_server"`
}

// Hold flag value globally so it's defined only ONCE
var configPath string

func init() {
    // Define -config flag exactly once
    flag.StringVar(&configPath, "config", "", "path to configuration file")
}

func MustLoad() *Config {
    // Read CONFIG_PATH environment variable
    if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
        configPath = envPath
    }

    // Read command-line flags only once
    flag.Parse()

    if configPath == "" {
        log.Fatalf("Config path is not set")
    }

    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        log.Fatalf("Config file does not exist: %s", err.Error())
    }

    var cfg Config
    if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
        log.Fatalf("cannot read config file: %v", err)
    }

    return &cfg
}