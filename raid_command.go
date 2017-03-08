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
	RAID = 1
	DATE = 2
	TIME = 3
	TMZ = 4
	
	DATE_FORMAT = "2006-01-02"
	
	EU = "eu"
	US = "us"

	RANCOR = "rancor"
	TANK = "tank"
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
	if len(chunks) > 2 {
		if !canSetRaids {
			return "Nope"
		}
		var buffer bytes.Buffer
		switch chunks[2] {
		case "set":

		case "delete":

		default:
			return "Your training is not complete yet. You must go to Dagobah."
		}
		return buffer.String()
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
