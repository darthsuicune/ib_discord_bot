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
	result := make([]string, 2)
	timeToStart := timeTilEvent(r.StartTime)
	timeToFfa := timeTilEvent(r.Ffa)

	if timeToFfa == "finished" {
		return []string{"The rancor has finished! Inglorious Officers, what's up with that?!"}
	}

	if timeToStart == "finished" {
		result[0] = "Rancor initial phase already started. " +
			"Remember, single toon on the first 12h, only 0 damage attacks then til FFA"
	} else {
		start := formatTime(timeToStart)
		result[0] = "Rancor " + start
	}
	end := formatTime(timeToFfa)
	result[1] = "FFA " + end
	return result
}

func (r Rancor) String() string {
	return strings.Join(r.Times(), "\n")
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

	if timeToFfa == "finished" {
		return []string{"The tank has finished! Inglorious Officers, what's up with that?!"}
	}

	var start string
	if timeToStart == "finished" {
		start = "The tank has already started."
	} else {
		start = "Tank " + formatTime(timeToStart)
	}

	ffa := "FFA " + formatTime(timeToFfa)

	if t.Phase2.Equal(t.StartTime) {
		return []string{start, ffa}
	} else {
		timeToP4 := timeTilEvent(t.Phase4)
		timeToP3 := timeTilEvent(t.Phase3)
		timeToP2 := timeTilEvent(t.Phase2)

		p2 := "Phase 2 " + formatTime(timeToP2)
		p3 := "Phase 3 " + formatTime(timeToP3)
		p4 := "Phase 4 " + formatTime(timeToP4)
		return []string{start, p2, p3, p4, ffa}
	}
}

func (t Tank) String() string {
	times := t.Times()
	return strings.Join(times, "\n")
}

func timeTilEvent(when time.Time) string {
	waitTime := when.Sub(time.Now())
	if waitTime < -5*time.Minute {
		return "finished"
	}
	wait := waitTime.String()
	minuteMark := strings.LastIndex(wait, "m") + 1
	return wait[0:minuteMark]
}

func formatTime(timeToFormat string) string {
	if timeToFormat == "finished" {
		return "already started"
	} else if timeToFormat != "" {
		return "in **" + timeToFormat + "**"
	} else {
		return "**NOW!**"
	}
}
