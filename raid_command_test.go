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

/**
 * DEFAULT
 */
func TestDefaultRaidEUSetsUpTheDefault(t *testing.T) {
	command := "!draid eu rancor 2012-12-13"
	expectedBase := time.Date(2012,12,13,20,0,0,0,euTime())
	expectedFfa := expectedBase.Add(24*time.Hour)
	base := trimTime(expectedBase.Sub(time.Now()).String())
	ffa := trimTime(expectedFfa.Sub(time.Now()).String())
	expectedResponse := fmt.Sprintf("Rancor in **%s**\nFFA in **%s**", base, ffa)
	
	res := setDefaultRaid(command)
	if(res != expectedResponse) {
		t.Error(fmt.Sprintf("Received %s, expected: %s", res, expectedResponse))
	}
}

func TestDefaultRaidUSSetsUpTheDefault(t *testing.T) {
	command := "!draid us tank 2012-12-13"
	expectedBase := time.Date(2012,12,13,22,0,0,0,usTime())
	expectedFfa := expectedBase.Add(46*time.Hour)
	base := trimTime(expectedBase.Sub(time.Now()).String())
	ffa := trimTime(expectedFfa.Sub(time.Now()).String())
	expectedResponse := fmt.Sprintf("Tank in **%s**\nFFA in **%s**", base, ffa)
	
	res := setDefaultRaid(command)
	if(res != expectedResponse) {
		t.Error(fmt.Sprintf("Received %s, expected: %s", res, expectedResponse))
	}
}

/**
 * NON-DEFAULT
 */
func TestRaidDisplaysAllGuildsAllRaids(t *testing.T) {
	command := "!raid"
	
	readRaidCommand(command, false)
}

func TestRaidEUDisplaysAllEURaids(t *testing.T) {
	command := "!raid eu"
	
	readRaidCommand(command, false)
}

func TestRaidUSDisplaysAllUSRaids(t *testing.T) {
	command := "!raid us"
	
	readRaidCommand(command, false)
}

func TestRaidSetReturnsNopeWithoutPermissions(t *testing.T) {
	command := "!raid set eu tank 2012-12-13 01:02:03 GMT"
	
	res := readRaidCommand(command, false)
	if res != "Nope" {
		t.Error(fmt.Sprintf("NOPE: %s", res))
	}
}

func TestRaidSetsUpTheRightValue(t *testing.T) {
	command := "!raid set eu tank 2012-12-13 01:02:03 GMT"
	
	readRaidCommand(command, true)
}

func TestDeleteRaidRemovesTheCurrentRaid(t *testing.T) {
	command := "!raid delete eu tank"
	
	readRaidCommand(command, true)
}
