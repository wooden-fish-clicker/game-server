package configs

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type App struct {
	PrefixUrl string

	AppName string

	LogSavePath string
	LogSaveName string
	LogFileExt  string

	MaxLogFiles int

	ImageStaticPath string
	ImageSavePath   string
}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type MySql struct {
	ConnString string
	Name       string

	MySqlBase MySqlBase
}

type MySqlBase struct {
	ConnMaxLifeTime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

type MongoDB struct {
	ConnString string
	Name       string
}

type Jwt struct {
	Secret         string
	ExpirationDays int
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

type Service struct {
	Addr string
}

type Config struct {
	App     App
	Server  Server
	MySql   MySql
	MongoDB MongoDB
	Jwt     Jwt
	Redis   Redis
	Service Service
}

var C Config

// Setup initialize the configuration instance
func Setup() {
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file not found")
		} else {
			log.Fatalf("Config file was found but another error was produced")
		}
	}

	viper.AutomaticEnv()

	err := viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	godotenv.Load()

}
