package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Log         LogConfig         `mapstructure:"log"`
	MailTmpl    MailTmplConfig    `mapstructure:"mailTmpl"`
	TimeCard    TimeCard          `mapstructure:"timecard"`
	MessageTmpl MessageTmplConfig `mapstructure:"messageTmpl"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

type MailTmplConfig struct {
	PaidLeave    string `mapstructure:"PaidLeave"`
	LateArrival  string `mapstructure:"LateArrival"`
	Absence      string `mapstructure:"Absence"`
	TimeCard     string `mapstructure:"TimeCard"`
	SpecialLeave string `mapstructure:"SpecialLeave"`
}

type TimeCard struct {
	Path string `mapstructure:"path"`
}

type MessageTmplConfig struct {
	Help string `mapstructure:"Help"`
}

// シングルトンインスタンスとロック
var (
	instance *AppConfig
	once     sync.Once
)

// 呼び出し時に一度だけ読み込み
func GetConfig() *AppConfig {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
		viper.AddConfigPath("config")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("設定ファイルの読み込み失敗: %v", err)
		}

		var cfg AppConfig
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Fatalf("設定の構造体変換失敗: %v", err)
		}

		instance = &cfg
	})

	return instance
}
