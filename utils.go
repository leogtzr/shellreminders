package main

import "github.com/spf13/viper"

func colorForMessages() map[string]string {
	return map[string]string{
		"red":    redHexColor,
		"yellow": yellowHexColor,
		"green":  greenHexColor,
	}
}

func readConfig(filename, configPath string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(configPath)
	v.SetConfigType("env")
	err := v.ReadInConfig()
	return v, err
}

func sendEmail(cfg *viper.Viper) error {
	return nil
}
