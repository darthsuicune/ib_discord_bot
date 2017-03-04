package main 

import (
	"bytes"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var guilds map[string]Guild

const (
	GUILD = 0
	RAID = 1
	DATE = 2
	TIME = 3
	TMZ = 4
)

func parseRaidCommand(s *discordgo.Session, m *discordgo.Message) string {
	chunks := strings.Split(m.Content, " ")
	if len(chunks) == 1 {
		return fullRaidInformation()
	}
	return parseChunks(chunks[1:])
}

func fullRaidInformation() string {
	var buffer bytes.Buffer
	for _, guild := range guilds {
		buffer.Write([]byte(guild.Rancor.Times()))
		buffer.Write([]byte(guild.Tank.Times()))
	}
	return buffer.String()
}

func parseDefaultRaidCommand(s * discordgo.Session, m *discordgo.Message) string {
	chunks := strings.Split(m.Content, " ")
	return parseDefaultRaid(chunks[1:])
}

func parseChunks(chunks []string) string {
	var buffer bytes.Buffer
	if len(chunks) == 1 {
		guild = strings.ToLower(chunks[GUILD])
		if guild != "eu" && guild != "us" {
			buffer.Write([]byte("You f\\*\\*\\*ing moron"))
		} else {
			buffer.Write([]byte(guilds[guild].Rancor.Times()))
			buffer.Write([]byte(guilds[guild].Tank.Times()))
		}
	}
	return buffer.String()
}

func parseDefaultRaid(chunks []string) string {
	if len(chunks) == 1 {
		
	}
}
