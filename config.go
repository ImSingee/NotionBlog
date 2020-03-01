package main

import (
	"github.com/spf13/viper"
	"log"
)

func setDefaultConfig() {
	viper.SetDefault("version", 1)
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
			log.Println("Cannot find config file.")
		} else {
			log.Fatal("Can not open or parse config.yml: ", err)
		}
	}

	_ = viper.BindEnv("token_v2", "NOTION_TOKEN")
	_ = viper.BindEnv("converter.force", "NOTION_CONVERTER_FORCE")
	_ = viper.BindEnv("render.checkbox", "NOTION_RENDER_CHECKBOX")
	_ = viper.BindEnv("database.post", "NOTION_DATABASE_POST")
	_ = viper.BindEnv("user.locale", "NOTION_USER_LOCALE")
	_ = viper.BindEnv("user.timezone", "NOTION_USER_TIMEZONE")
}
