package discohandler

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/yamaki-87/mailbot/src/config"
	"github.com/yamaki-87/mailbot/src/domain"
	"github.com/yamaki-87/mailbot/src/mail"
	mailtmpl "github.com/yamaki-87/mailbot/src/mail_tmpl"
)

const FileName = "勤務表.pdf"

func FileHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	// メッセージの添付ファイルを確認
	if len(m.Attachments) == 0 {
		return
	}
	config := config.GetConfig()

	for _, attachment := range m.Attachments {
		log.Info().Msgf("受け取った添付ファイル: %s, URL: %s", attachment.Filename, attachment.URL)
		// ここで添付ファイルの処理を行うことができます
		if strings.HasPrefix(m.Content, "!勤務表") && attachment.ContentType == "application/pdf" {
			// 例えば、PDFファイルを保存するなどの処理
			s.ChannelMessageSend(m.ChannelID, "勤務表のPDFファイルを受け取りました。")
			resp, err := s.Request("GET", attachment.URL, nil)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "⚠️ 勤務表のPDFファイルの取得に失敗しました。")
				log.Err(err).Msg("勤務表のPDFファイルの取得に失敗しました")
				return
			}

			err = os.WriteFile(config.TimeCard.Path, resp, 0644)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "⚠️ 勤務表のPDFファイルの保存に失敗しました。")
				log.Err(err).Msg("勤務表のPDFファイルの保存に失敗しました")
				return
			}

			s.ChannelMessageSend(m.ChannelID, "勤務表のPDFファイルを保存しました。")
			log.Info().Msgf("勤務表のPDFファイルを保存しました: %s", config.TimeCard.Path)
		}
	}
}

func MailHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	mailService := mail.MailService{}
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
}

func DiscordBootstrap() {
	token := os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal().Msgf("Bot作成失敗:%v", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
	dg.AddHandler(MailHandler)
	dg.AddHandler(FileHandler)

	err = dg.Open()
	if err != nil {
		log.Fatal().Msgf("接続失敗: %v", err)
	}
	defer dg.Close()
	log.Info().Msg("Bot起動中。Ctrl+Cで終了")
	select {} // 無限待機
}
