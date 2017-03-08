package main

import (
	"fmt"
	"testing"
	"time"
)

/**
 * Test the raid parser. Possible commands:
 * !raid
 * !raid eu/us
 * !draid eu/us <type> <date>
 * !raid set eu/us <type> <date> <time> <TZ>
 * !raid delete eu/us <type>
 */

var raidYear = time.Now().Add(366 * 24 * time.Hour).Year()
var raids map[string]*Guild
var euRancor = "**EU** Rancor in **23h59m**\nFFA in **47h59m**\n"
var euTank = "Tank in **23h59m**\n" +
	"Phase 2 in **47h59m**\n" +
	"Phase 3 already started\n" +
	"Phase 4 already started\n" +
	"FFA in **47h59m**\n"
var usRancor = "**US** Rancor in **23h59m**\nFFA in **47h59m**\n"
var usTank = "Tank in **23h59m**\nFFA in **47h59m**\n"

func init() {
	raids = make(map[string]*Guild)
	raids["eu"] = &Guild{}
	raids["us"] = &Guild{}
}

/**
 * DEFAULT
 */
func TestDefaultRaidEUSetsUpTheDefault(t *testing.T) {
	command := fmt.Sprintf("!draid eu rancor %d-12-13", raidYear)
	expectedBase := time.Date(raidYear, 12, 13, 20, 0, 0, 0, euTime())
	expectedFfa := expectedBase.Add(24 * time.Hour)
	base := trimTime(expectedBase.Sub(time.Now()).String())
	ffa := trimTime(expectedFfa.Sub(time.Now()).String())
	expectedResponse := fmt.Sprintf("Rancor in **%s**\nFFA in **%s**", base, ffa)

	res := setDefaultRaid(command)

	if res != expectedResponse {
		t.Error(fmt.Sprintf("Received %s, expected: %s", res, expectedResponse))
	}
}

func TestDefaultRaidUSSetsUpTheDefault(t *testing.T) {
	command := fmt.Sprintf("!draid us tank %d-12-13", raidYear)
	expectedBase := time.Date(raidYear, 12, 13, 22, 0, 0, 0, usTime())
	expectedFfa := expectedBase.Add(46 * time.Hour)
	base := trimTime(expectedBase.Sub(time.Now()).String())
	ffa := trimTime(expectedFfa.Sub(time.Now()).String())
	expectedResponse := fmt.Sprintf("Tank in **%s**\nFFA in **%s**", base, ffa)

	res := setDefaultRaid(command)

	if res != expectedResponse {
		t.Error(fmt.Sprintf("Received %s, expected: %s", res, expectedResponse))
	}
}

/**
 * NON-DEFAULT
 */
func TestRaidDisplaysAllGuildsAllAvailableRaids(t *testing.T) {
	command := "!raid"
	addRaid(RANCOR, EU)
	addRaid(RANCOR, US)
	addRaid(TANK, US)

	res := readRaidCommand(command, false)

	if res != euRancor + usRancor + usTank &&
	res != usRancor + usTank + euRancor {
		t.Error(res)
	}
}

func addRaid(raidType string, guild string) {
	startTime := time.Now().Add(24 * time.Hour)
	ffaTime := time.Now().Add(48 * time.Hour)
	if raidType == RANCOR {
		guilds[guild].Rancor = &Rancor{StartTime: startTime, Ffa: ffaTime}
	} else {
		if guild == EU {
			guilds[guild].Tank = &Tank{StartTime: startTime, Phase2: ffaTime, Ffa: ffaTime}
		} else {
			guilds[guild].Tank = &Tank{StartTime: startTime, Phase2: startTime, Ffa: ffaTime}
		}
	}
}

func TestRaidEUDisplaysAllEURaids(t *testing.T) {
	command := "!raid eu"
	addRaid(RANCOR, EU)
	addRaid(TANK, EU)

	res := readRaidCommand(command, false)
	if res != euRancor + euTank {
		t.Error(res)
	}
}

func TestRaidUSDisplaysAllUSRaids(t *testing.T) {
	command := "!raid us"
	addRaid(RANCOR, US)
	addRaid(TANK, US)

	res := readRaidCommand(command, false)
	if res != usRancor + usTank {
		t.Error(res)
	}
}

func TestRaidSetReturnsNopeWithoutPermissions(t *testing.T) {
	command := "!raid set eu tank 2012-12-13 01:02:03 GMT"

	res := readRaidCommand(command, false)

	if res != "Nope" {
		t.Error(fmt.Sprintf("NOPE: %s", res))
	}
}

func TestRaidSetsUpTheRightValue(t *testing.T) {
	command := fmt.Sprintf("!raid set eu tank %d-12-13 01:02:03 GMT", raidYear)
	expectedTiming := time.Date(raidYear, 12, 13, 1, 2, 3, 0, euTime())

	expectedP2 := expectedTiming.Add(10 * time.Hour)
	expectedP3 := expectedTiming.Add(34 * time.Hour)
	expectedP4 := expectedTiming.Add(44 * time.Hour)
	expectedFfa := expectedTiming.Add(46 * time.Hour)

	readRaidCommand(command, true)


	if !guilds[EU].Tank.StartTime.Equal(expectedTiming) {
		t.Error("Wrong start time: " + guilds[EU].Tank.StartTime.String())
	}
	if !guilds[EU].Tank.Phase2.Equal(expectedP2) {
		t.Error("Wrong phase 2 time: " + guilds[EU].Tank.Phase2.String())
	}
	if !guilds[EU].Tank.Phase3.Equal(expectedP3) {
		t.Error("Wrong phase 3 time: " + guilds[EU].Tank.Phase3.String())
	}
	if !guilds[EU].Tank.Phase4.Equal(expectedP4) {
		t.Error("Wrong phase 4 time: " + guilds[EU].Tank.Phase4.String())
	}
	if !guilds[EU].Tank.Ffa.Equal(expectedFfa) {
		t.Error("Something went very wrong. " +
			"Current info: Raid on " + guilds[EU].Tank.Ffa.String() +
			"\nExpected: " + expectedFfa.String())
	}
}
func TestUSRaidSetsUpTheRightValue(t *testing.T) {
	command := fmt.Sprintf("!raid set us tank %d-12-13 01:02:03 EST", raidYear)
	expectedTiming := time.Date(raidYear, 12, 13, 1, 2, 3, 0, usTime())
	expectedFfa := expectedTiming.Add(46 * time.Hour)

	readRaidCommand(command, true)


	if !guilds[US].Tank.StartTime.Equal(expectedTiming) {
		t.Error("Wrong start time: " + guilds[US].Tank.StartTime.String())
	}
	if !guilds[US].Tank.Ffa.Equal(expectedFfa) {
		t.Error("Something went very wrong. " +
			"Current info: Raid on " + guilds[US].Tank.Ffa.String() +
			"\nExpected: " + expectedFfa.String())
	}
}

func TestDeleteRaidRemovesTheCurrentRaid(t *testing.T) {
	addRaid(RANCOR, EU)
	addRaid(RANCOR, US)
	command := "!raid delete eu rancor"

	if guilds[EU].Rancor == nil {
		t.Error("Wait wat? EU Rancor should exist!")
	}
	if guilds[US].Rancor == nil {
		t.Error("Wait wat? US Rancor should exist!")
	}
	readRaidCommand(command, true)

	if guilds[EU].Rancor != nil {
		t.Error("Wait wat? Rancor should NOT exist!")
	}
	if guilds[US].Rancor == nil {
		t.Error("Wait wat? US Rancor should still exist!")
	}
}