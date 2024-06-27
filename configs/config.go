package configs

import "github.com/spf13/viper"

type conf struct {
	ReqsPerSecond int
	TokenConfigs  map[string]int
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
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