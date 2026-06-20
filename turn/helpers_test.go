package turn

import (
	"testing"

	"github.com/trancecode/go-srd5e/dice"
)

func TestInitiative(t *testing.T) {
	// Dex 14 (+2): static initiative 12.
	if StaticInitiative(14) != 12 {
		t.Errorf("StaticInitiative(14) = %d, want 12", StaticInitiative(14))
	}
	// rolled d20 shows 15, +2 mod = 17.
	if got := RollInitiative(2, dice.Constant(15), dice.VantageNone); got != 17 {
		t.Errorf("RollInitiative = %d, want 17", got)
	}
}

func TestMovementRules(t *testing.T) {
	if DifficultTerrainCost(10) != 20 {
		t.Errorf("difficult terrain = %d, want 20", DifficultTerrainCost(10))
	}
	if DashDistance(30) != 30 {
		t.Errorf("dash = %d, want 30", DashDistance(30))
	}
	if StandUpCost(30) != 15 {
		t.Errorf("stand up = %d, want 15", StandUpCost(30))
	}
}

func TestJumps(t *testing.T) {
	// long jump: STR score feet with run-up, half without.
	if LongJump(15, true) != 15 || LongJump(15, false) != 7 {
		t.Errorf("long jump = %d/%d, want 15/7", LongJump(15, true), LongJump(15, false))
	}
	// high jump: 3 + STR mod feet with run-up, half without.
	if HighJump(3, true) != 6 || HighJump(3, false) != 3 {
		t.Errorf("high jump = %d/%d, want 6/3", HighJump(3, true), HighJump(3, false))
	}
}
