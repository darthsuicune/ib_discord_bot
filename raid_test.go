package main

import (
	"strings"
	"testing"
	"time"
)

func TestPastStartTime(t *testing.T) {
	r := Rancor{StartTime: time.Now().Add(-10 * time.Minute), Ffa: time.Now().Add(24 * time.Hour)}
	ti := r.Times()
	if ti[0] != "Rancor initial phase already started. Remember, single toon on the first 12h, only 0 damage attacks then til FFA" {
		t.Error("wrong start time: " + ti[0])
	}
}

func TestPastFfaTime(t *testing.T) {
	r := Rancor{StartTime: time.Now().Add(-10 * time.Minute).Add(-24 * time.Hour), Ffa: time.Now().Add(-10 * time.Minute)}
	ti := r.Times()
	if ti[0] != "The rancor has finished! Inglorious Officers, what's up with that?!" {
		t.Error("wrong end time: " + ti[0])
	}
}

func TestTankPastStartTime(t *testing.T) {
	r := Tank{StartTime: time.Now().Add(-10 * time.Minute), Ffa: time.Now().Add(24 * time.Hour)}
	ti := r.Times()
	if ti[0] != "The tank has already started." {
		t.Error("wrong start time: " + ti[0])
	}
}

func TestTankPastFfaTime(t *testing.T) {
	r := Tank{StartTime: time.Now().Add(-10 * time.Minute).Add(-24 * time.Hour), Ffa: time.Now().Add(-10 * time.Minute)}
	ti := r.Times()
	if ti[0] != "The tank has finished! Inglorious Officers, what's up with that?!" {
		t.Error("wrong end time: " + ti[0])
	}
}

func TestRancorReturnsProperTiming(t *testing.T) {
	r := createRancor()
	ti := r.Times()
	if len(ti) != 2 {
		t.Error("Invalid data returned: " + strings.Join(ti, ", "))
	}
	if ti[0] != "Rancor **NOW!**" {
		t.Error("wrong start time: " + ti[0])
	}
	if ti[1] != "FFA in **23h59m**" {
		t.Error("wrong end time: " + ti[1])
	}
}

func createRancor() Raid {
	return Rancor{StartTime: time.Now(), Ffa: time.Now().Add(24 * time.Hour)}
}

func TestTankReturnsProperTiming(t *testing.T) {
	r := createEUTank()
	ti := r.Times()
	if len(ti) != 5 {
		t.Error("Invalid data returned: " + strings.Join(ti, ", "))
	}
	if ti[0] != "Tank **NOW!**" {
		t.Error("Wrong start time: " + ti[0])
	}
	if ti[1] != "Phase 2 in **5h59m**" {
		t.Error("Wrong phase 2 time: " + ti[1])
	}
	if ti[2] != "Phase 3 in **11h59m**" {
		t.Error("Wrong phase 3 time: " + ti[2])
	}
	if ti[3] != "Phase 4 in **17h59m**" {
		t.Error("Wrong phase 4 time: " + ti[3])
	}
	if ti[4] != "FFA in **45h59m**" {
		t.Error("Wrong ffa time: " + ti[4])
	}
}

func createEUTank() Raid {
	t := time.Now()
	return Tank{StartTime: t, Phase2: t.Add(6 * time.Hour), Phase3: t.Add(12 * time.Hour), Phase4: t.Add(18 * time.Hour), Ffa: t.Add(46 * time.Hour)}
}

func TestForUSReturnsOnlyStartAndFfa(t *testing.T) {
	r := createUSTank()

	ti := r.Times()
	if len(ti) != 2 {
		t.Error("Invalid data " + strings.Join(ti, ", "))
	}
}

func createUSTank() Raid {
	t := time.Now()
	return Tank{StartTime: t, Phase2: t, Phase3: t, Phase4: t, Ffa: t.Add(46 * time.Hour)}
}
