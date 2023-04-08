package config

import "github.com/spf13/viper"

type Config struct {
	Service ServiceConfig
	Mysql   mysqlConfig
	Email   EmailConfig
}
type ServiceConfig struct {
	Version               string
	Mode                  string
	Port                  string
	Secret                string
	PasswordTokenExpired  int
	AccessTokenExpiredIn  int
	RefreshTokenExpiredIn int
}
type EmailConfig struct {
	AppEmail    string
	AppPassword string
}
type mysqlConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func LoadConfig(fileName string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(fileName)
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var cf Config
	if err := v.Unmarshal(&cf); err != nil {
		return nil, err
	}

	return &cf, nil
}
