package config

import "C"
import (
	"github.com/esonhugh/go-cli-template-v2/utils/Error"
	"github.com/spf13/viper"
)

// Config struct is a wrapper of viper

// GlobalConfig default Global Variable for Config
var GlobalConfig *viper.Viper

// SpecificInit func is Init with specific Config file
func Init(file string) {
	GlobalConfig = viper.New()
	GlobalConfig.SetConfigFile(file)
	GlobalConfig.SetConfigType("yaml")
	err := GlobalConfig.ReadInConfig()
	Error.HandleFatal(err, "Config File "+file+" can't read or not exist.")
	// Error.HandleError(GlobalConfig.Unmarshal(&C))
}

func SaveConfig() {
	// GlobalConfig.Set("log_level", C.LogLevel)
	err := GlobalConfig.WriteConfig()
	Error.HandleFatal(err)
}
