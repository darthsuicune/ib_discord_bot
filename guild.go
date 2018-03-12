package main

import (
	"fmt"
	"time"
	"bytes"
	"strings"
)

type Guilds struct {
	Eu *Guild
	Us *Guild
}

func (g *Guilds) Raids() string {
	return fmt.Sprintf("**%s**\n%s\n**%s**\n%s", "EU", g.Eu.Raids(), "US", g.Us.Raids())
}

func (g *Guilds) Guild(code string) *Guild {
	if strings.ToLower(code) == EU {
		return guilds.Eu
	} else if strings.ToLower(code) == US {
		return guilds.Us
	} else {
		panic("WHAT THE FUCK DID YOU DO!")
	}
}

type Guild struct {
	Rancor *Rancor
	Tank   *Tank
	Sith   *Sith
	Location *time.Location
}

/**
 * US Default times:
 * -Rancor: 20:30 EST, FFA: +24h
 * -Tank: 22:00 EST, FFA: +46h
 *
 * EU Default times:
 * -Rancor: 20:00 GMT, FFA: +24h
 * -Tank: 21:00 GMT, P2: +10h, P3: +34h, P4: +44h, FFA: +46h
 */

func usTime() *time.Location {
	l, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println(err)
	}
	return l
}

func euTime() *time.Location {
	l, err := time.LoadLocation("Europe/London")
	if err != nil {
		fmt.Println(err)
	}
	return l
}

func (g *Guild) SetDefaultUSRancor(startTime time.Time) {
	newTime := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 20, 30, 0, 0, usTime())
	g.SetRancor(newTime)
}

func (g *Guild) SetRancor(startTime time.Time) {
	g.Rancor = &Rancor{StartTime: startTime, Ffa: startTime.Add(24 * time.Hour)}
}

func (g *Guild) SetDefaultEURancor(startTime time.Time) {
	newTime := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 20, 0, 0, 0, euTime())
	g.SetRancor(newTime)
}

func (g *Guild) SetDefaultUSTank(startTime time.Time) {
	newTime := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 22, 0, 0, 0, usTime())
	g.SetUSTank(newTime)
}

func (g *Guild) SetUSTank(startTime time.Time) {
	g.Tank = &Tank{StartTime: startTime, Phase2: startTime, Phase3: startTime, Phase4: startTime, Ffa: startTime.Add(24 * time.Hour)}
}

func (g *Guild) SetDefaultEUTank(startTime time.Time) {
	newTime := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 21, 0, 0, 0, euTime())
	g.SetEUTank(newTime)
}

func (g *Guild) SetEUTank(startTime time.Time) {
	g.Tank = &Tank{StartTime: startTime, Phase2: startTime.Add(10 * time.Hour), Phase3: startTime.Add(34 * time.Hour), Phase4: startTime.Add(44 * time.Hour), Ffa: startTime.Add(46 * time.Hour)}
}

func (g *Guild) DeleteTank() {
	g.Tank = nil
}

func (g *Guild) DeleteRancor() {
	g.Rancor = nil
}

func (g *Guild) Raids() string {
	var buffer bytes.Buffer
	if g.Rancor != nil {
		buffer.WriteString(g.Rancor.String())
		buffer.WriteString("\n")
	}
	if g.Tank != nil {
		buffer.WriteString(g.Tank.String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}
