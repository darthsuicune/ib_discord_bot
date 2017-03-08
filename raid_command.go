package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

var guilds map[string]*Guild

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

func init() {
	guilds = make(map[string]*Guild)
	guilds[EU] = &Guild{}
	guilds[US] = &Guild{}
}

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
		return fullRaidInformation()
	}
	if len(chunks) == 2 {
		var buffer bytes.Buffer
		guild := strings.ToLower(chunks[1])
		if guild != EU && guild != US {
			buffer.WriteString("You f\\*\\*\\*ing moron")
		} else {
			buffer.WriteString(fmt.Sprintf("**%s** %s", strings.ToUpper(guild), guilds[guild].Raids()))
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

func fullRaidInformation() string {
	var buffer bytes.Buffer
	for guildCode, guild := range guilds {
		buffer.WriteString(fmt.Sprintf("**%s** ", strings.ToUpper(guildCode)))
		buffer.WriteString(guild.Raids())
	}
	return buffer.String()
}

func setCustomRaid(s []string) string {
	guild := s[GUILD]
	raidType := s[RAID]
	date := s[DATE]
	timing := s[TIME]
	tmz := s[TMZ]

	if guild != EU && guild != US {
		return "You're an idiot. Only the EU and US exist, the rest of the world doesn't."
	}
	if raidType != RANCOR && raidType != TANK {
		return "Only rancors and tanks m8. The FBI might raid your home to search for those drugs tho."
	}
	fullTiming := fmt.Sprintf("%s %s %s", date, timing, tmz)

	var location *time.Location
	switch guild {
	case EU:
		location = euTime()
	case US:
		location = usTime()
	}
	startTime, err := time.ParseInLocation("2006-01-02 03:04:05 MST", fullTiming, location)
	if err != nil {return err.Error()}
	switch raidType {
	case RANCOR:
		guilds[guild].SetRancor(startTime)
	case TANK:
		if guild == EU {
			guilds[guild].SetEUTank(startTime)
		} else {
			guilds[guild].SetUSTank(startTime)
		}
	}
	return "Get your lightswords ready. The Rancor will appear " + formatTime(timeTilEvent(startTime))
}

func deleteRaid(s []string) string {
	guild := s[GUILD]
	raidType := s[RAID]

	if guild != EU && guild != US {
		return "You're an idiot. Only the EU and US exist, the rest of the world doesn't."
	}
	if raidType != RANCOR && raidType != TANK {
		return "Only rancors and tanks m8. The FBI might raid your home to search for those drugs tho."
	}
	switch raidType {
	case RANCOR:
		guilds[guild].DeleteRancor()
	case TANK:
		guilds[guild].DeleteTank()
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
	case "rancor":
		res = createDefaultRancor(guild, date)
	case "tank":
		res = createDefaultTank(guild, date)
	}
	return fmt.Sprintf("%s", res)
}

func createDefaultRancor(guild string, dateString string) string {
	var date time.Time
	var err error

	switch guild {
	case EU:
		date, err = time.ParseInLocation(DATE_FORMAT, dateString, euTime())
		guilds[guild].SetDefaultEURancor(date)
	case US:
		date, err = time.ParseInLocation(DATE_FORMAT, dateString, usTime())
		guilds[guild].SetDefaultUSRancor(date)
	}
	if err != nil {
		return fmt.Sprintf("Error while creating date: %s", err.Error())
	}
	return guilds[guild].Rancor.String()
}

func createDefaultTank(guild string, dateString string) string {
	var date time.Time
	var err error

	switch guild {
	case EU:
		date, err = time.ParseInLocation(DATE_FORMAT, dateString, euTime())
		guilds[guild].SetDefaultEUTank(date)
	case US:
		date, err = time.ParseInLocation(DATE_FORMAT, dateString, usTime())
		guilds[guild].SetDefaultUSTank(date)
	}
	if err != nil {
		return fmt.Sprintf("Error while creating date: %s", err.Error())
	}
	return guilds[guild].Tank.String()
}
