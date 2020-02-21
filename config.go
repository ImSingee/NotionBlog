package main

import (
	"github.com/spf13/viper"
	"log"
)

func setDefaultConfig() {
	viper.SetDefault("version", 1)
	viper.SetDefault("converter.version", 0)
	viper.SetDefault("user.locale", "en")
	viper.SetDefault("user.timezone", "Etc/UTC")
	viper.SetDefault("render.checkbox", false)
}

func loadConfig() {
	setDefaultConfig()

	viper.SetConfigName("config")
	viper.AddConfigPath(notionDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore
			log.Fatal("Cannot find config file.")
		} else {
			log.Fatal("Can not open or parse config.yml: ", err)
		}
	}
}

func saveConfig() {
	err := viper.WriteConfig()
	if err != nil {
		log.Println("Cannot refresh config")
	}

}
