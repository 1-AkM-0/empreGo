package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
}

func (b *Bot) SendMessage(channelID, content string) (*discordgo.Message, error) {
	return b.session.ChannelMessageSend(channelID, content)
}

func (b *Bot) Close() error {
	err := b.session.Close()
	if err != nil {
		return fmt.Errorf("erro ao tentar finalizar sessao do bot")
	}
	return nil
}

func NewBot(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("erro ao tentar criar a sessao do bot: %w", err)
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages

	err = session.Open()
	if err != nil {
		return nil, fmt.Errorf("erro ao tentar abrir conexao com o discord: %w", err)
	}

	bot := &Bot{session: session}

	return bot, nil
}
