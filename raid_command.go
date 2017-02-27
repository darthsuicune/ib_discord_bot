package main

import (
	"bytes"
	"errors"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

var raids map[string]map[string]Raid

func parseRaidCommand(s *discordgo.Session, m *discordgo.Message) string {
	chunks := strings.Split(m.Content, " ")

	var buffer bytes.Buffer
	if len(chunks) == 1 {
		for _, guildRaids := range raids {
			for _, items := range guildRaids {
				buffer.Write([]byte(items.String()))
			}
		}
	} else if chunks[1] == "eu" || chunks[1] == "us" {
		for _, guildRaids := range raids[chunks[1]] {
			buffer.Write([]byte(guildRaids.String()))
		}
	} else if chunks[1] == "set" {
		channel, _ := s.Channel(m.ChannelID)
		member, _ := s.GuildMember(channel.GuildID, m.Author.ID)
		if canSetRaids(s, channel.GuildID, member.Roles) {
			if len(chunks) >= 5 {
				err := setRaid(chunks[2:])
				if err != nil {
					buffer.Write([]byte("Error when creating the raid " + err.Error()))
				} else {
					buffer.Write([]byte("I've altered the deal. New details:\n"))
					buffer.Write([]byte(raids[chunks[2]][chunks[3]].String()))
				}
			} else {
				buffer.Write([]byte("You are missing like half the stuff, dude. Do it like this: !raid set us Rancor 1970-01-01 00:00 TZ"))
			}
		} else {
			buffer.Write([]byte("You can't do that"))
		}
	} else if chunks[1] == "delete" {
		channel, _ := s.Channel(m.ChannelID)
		member, _ := s.GuildMember(channel.GuildID, m.Author.ID)
		if canSetRaids(s, channel.GuildID, member.Roles) {
			if len(chunks) < 4 || len(chunks) > 5 {
				buffer.Write([]byte("Wrong parameters, use !raid delete <us/eu> <raid_name>"))
			} else {
				delete(raids[chunks[2]], chunks[3])
				buffer.Write([]byte("I've altered the deal again. Pray that I don't alter it any further."))
			}
		}
	}
	return buffer.String()
}

func setRaid(chunks []string) error {
	if chunks[0] != "eu" && chunks[0] != "us" {
		return errors.New("Guild must be eu or us")
	}
	guild := strings.TrimSpace(chunks[0])
	if _, ok := raids[guild]; !ok {
		raids[guild] = make(map[string]Raid)
	}
	raidType := strings.TrimSpace(strings.ToLower(chunks[1]))
	when := strings.TrimSpace(strings.Join(chunks[2:], " "))

	var location *time.Location
	if len(chunks) > 4 {
		location, _ = time.LoadLocation(chunks[4])
	} else {
		location, _ = time.LoadLocation("GMT")
	}
	raidTime, err := time.ParseInLocation("2006-01-02 15:04 MST", when, location)
	//Only set the new raid if there wasn't any error while parsing the date.
	if err == nil {
		raids[guild][raidType] = Raid1{Guild: guild, Type: raidType, Timing: raidTime}
	}
	return err
}

//AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
type Raid1 struct {
	Guild  string
	Type   string
	Timing time.Time
}

func (r Raid1) Times() []string {
	s := make([]string, 3)
	return s
}
func (r Raid1) String() string { return "" }

func canSetRaids(s *discordgo.Session, guildID string, roles []string) bool {
	for _, role := range roles {
		r, _ := s.State.Role(guildID, role)
		if r.Name == "Inglorious Leaders" || r.Name == "Inglorious Officers" {
			return true
		}
	}
	return false
}
