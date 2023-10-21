package d

import (
	"brand/b"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/apricot/neko"
)

var (
	Token = os.Getenv("TOKEN_DISCORD_APPLICATION")

	ChannelMarguerite = "1127907325202665482"
)

func Run() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		logs.Panic(err)
	}
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	dg.AddHandler(PingPong)
	dg.AddHandler(Inquiry)

	go func() {
		for {
			time.Sleep(time.Hour)
			if reply, update := b.Poll(); update {
				if _, err := dg.ChannelMessageSend(ChannelMarguerite, fmt.Sprint("Routinely check: ", reply)); err != nil {
					logs.Warning(err)
				}
			}
		}
	}()

	if err := dg.Open(); err != nil {
		logs.Panic(err)
	}
	defer dg.Close()
	logs.Info("Discord Application running.")
	neko.Block()
}

func PingPong(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	logs.Debug(m.ChannelID)

	switch m.Content {
	case "Ping":
		logs.Info("Pong!")
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "Pong":
		logs.Info("Ping?")
		s.ChannelMessageSend(m.ChannelID, "Ping?")
	}
}

func Inquiry(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != "Inquiry" {
		return
	}
	logs.Info("Inquiry.")
	reply, _ := b.Poll()
	if _, err := s.ChannelMessageSend(m.ChannelID, reply); err != nil {
		logs.Warning(err)
	}
}
