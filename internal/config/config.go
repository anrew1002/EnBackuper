package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

// type Config struct {
// 	TFTP     string
// 	net_file string
// 	test     bool
// 	tftpData string
// }

type Config struct {
	TFTP          string `env:"TFTP" env-required:"true"`
	TFTPData      string `env:"TFTP_SAVEPATH" env-required:"true"`
	NETWORKS_FILE string `env:"HOST" env-default:"networks.txt"`
	Test          bool
}

func MustLoad() *Config {

	var cfg Config

	// Читаем конфиг-файл и заполняем нашу структуру
	err := cleanenv.ReadConfig("./.env", &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
