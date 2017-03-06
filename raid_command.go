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
)

func init() {
	guilds = make(map[string]*Guild)
	eu := Guild{}
	us := Guild{}
	guilds["eu"] = &eu
	guilds["us"] = &us
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
	return parseChunks(chunks[1:])
}

func fullRaidInformation() string {
	var buffer bytes.Buffer
	for _, guild := range guilds {
		buffer.Write([]byte(guild.Rancor.String()))
		buffer.Write([]byte(guild.Tank.String()))
	}
	return buffer.String()
}

func parseChunks(chunks []string) string {
	var buffer bytes.Buffer
	if len(chunks) == 1 {
		guild := strings.ToLower(chunks[GUILD])
		if guild != "eu" && guild != "us" {
			buffer.Write([]byte("You f\\*\\*\\*ing moron"))
		} else {
			buffer.Write([]byte(guilds[guild].Rancor.String()))
			buffer.Write([]byte(guilds[guild].Tank.String()))
		}
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
	if guild != "eu" && guild != "us" {
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
	
	switch(guild) {
		case "eu":
			date, err = time.ParseInLocation(DATE_FORMAT, dateString, euTime())
			guilds[guild].SetDefaultEURancor(date)
		case "us":
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
	
	switch(guild) {
		case "eu":
			date, err = time.ParseInLocation(DATE_FORMAT, dateString, euTime())
			guilds[guild].SetDefaultEUTank(date)
		case "us":
			date, err = time.ParseInLocation(DATE_FORMAT, dateString, usTime())
			guilds[guild].SetDefaultUSTank(date)
	}
	if err != nil {
		return fmt.Sprintf("Error while creating date: %s", err.Error())
	}
	return guilds[guild].Tank.String()
}
