package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

const (
	GUILD = 0
	RAID  = 1
	DATE  = 2
	TIME  = 3
	TMZ   = 4

	DATE_FORMAT = "2006-01-02"

	EU = "eu"
	US = "us"

	RANCOR = "rancor"
	TANK   = "tank"
)

/**
 * Raid command should be any of:
 * !raid
 * !raid eu/us
 * !raid set eu/us <type> <date> <time> <TZ>
 * !raid delete eu/us <type>
 */
func readRaidCommand(s string, canSetRaids bool) string {
	chunks := strings.Split(s, " ")
	if len(chunks) == 1 {
		return guilds.Raids()
	}
	if len(chunks) == 2 {
		var buffer bytes.Buffer
		guild := strings.ToLower(chunks[1])
		if guild != EU && guild != US {
			buffer.WriteString("You f\\*\\*\\*ing moron")
		} else {
			buffer.WriteString(fmt.Sprintf("**%s**\n%s", strings.ToUpper(guild), guilds.Guild(guild).Raids()))
		}
		return buffer.String()
	}
	if len(chunks) == 7 {
		if !canSetRaids {
			return "Nope"
		}
		if chunks[1] != "set" {
			return "Your training is not complete yet. " +
				"You must head to planet Dagobah to learn about \"set\""
		}
		return setCustomRaid(chunks[2:])
	}
	if len(chunks) == 4 {
		if !canSetRaids {
			return "Nope"
		}
		if chunks[1] != "delete" {
			return "Your training is not complete yet. " +
				"You must head to planet Dagobah to learn about \"delete\""
		}
		return deleteRaid(chunks[2:])
	}
	return "Wait, wat?"
}

func setCustomRaid(s []string) string {
	g := s[GUILD]
	raidType := s[RAID]
	date := s[DATE]
	timing := s[TIME]
	tmz := s[TMZ]

	if g != EU && g != US {
		return "You're an idiot. Only the EU and US exist, the rest of the world doesn't."
	}
	if raidType != RANCOR && raidType != TANK {
		return "Only rancors and tanks m8. The FBI might raid your home to search for those drugs tho."
	}

	fullTiming := fmt.Sprintf("%s %s %s", date, timing, tmz)

	guild := guilds.Guild(g)
	startTime, err := time.ParseInLocation("2006-01-02 15:04 MST", fullTiming, guild.Location)
	if err != nil {return err.Error()}
	switch raidType {
	case RANCOR:
		guild.SetRancor(startTime)
	case TANK:
		if g == EU {
			guild.SetEUTank(startTime)
		} else {
			guild.SetUSTank(startTime)
		}
	}
	return "Get your lightswords ready. The Rancor will appear " + formatTime(timeTilEvent(startTime))
}

func deleteRaid(s []string) string {
	g := s[GUILD]
	raidType := s[RAID]

	if g != EU && g != US {
		return "You're an idiot. Only the EU and US exist, the rest of the world doesn't."
	}
	if raidType != RANCOR && raidType != TANK {
		return "Only rancors and tanks m8. The FBI might raid your home to search for those drugs tho."
	}
	guild := guilds.Guild(g)
	switch raidType {
	case RANCOR:
		guild.DeleteRancor()
	case TANK:
		guild.DeleteTank()
	}
	return fmt.Sprintf("That %s was not the one you were looking for", raidType)
}

/**
 * Raid command should be
 * !draid eu/us <type> <date>
 */
func setDefaultRaid(s string) string {
	chunks := strings.Split(s, " ")
	if len(chunks) != 4 {
		return "Wrong command passed. Please use !draid eu/us <type> <date>"
	}
	request := chunks[1:]
	guild := strings.ToLower(request[GUILD])
	if guild != EU && guild != US {
		return "You moron"
	}
	raidType := strings.ToLower(request[RAID])
	if raidType != "rancor" && raidType != "tank" {
		return "It's not that difficult mate"
	}
	date := request[DATE]
	var res string
	switch raidType {
	case RANCOR:
		res = createDefaultRancor(guild, date)
	case TANK:
		res = createDefaultTank(guild, date)
	}
	return fmt.Sprintf("%s", res)
}

func createDefaultRancor(g string, dateString string) string {
	guild := guilds.Guild(g)
	date, err := time.ParseInLocation(DATE_FORMAT, dateString, guild.Location)
	switch g {
	case EU:
		guild.SetDefaultEURancor(date)
	case US:
		guild.SetDefaultUSRancor(date)
	}
	if err != nil {
		return fmt.Sprintf("Error while creating date: %s", err.Error())
	}
	return guild.Rancor.String()
}

func createDefaultTank(g string, dateString string) string {
	guild := guilds.Guild(g)

	date, err := time.ParseInLocation(DATE_FORMAT, dateString, guild.Location)
	switch g {
	case EU:
		guild.SetDefaultEUTank(date)
	case US:
		guild.SetDefaultUSTank(date)
	}
	if err != nil {
		return fmt.Sprintf("Error while creating date: %s", err.Error())
	}
	return guild.Tank.String()
}
