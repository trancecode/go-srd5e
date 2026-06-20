package core

import "testing"

func TestCarryingCapacity(t *testing.T) {
	// STR 15, Medium: 15*15 = 225; push/drag/lift = 450.
	if got := CarryingCapacity(15, SizeMedium); got != 225 {
		t.Errorf("CarryingCapacity Medium = %d, want 225", got)
	}
	if got := PushDragLift(15, SizeMedium); got != 450 {
		t.Errorf("PushDragLift Medium = %d, want 450", got)
	}
	// Large doubles, Tiny halves.
	if got := CarryingCapacity(15, SizeLarge); got != 450 {
		t.Errorf("CarryingCapacity Large = %d, want 450", got)
	}
	if got := CarryingCapacity(10, SizeTiny); got != 75 {
		t.Errorf("CarryingCapacity Tiny = %d, want 75", got)
	}
}

func TestEncumbranceTier(t *testing.T) {
	// STR 10: encumbered > 50 (STR*5), heavily > 100 (STR*10).
	if got := EncumbranceTier(10, 40); got != Unencumbered {
		t.Errorf("tier(40) = %v, want Unencumbered", got)
	}
	if got := EncumbranceTier(10, 60); got != Encumbered {
		t.Errorf("tier(60) = %v, want Encumbered", got)
	}
	if got := EncumbranceTier(10, 120); got != HeavilyEncumbered {
		t.Errorf("tier(120) = %v, want HeavilyEncumbered", got)
	}
}

func TestCarryingCapacityRequiresSize(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("CarryingCapacity with SizeUnspecified should panic")
		}
	}()
	CarryingCapacity(10, SizeUnspecified)
}
