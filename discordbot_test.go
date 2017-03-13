package main

import (
	"testing"
	"io/ioutil"
	"encoding/json"
)

func TestSaveState(t *testing.T) {
	addRaid(RANCOR,EU)
	addRaid(TANK,US)
	err := saveState()
	if err != nil {
		t.Error(err.Error())
	}
	res, err:= ioutil.ReadFile(RAIDS_FILE_NAME)
	if err != nil {
		t.Error(err.Error())
	}
	expected, err := json.Marshal(guilds)
	if err != nil {
		t.Error(err.Error())
	}
	if string(res) != string(expected) {
		t.Error(err.Error() + "Result: " + string(res) + "Expected: " + string(expected))
	}
}

func TestSaveRestoreState(t *testing.T) {
	addRaid(RANCOR,EU)
	addRaid(TANK,US)
	err := saveState()
	if err != nil {
		t.Error(err.Error())
	}
	guilds.Eu.DeleteRancor()
	guilds.Us.DeleteTank()
	loadSavedState()
	if guilds == nil || guilds.Eu == nil || guilds.Us == nil || guilds.Eu.Rancor == nil || guilds.Us.Tank == nil {
		t.Error("Didn't recover the data properly")
	}
}