package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var g *Guild
var testYear = time.Now().Year() + 1
var pastYear = 2015

func init() {
	g = createGuild()
}

func createGuild() *Guild {
	return &Guild{}
}

func TestTimeZones(t *testing.T) {
	british, err := time.LoadLocation("Europe/London")
	if err != nil {
		t.Error(err.Error())
	}
	eu, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Error(err.Error())
	}
	timing := time.Date(testYear, 2, 3, 4, 5, 6, 7, eu).Add(1 * time.Hour)
	expected := time.Date(testYear,2,3,4,5,6,7, british)
	if !timing.Equal(expected) {
		t.Error("Times aren't equal. Got: " + timing.String() + ", expected: " + expected.String())
	}
}

func TestRaidFinished(t *testing.T) {
	timing := time.Date(pastYear, 1, 2, 15, 04, 05, 0, usTime())

	g.SetDefaultUSRancor(timing)

	r := g.Rancor.Times()
	if len(r) > 1 {
		t.Error("Nope")
	} else {
		if !strings.Contains(r[0], "finished") {
			t.Error("Tampocoe")
		}
	}
}

func TestSetDefaultUSRancor(t *testing.T) {
	timing := time.Date(testYear, 1, 2, 15, 04, 05, 0, usTime())
	expectedTiming := time.Date(testYear, 1, 2, 20, 30, 0, 0, usTime())
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

func trimTime(timing string) string {
	mark := strings.LastIndex(timing, "m") + 1
	return timing[:mark]
}

func TestSetUSRancor(t *testing.T) {
	timing := time.Date(testYear, 1, 2, 15, 04, 05, 0, usTime())
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
	timing := time.Date(testYear, 1, 2, 15, 04, 05, 0, usTime())
	expectedTiming := time.Date(testYear, 1, 2, 22, 0, 0, 0, usTime())
	g.SetDefaultUSTank(timing)

	tank := g.Tank.Times()
	if len(tank) != 2 {
		t.Error("baaad result: " + strings.Join(tank, ", "))
	}
	expectedStart := trimTime(expectedTiming.Sub(time.Now()).String())
	if tank[0] != fmt.Sprintf("Tank in **%s**", expectedStart) {
		t.Error("Wrong start time: " + tank[0] + ", expected: " + expectedStart)
	}
	expectedFfa := trimTime(expectedTiming.Add(24 * time.Hour).Sub(time.Now()).String())
	if tank[1] != fmt.Sprintf("FFA in **%s**", expectedFfa) {
		t.Error("Wrong ffa time: " + tank[1] + ", expected: " + expectedFfa)
	}
}

func TestSetUSTank(t *testing.T) {
	timing := time.Date(testYear, 1, 2, 15, 04, 05, 0, usTime())
	g.SetUSTank(timing)

	tank := g.Tank.Times()
	if len(tank) != 2 {
		t.Error("baaad result: " + strings.Join(tank, ", "))
	}
	expectedStart := trimTime(timing.Sub(time.Now()).String())
	if tank[0] != fmt.Sprintf("Tank in **%s**", expectedStart) {
		t.Error("Wrong start time: " + tank[0] + ", expected: " + expectedStart)
	}
	expectedFfa := trimTime(timing.Add(24 * time.Hour).Sub(time.Now()).String())
	if tank[1] != fmt.Sprintf("FFA in **%s**", expectedFfa) {
		t.Error("Wrong ffa time: " + tank[1] + ", expected: " + expectedFfa)
	}
}

func TestSetDefaultEURancor(t *testing.T) {
	timing := time.Date(testYear, 1, 2, 15, 04, 05, 0, euTime())
	expectedTiming := time.Date(testYear, 1, 2, 20, 0, 0, 0, euTime())
	g.SetDefaultEURancor(timing)
	
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

func TestSetEURancor(t *testing.T) {
	timing := time.Date(testYear, 1, 2, 15, 04, 05, 0, euTime())
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

func TestSetDefaultEUTank(t *testing.T) {
	timing := time.Date(testYear, 1, 2, 15, 04, 05, 0, euTime())
	expectedStartTime := time.Date(testYear, 1, 2, 21, 0, 0, 0, euTime())
	expectedP2Time := expectedStartTime.Add(10*time.Hour)
	expectedP3Time := expectedStartTime.Add(34*time.Hour)
	expectedP4Time := expectedStartTime.Add(44*time.Hour)
	expectedFfaTime := expectedStartTime.Add(46*time.Hour)
	g.SetDefaultEUTank(timing)

	tank := g.Tank.Times()
	if len(tank) != 5 {
		t.Error("baaad result: " + strings.Join(tank, ", "))
	}
	expectedStart := trimTime(expectedStartTime.Sub(time.Now()).String())
	if tank[0] != fmt.Sprintf("Tank in **%s**", expectedStart) {
		t.Error("Wrong start time: " + tank[0] + ", expected: " + expectedStart)
	}
	expectedP2 := trimTime(expectedP2Time.Sub(time.Now()).String())
	if tank[1] != fmt.Sprintf("Phase 2 in **%s**", expectedP2) {
		t.Error("Wrong phase 2 time: " + tank[1] + ", expected: " + expectedP2)
	}
	expectedP3 := trimTime(expectedP3Time.Sub(time.Now()).String())
	if tank[2] != fmt.Sprintf("Phase 3 in **%s**", expectedP3) {
		t.Error("Wrong phase 3 time: " + tank[2] + ", expected: " + expectedP3)
	}
	expectedP4 := trimTime(expectedP4Time.Sub(time.Now()).String())
	if tank[3] != fmt.Sprintf("Phase 4 in **%s**", expectedP4) {
		t.Error("Wrong phase 4 time: " + tank[3] + ", expected: " + expectedP4)
	}
	expectedFfa := trimTime(expectedFfaTime.Sub(time.Now()).String())
	if tank[4] != fmt.Sprintf("FFA in **%s**", expectedFfa) {
		t.Error("Wrong ffa time: " + tank[4] + ", expected: " + expectedFfa)
	}
}

func TestSetEUTank(t *testing.T) {
	timing := time.Date(testYear, 1, 2, 15, 04, 05, 0, euTime())
	expectedP2Time := timing.Add(10*time.Hour)
	expectedP3Time := timing.Add(34*time.Hour)
	expectedP4Time := timing.Add(44*time.Hour)
	expectedFfaTime := timing.Add(46*time.Hour)
	g.SetEUTank(timing)

	tank := g.Tank.Times()
	if len(tank) != 5 {
		t.Error("baaad result: " + strings.Join(tank, ", "))
	}
	expectedStart := trimTime(timing.Sub(time.Now()).String())
	if tank[0] != fmt.Sprintf("Tank in **%s**", expectedStart) {
		t.Error("Wrong start time: " + tank[0] + ", expected: " + expectedStart)
	}
	expectedP2 := trimTime(expectedP2Time.Sub(time.Now()).String())
	if tank[1] != fmt.Sprintf("Phase 2 in **%s**", expectedP2) {
		t.Error("Wrong phase 2 time: " + tank[1] + ", expected: " + expectedP2)
	}
	expectedP3 := trimTime(expectedP3Time.Sub(time.Now()).String())
	if tank[2] != fmt.Sprintf("Phase 3 in **%s**", expectedP3) {
		t.Error("Wrong phase 3 time: " + tank[2] + ", expected: " + expectedP3)
	}
	expectedP4 := trimTime(expectedP4Time.Sub(time.Now()).String())
	if tank[3] != fmt.Sprintf("Phase 4 in **%s**", expectedP4) {
		t.Error("Wrong phase 4 time: " + tank[3] + ", expected: " + expectedP4)
	}
	expectedFfa := trimTime(expectedFfaTime.Sub(time.Now()).String())
	if tank[4] != fmt.Sprintf("FFA in **%s**", expectedFfa) {
		t.Error("Wrong ffa time: " + tank[4] + ", expected: " + expectedFfa)
	}
}
