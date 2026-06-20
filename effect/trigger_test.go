package effect

import "testing"

func boolPtr(b bool) *bool { return &b }

func triggerMatches(tr Trigger, o Outcome) bool {
	effs := []ConditionalEffect{{Trigger: tr, Effect: Effect{HalfOnSave: true}}}
	return len(Triggered(effs, o)) == 1
}

func TestTriggered(t *testing.T) {
	hit := Outcome{Hit: true}
	crit := Outcome{Hit: true, Crit: true}
	miss := Outcome{}
	saveFail := Outcome{Saved: boolPtr(false)}
	saveOk := Outcome{Saved: boolPtr(true)}
	contestWin := Outcome{ContestWon: boolPtr(true)}
	contestLose := Outcome{ContestWon: boolPtr(false)}

	cases := []struct {
		tr   Trigger
		o    Outcome
		want bool
	}{
		{OnHit, hit, true},
		{OnHit, crit, true}, // a crit is a hit
		{OnHit, miss, false},
		{OnCrit, crit, true},
		{OnCrit, hit, false},
		{OnMiss, miss, true},
		{OnMiss, hit, false},
		{OnMiss, saveFail, false}, // a save outcome is not an attack miss
		{OnSaveFail, saveFail, true},
		{OnSaveFail, saveOk, false},
		{OnSaveSuccess, saveOk, true},
		{OnSave, saveFail, true},
		{OnSave, saveOk, true},
		{OnSave, hit, false},
		{OnContestWin, contestWin, true},
		{OnContestLose, contestLose, true},
		{OnContestWin, contestLose, false},
		{Always, miss, true},
		{Always, hit, true},
		{TriggerUnspecified, hit, false},
	}
	for _, c := range cases {
		if got := triggerMatches(c.tr, c.o); got != c.want {
			t.Errorf("trigger %d outcome %+v = %v, want %v", c.tr, c.o, got, c.want)
		}
	}
}

func TestTriggeredCollectsInOrder(t *testing.T) {
	effs := []ConditionalEffect{
		{Trigger: OnHit, Effect: Effect{HalfOnSave: true}},
		{Trigger: OnMiss, Effect: Effect{}},
		{Trigger: Always, Effect: Effect{HalfOnSave: true}},
	}
	got := Triggered(effs, Outcome{Hit: true})
	if len(got) != 2 { // OnHit + Always; OnMiss excluded
		t.Fatalf("got %d effects, want 2", len(got))
	}
	if !got[0].HalfOnSave || !got[1].HalfOnSave {
		t.Error("collected effects not in source order")
	}
}
