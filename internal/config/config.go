package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)





type Config struct{
	Env 	string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	DBConfig	`yaml:"db"`
}

type HTTPServer struct{
	Adress		string `yaml:"adress" env-default:"localhost:8084"`
	Timeout 	string `yaml:"timeout" env-default:"4s"`
	IdleTimeout	string `yaml:"idle_timeout" env-default:"60s"`
}

type DBConfig struct{
	Host		string `yaml:"host"`
	Port		string `yaml:"port"`
	User		string `yaml:"user"`
	Password 	string `yaml:"password"`
	DBName		string `yaml:"dbname"`
	SSLMode		string `yaml:"sslmode"`
}


//exit program if config not found
func MustLoad() *Config{
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == ""{
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("config file not exists: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg);err!=nil{
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg

}