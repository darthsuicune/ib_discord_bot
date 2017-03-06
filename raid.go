package main

import (
	"strings"
	"time"
)

type Raid interface {
	Times() []string
	String() string
}

type Rancor struct {
	StartTime time.Time
	Ffa       time.Time
}

func (r Rancor) Times() []string {
	timeToStart := timeTilEvent(r.StartTime)
	timeToEnd := timeTilEvent(r.Ffa)

	start := formatTime(timeToStart)
	end := formatTime(timeToEnd)
	res := []string{"Rancor " + start, "FFA " + end}
	return res
}

func (r Rancor) String() string {
	return strings.Join(r.Times(),"\n")
}

type Tank struct {
	StartTime time.Time
	Phase2    time.Time
	Phase3    time.Time
	Phase4    time.Time
	Ffa       time.Time
}

func (t Tank) Times() []string {
	timeToStart := timeTilEvent(t.StartTime)
	timeToFfa := timeTilEvent(t.Ffa)

	start := "Tank " + formatTime(timeToStart)
	ffa := "FFA " + formatTime(timeToFfa)

	if t.Phase2.Equal(t.StartTime) {
		return []string{start, ffa}
	} else {
		timeToP2 := timeTilEvent(t.Phase2)
		timeToP3 := timeTilEvent(t.Phase3)
		timeToP4 := timeTilEvent(t.Phase4)

		p2 := "Phase 2 " + formatTime(timeToP2)
		p3 := "Phase 3 " + formatTime(timeToP3)
		p4 := "Phase 4 " + formatTime(timeToP4)
		return []string{start, p2, p3, p4, ffa}
	}
}

func (t Tank) String() string {
	return strings.Join(t.Times(),"\n")
}

func timeTilEvent(when time.Time) string {
	waitTime := when.Sub(time.Now()).String()
	minuteMark := strings.LastIndex(waitTime, "m") + 1
	return waitTime[0:minuteMark]
}

func formatTime(timeToFormat string) string {
	if timeToFormat != "" {
		return "in **" + timeToFormat + "**"
	} else {
		return "**NOW!**"
	}
}
