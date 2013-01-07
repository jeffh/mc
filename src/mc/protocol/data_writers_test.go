package protocol

import (
	. "github.com/jeffh/goexpect"
	"testing"
)

func TestProtocolEntityMetadataSliceWriter(t *testing.T) {
	w, b := createProtocolWriter()
	metadata := []EntityMetadata{
		{EntityFlags, EntityMetadataByte, byte(6)},
		{EntityDrowning, EntityMetadataShort, int16(42)},
		{EntityUnderPotionFX, EntityMetadataInt, int32(4432)},
		{EntityAnimalCounter, EntityMetadataFloat, float32(0.5)},
		{EntityState1, EntityMetadataString, "There is no spoon"},
		{EntityState2, EntityMetadataSlot, EmptySlot},
		{EntityState3, EntityMetadataPosition, Position{2, 4, 6}},
	}

	err := ProtocolWriteEntityMetadataSlice(w, metadata)
	Expect(t, err, ToBeNil)

	Expect(t, b.Bytes(), ToEqualBytes,
		// these don't represent actual data types
		// but we're just testing our flexibility to parse everything
		entityKey(EntityFlags, EntityMetadataByte), byte(6),
		entityKey(EntityDrowning, EntityMetadataShort), int16(42),
		entityKey(EntityUnderPotionFX, EntityMetadataInt), int32(4432),
		entityKey(EntityAnimalCounter, EntityMetadataFloat), float32(0.5),
		entityKey(EntityState1, EntityMetadataString), "There is no spoon",
		entityKey(EntityState2, EntityMetadataSlot), int16(-1),
		entityKey(EntityState3, EntityMetadataPosition), int32(2), int32(4), int32(6),
		byte(127))
}

func TestProtocolEntityMetadataEmptySliceWriter(t *testing.T) {
	w, b := createProtocolWriter()
	metadata := []EntityMetadata{}

	err := ProtocolWriteEntityMetadataSlice(w, metadata)
	Expect(t, err, ToBeNil)

	var term byte
	err = readBytes(b, &term)
	Expect(t, err, ToBeNil)
	Expect(t, term, ToEqual, byte(127))
}

func TestProtocolSlotSliceWriter(t *testing.T) {
	w, b := createProtocolWriter()
	slot := Slot{
		ID:         1,
		Count:      5,
		Damage:     10,
		GzippedNBT: []byte{},
	}

	err := ProtocolWriteSlotSlice(w, []Slot{slot})
	Expect(t, err, ToBeNil)

	var size, id, dmg, dataTerm int16
	var count int8
	err = readBytes(b, &size, &id, &count, &dmg, &dataTerm)
	Expect(t, id, ToEqual, int16(1))
	Expect(t, id, ToEqual, int16(1))
	Expect(t, count, ToEqual, int8(5))
	Expect(t, dmg, ToEqual, int16(10))
	Expect(t, dataTerm, ToEqual, int16(-1))
	Expect(t, b.Len(), ToEqual, 0)
}

func TestProtocolSlotWriterForEmptySlot(t *testing.T) {
	w, b := createProtocolWriter()
	err := ProtocolWriteSlot(w, EmptySlot)
	Expect(t, err, ToBeNil)

	var id int16
	err = readBytes(b, &id)
	Expect(t, err, ToBeNil)
	Expect(t, id, ToEqual, int16(-1))
	Expect(t, b.Len(), ToEqual, 0)
}

func TestProtocolSlotWriterForEmptyGzippedNBT(t *testing.T) {
	// zero-length data should write -1 instead of 0 for size
	w, b := createProtocolWriter()
	slot := Slot{
		ID:         1,
		Count:      5,
		Damage:     10,
		GzippedNBT: []byte{},
	}

	err := ProtocolWriteSlot(w, slot)
	Expect(t, err, ToBeNil)
	var id, dmg int16
	var count int8
	var dataTerm int16
	err = readBytes(b, &id, &count, &dmg, &dataTerm)
	Expect(t, err, ToBeNil)
	Expect(t, id, ToEqual, int16(1))
	Expect(t, count, ToEqual, int8(5))
	Expect(t, dmg, ToEqual, int16(10))
	Expect(t, dataTerm, ToEqual, int16(-1))
	Expect(t, b.Len(), ToEqual, 0)
}

func TestProtocolBoolWriter(t *testing.T) {
	w, buf := createProtocolWriter()
	value := true
	err := ProtocolWriteBool(w, value)
	Expect(t, err, ToBeNil)

	var result byte
	err = readBytes(buf, &result)
	Expect(t, err, ToBeNil)
	Expect(t, result, ToEqual, byte(1))
}

func TestProtocolStringWriter(t *testing.T) {
	w, buf := createProtocolWriter()
	data := "foobar"
	err := ProtocolWriteString(w, data)
	Expect(t, err, ToBeNil)

	var result string
	err = readBytes(buf, &result)
	Expect(t, err, ToBeNil)
	Expect(t, result, ToEqual, data)
}

func TestProtocolBytesWriter(t *testing.T) {
	w, buf := createProtocolWriter()
	data := []byte{1, 2, 3}
	err := ProtocolWriteByteSlice(w, data)
	Expect(t, err, ToBeNil)

	var size int16
	err = readBytes(buf, &size)
	Expect(t, err, ToBeNil)
	Expect(t, size, ToEqual, int16(len(data)))
	for _, ch := range data {
		var b byte
		err = readBytes(buf, &b)
		Expect(t, err, ToBeNil)
		Expect(t, ch, ToEqual, b)
	}
}
