package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/yamaki-87/mailbot/src/config"
	"github.com/yamaki-87/mailbot/src/domain"
	"github.com/yamaki-87/mailbot/src/logger"
	"github.com/yamaki-87/mailbot/src/mail"
	mailtmpl "github.com/yamaki-87/mailbot/src/mail_tmpl"
)

func main() {
	logger.Init()
	token := os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal().Msgf("Bot作成失敗:%v", err)
	}

	mailService := mail.MailService{}
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return
		}

		mailSendType, err := mail.ParseMailSendType(m.Content)
		if err != nil {
			log.Err(err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ %v...", err))
			return
		}
		if mailSendType.Type == domain.Unknown {
			return
		}
		log.Info().Msgf("%+v", mailSendType)

		mailTmpl := mailService.GetTmplPath(mailSendType.Type)
		if mailTmpl == "" {
			log.Fatal().Msg("tmplが見つかりません")
		}
		mailBind := mailService.BindTemplate(mailSendType)
		mailS, err := mailtmpl.CreateMailTmpl(mailBind, mailTmpl)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⚠️ メールテンプレート処理失敗...")
			log.Error().Err(err)
			return
		}

		s.ChannelMessageSend(m.ChannelID, "📤 有給メールを送信します...")
		err = mail.SendMail(mailS)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⚠️ メール送信処理失敗...")
			log.Error().Err(err)
			return
		}

		log.Info().Msg("✅ メール送信完了！")
		s.ChannelMessageSend(m.ChannelID, "✅ メール送信完了！")
	})

	err = dg.Open()
	if err != nil {
		log.Fatal().Msgf("接続失敗: %v", err)
	}
	defer dg.Close()
	log.Info().Msg("Bot起動中。Ctrl+Cで終了")
	select {} // 無限待機
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf(".envファイルの読み込みに失敗しました：%v", err)
	}

	_ = config.GetConfig()

}
