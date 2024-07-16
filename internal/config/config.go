package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/richardbizik/gommentary/internal/database"
	"github.com/richardbizik/gommentary/internal/profile"
	"github.com/richardbizik/gommentary/internal/rest/middleware"
)

// Config serves as a struct to which configuration is mapped into from ENV or configuration file
// it is a global variable that can be accessed from anywhere in a program
type Config struct {
	RestApi    RestApi              `yaml:"rest" json:"rest"`
	Database   database.Config      `yaml:"database" json:"database"`
	EnableOTEL bool                 `yaml:"enableOtel" json:"enableOtel" default:"false"`
	JWT        middleware.JWTConfig `yaml:"jwt" json:"jwt"`
}

type RestApi struct {
	Context string `yaml:"context" json:"context" env:"REST_API_CONTEXT" env-default:"/"`
	Port    int    `yaml:"port" json:"port" env:"REST_API_PORT" env-default:"8080"`
}

var (
	Conf Config
)

func InitConfig() Config {
	c := Config{}
	var fileName string
	confFile := os.Getenv("CONFIG_FILE")
	if confFile == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fileName = fmt.Sprintf("%s/conf-%s.yaml", wd, strings.ToLower(string(profile.Current)))
	} else {
		fileName = confFile
	}
	err := cleanenv.ReadConfig(fileName, &c)
	if err != nil {
		fmt.Printf("Error occurred while reading the config file: %s Error: %v", fileName, err)
		panic(err)
	}
	Conf = c
	return c
}
