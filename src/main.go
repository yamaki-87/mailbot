package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/yamaki-87/mailbot/src/config"
	discohandler "github.com/yamaki-87/mailbot/src/disco_handler"
	"github.com/yamaki-87/mailbot/src/logger"
)

func main() {
	logger.Init()
	discohandler.DiscordBootstrap()
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf(".envファイルの読み込みに失敗しました：%v", err)
	}

	_ = config.GetConfig()

}
