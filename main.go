package main

import (
	"log"
	"sport-scraping/cmd"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
