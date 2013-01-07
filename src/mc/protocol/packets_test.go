package protocol

import (
	"bytes"
	"compress/gzip"
	. "describe"
	"io"
	"io/ioutil"
	"testing"
)

func TestSlotCanReturnReaderOfItsTrueData(t *testing.T) {
	expectedData := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	rawReader := bytes.NewBuffer(expectedData)
	compressedOutput := bytes.NewBuffer([]byte{})
	reader := gzip.NewWriter(compressedOutput)
	_, err := io.Copy(reader, rawReader)
	reader.Close()
	Expect(t, err, ToBeNil)

	compressed, err := ioutil.ReadAll(compressedOutput)
	Expect(t, err, ToBeNil)

	s := Slot{GzippedNBT: compressed}
	r, err := s.NewReader()
	Expect(t, err, ToBeNil)

	data, err := ioutil.ReadAll(r)
	Expect(t, err, ToBeNil)
	Expect(t, data, ToEqual, expectedData)
}

func TestSlotIsEmptyIfIDIsNegOne(t *testing.T) {
	// negative -1 is empty
	s := Slot{ID: -1}
	Expect(t, s.IsEmpty(), ToBeTrue)

	s = EmptySlot // alias to one above
	Expect(t, s.IsEmpty(), ToBeTrue)

	s = Slot{ID: 1}
	Expect(t, s.IsEmpty(), Not(ToBeTrue))
}

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
