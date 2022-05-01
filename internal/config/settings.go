package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

//Settings represents the configuration that we can provide
//from the outside in order to run the application in different ways.
type Settings struct {
	Host            string `mapstructure:"HOST"`
	Port            string `mapstructure:"PORT"`
	Env             string `mapstructure:"ENV"`
	DriverName      string `mapstructure:"SQL_DRIVER_NAME"`
	ConnString      string `mapstructure:"SQL_CONN_STRING"`
	ConnMaxLifetime int64  `mapstructure:"SQL_CONN_MAX_LIFETIME"`
	ConnMaxOpen     int    `mapstructure:"SQL_CONN_MAX_OPEN_CONN"`
	ConnMaxIdle     int    `mapstructure:"SQL_CONN_MAX_IDLE_CONN"`
	HttpReqTimeout  int    `mapstructure:"HTTP_REQUEST_TIMEOUT"`
}

func New(envPath string) (*Settings, error) {
	var cfg Settings
	viper.SetConfigFile(envPath)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "No env file found")
	}
	//try to assign read variables into golang struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "Error trying to unmarshal configuration")
	}
	return &cfg, nil
}
