package turn

import (
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

// StaticInitiative is the take-10 initiative: 10 plus the Dexterity modifier.
func StaticInitiative(dexScore core.AbilityScore) int {
	return 10 + int(core.AbilityModifier(dexScore))
}

// RollInitiative rolls a d20 (with vantage) and adds the Dexterity modifier.
func RollInitiative(dexMod core.Modifier, r dice.Roller, v dice.Vantage) int {
	return dice.RollD20(r, v).Total + int(dexMod)
}

// DifficultTerrainCost is the movement budget cost of a measured distance in
// difficult terrain: one extra foot per foot (×2).
func DifficultTerrainCost(distance core.Distance) core.Distance { return distance * 2 }

// DashDistance is the extra movement the Dash action grants: your speed again.
func DashDistance(speed core.Distance) core.Distance { return speed }

// StandUpCost is the movement cost to stand from prone: half your speed.
func StandUpCost(speed core.Distance) core.Distance { return speed / 2 }

// LongJump is how far you can leap horizontally: your Strength score in feet
// with a running start, half that from a standstill.
func LongJump(str core.AbilityScore, runningStart bool) core.Distance {
	d := core.Distance(str)
	if !runningStart {
		d /= 2
	}
	return d
}

// HighJump is how high you can leap: 3 plus your Strength modifier in feet with
// a running start, half that from a standstill.
func HighJump(strMod core.Modifier, runningStart bool) core.Distance {
	d := core.Distance(3 + int(strMod))
	if !runningStart {
		d /= 2
	}
	return d
}
