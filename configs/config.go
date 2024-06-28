package configs

import "github.com/spf13/viper"

type PersistenceConfigs struct {
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
}

type Conf struct {
	ServerPort    string
	Persistence   PersistenceConfigs
	ReqsPerSecond int
	TokenConfigs  map[string]int
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.SetConfigFile("config.yaml")
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
