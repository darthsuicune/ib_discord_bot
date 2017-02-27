package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var g *Guild

func init() {
	g = &Guild{}
}

func createGuild() *Guild {
	return &Guild{}
}

func TestSetDefaultUSRancor(t *testing.T) {
	timing := time.Date(2006, 1, 2, 15, 04, 05, 0, usTime())
	expectedTiming := time.Date(2006, 1, 2, 20, 30, 0, 0, usTime())
	g.SetDefaultUSRancor(timing)

	rancor := g.Rancor.Times()
	if len(rancor) != 2 {
		t.Error("baaad result: " + strings.Join(rancor, ", "))
	}
	expectedStart := trimTime(expectedTiming.Sub(time.Now()).String())
	if rancor[0] != fmt.Sprintf("Rancor in **%s**", expectedStart) {
		t.Error("Wrong start time: " + rancor[0] + ", expected: " + expectedStart)
	}
	expectedFfa := trimTime(expectedTiming.Add(24 * time.Hour).Sub(time.Now()).String())
	if rancor[1] != fmt.Sprintf("FFA in **%s**", expectedFfa) {
		t.Error("Wrong ffa time: " + rancor[1] + ", expected: " + expectedFfa)
	}
}

func TestSetUSRancor(t *testing.T) {
	timing := time.Date(2006, 1, 2, 15, 04, 05, 0, usTime())
	g.SetRancor(timing)

	rancor := g.Rancor.Times()
	if len(rancor) != 2 {
		t.Error("baaad result: " + strings.Join(rancor, ", "))
	}
	expectedStart := trimTime(timing.Sub(time.Now()).String())
	if rancor[0] != fmt.Sprintf("Rancor in **%s**", expectedStart) {
		t.Error("Wrong start time: " + rancor[0] + ", expected: " + expectedStart)
	}
	expectedFfa := trimTime(timing.Add(24 * time.Hour).Sub(time.Now()).String())
	if rancor[1] != fmt.Sprintf("FFA in **%s**", expectedFfa) {
		t.Error("Wrong ffa time: " + rancor[1] + ", expected: " + expectedFfa)
	}
}

func TestSetDefaultUSTank(t *testing.T) {
	timing := time.Date(2006, 1, 2, 15, 04, 05, 0, usTime())
	expectedTiming := time.Date(2006, 1, 2, 22, 0, 0, 0, usTime())
	g.SetDefaultUSTank(timing)

	tank := g.Tank.Times()
	if len(tank) != 2 {
		t.Error("baaad result: " + strings.Join(tank, ", "))
	}
	expectedStart := trimTime(expectedTiming.Sub(time.Now()).String())
	if tank[0] != fmt.Sprintf("Tank in **%s**", expectedStart) {
		t.Error("Wrong start time: " + tank[0] + ", expected: " + expectedStart)
	}
	expectedFfa := trimTime(expectedTiming.Add(46 * time.Hour).Sub(time.Now()).String())
	if tank[1] != fmt.Sprintf("FFA in **%s**", expectedFfa) {
		t.Error("Wrong ffa time: " + tank[1] + ", expected: " + expectedFfa)
	}
}

func trimTime(timing string) string {
	mark := strings.LastIndex(timing, "m") + 1
	return timing[:mark]
}

func TestSetUSTank(t *testing.T) {
	timing := time.Date(2006, 1, 2, 15, 04, 05, 0, usTime())
	g.SetUSTank(timing)

	tank := g.Tank.Times()
	if len(tank) != 2 {
		t.Error("baaad result: " + strings.Join(tank, ", "))
	}
	expectedStart := trimTime(timing.Sub(time.Now()).String())
	if tank[0] != fmt.Sprintf("Tank in **%s**", expectedStart) {
		t.Error("Wrong start time: " + tank[0] + ", expected: " + expectedStart)
	}
	expectedFfa := trimTime(timing.Add(46 * time.Hour).Sub(time.Now()).String())
	if tank[1] != fmt.Sprintf("FFA in **%s**", expectedFfa) {
		t.Error("Wrong ffa time: " + tank[1] + ", expected: " + expectedFfa)
	}
}

func TestSetDefaultEURancor(t *testing.T) {
	timing := time.Date(2006, 1, 2, 15, 04, 05, 0, euTime())
	expectedTiming := time.Date(2006, 1, 2, 20, 30, 0, 0, euTime())
}
func TestSetEURancor(t *testing.T) {}
func TestSetDefaultEUTank(t *testing.T) {}
func TestSetEUTank(t *testing.T) {}
