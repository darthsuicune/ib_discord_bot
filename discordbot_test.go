package main

import (
	"testing"
)

func TestGuild_Save(t *testing.T) {
	addRaid(RANCOR, EU)
	savedData := guilds[EU].Save()
	if savedData == nil {
		t.Error("Wat")
	}
}

func TestGuild_Restore(t *testing.T) {

}