package protocol

import (
	. "describe"
	"testing"
)

func TestGameModeSurvival(t *testing.T) {
	gm := SurvivalMode
	Expect(t, gm.IsSurvival(), ToBeTrue)
	Expect(t, gm.IsCreative(), Not(ToBeTrue))
	Expect(t, gm.IsAdventure(), Not(ToBeTrue))
	Expect(t, gm.IsHardcore(), Not(ToBeTrue))
}

func TestGameModeCreative(t *testing.T) {
	gm := CreativeMode
	Expect(t, gm.IsSurvival(), Not(ToBeTrue))
	Expect(t, gm.IsCreative(), ToBeTrue)
	Expect(t, gm.IsAdventure(), Not(ToBeTrue))
	Expect(t, gm.IsHardcore(), Not(ToBeTrue))
}

func TestGameModeAdventure(t *testing.T) {
	gm := AdventureMode
	Expect(t, gm.IsSurvival(), Not(ToBeTrue))
	Expect(t, gm.IsCreative(), Not(ToBeTrue))
	Expect(t, gm.IsAdventure(), ToBeTrue)
	Expect(t, gm.IsHardcore(), Not(ToBeTrue))
}

func TestGameModeSurvivalFlag(t *testing.T) {
	gm := GameMode(AdventureMode | HardcoreModeFlag)
	Expect(t, gm.IsSurvival(), Not(ToBeTrue))
	Expect(t, gm.IsCreative(), Not(ToBeTrue))
	Expect(t, gm.IsAdventure(), ToBeTrue)
	Expect(t, gm.IsHardcore(), ToBeTrue)
}
