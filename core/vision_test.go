package core

import "testing"

func TestEffectiveLight(t *testing.T) {
	// Darkvision in range: dark -> dim, dim -> bright.
	if got := EffectiveLight(LightDark, VisionDarkvision, true); got != LightDim {
		t.Errorf("darkvision dark = %v, want Dim", got)
	}
	if got := EffectiveLight(LightDim, VisionDarkvision, true); got != LightBright {
		t.Errorf("darkvision dim = %v, want Bright", got)
	}
	// Out of range: unchanged.
	if got := EffectiveLight(LightDark, VisionDarkvision, false); got != LightDark {
		t.Errorf("darkvision out of range = %v, want Dark", got)
	}
	// Blindsight sees regardless: treat as bright.
	if got := EffectiveLight(LightDark, VisionBlindsight, true); got != LightBright {
		t.Errorf("blindsight = %v, want Bright", got)
	}
}

func TestSightVisibility(t *testing.T) {
	if got := SightVisibility(LightBright, false, false); got != VisibilityClear {
		t.Errorf("bright = %v, want Clear", got)
	}
	if got := SightVisibility(LightDim, false, false); got != VisibilityObscured {
		t.Errorf("dim = %v, want Obscured", got)
	}
	if got := SightVisibility(LightDark, false, false); got != VisibilityBlocked {
		t.Errorf("dark = %v, want Blocked", got)
	}
	if got := SightVisibility(LightBright, true, false); got != VisibilityBlocked {
		t.Errorf("blinded = %v, want Blocked", got)
	}
	if got := SightVisibility(LightBright, false, true); got != VisibilityBlocked {
		t.Errorf("invisible target = %v, want Blocked", got)
	}
}
