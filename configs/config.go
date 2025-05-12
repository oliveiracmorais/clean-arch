package configs

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
	AmpqURL           string `mapstructure:"AMPQ_URL"`
	AmpqPort          string `mapstructure:"AMPQ_PORT"`
}

func NewConfig() (*conf, error) {
	_, err := os.Stat(".env")

	if os.IsExist(err) {
		fmt.Println("Tentou carregar o arquivo .env")
		return loadConfig(".")
	}
	fmt.Println("Carregando variáveis de ambiente")
	return loadConfigEnvironment(), nil
}

func loadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}

func loadConfigEnvironment() *conf {
	cfg := &conf{
		DBDriver:          os.Getenv("DB_DRIVER"),
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		WebServerPort:     os.Getenv("WEB_SERVER_PORT"),
		GRPCServerPort:    os.Getenv("GRPC_SERVER_PORT"),
		GraphQLServerPort: os.Getenv("GRAPHQL_SERVER_PORT"),
		AmpqURL:           "amqp://guest:guest@rabbitmq",
		AmpqPort:          os.Getenv("AMPQ_PORT"),
	}

	return cfg
}
