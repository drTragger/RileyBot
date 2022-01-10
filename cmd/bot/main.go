package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/drTragger/rileyBot/internal/app/riley"
	"github.com/yanzay/tbot/v2"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "path", "configs/riley.toml", "Path to config file in .toml format")
}

func main() {
	config := riley.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Println("Could not find configs file. Using default values:", err)
	}
	tgBot := tbot.New(config.BotToken)
	server := riley.New(config, tgBot.Client(), tgBot)
	log.Fatal(server.Start())
}
