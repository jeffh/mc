package protocol

import (
	. "describe"
	"testing"
)

func TestProtocolStringReader(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, "hello world")
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadString(r)
	Expect(t, err, ToBeNil)
	Expect(t, v, ToEqual, "hello world")
}

func TestProtocolSlotSliceReader(t *testing.T) {
	r, b := createProtocolReader()
	// last two int8s are just arbitrary binary bits for now
	// ID, Count, Damage, DataSize, Gzipped Data...
	err := writeBytes(b, uint16(1), int16(2), int8(100), int16(99), int16(2), int8(2), int8(3))
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadSlotSlice(r)
	Expect(t, err, ToBeNil)
	slots, ok := v.([]Slot)
	Expect(t, ok, ToBeTrue)
	Expect(t, slots, ToBeLengthOf, 1)
	slot := slots[0]
	Expect(t, slot.ID, ToEqual, int16(2))
	Expect(t, slot.Count, ToEqual, int8(100))
	Expect(t, slot.Damage, ToEqual, int16(99))
	Expect(t, slot.CompressedNBT, ToBeLengthOf, 2)
	Expect(t, b.Len(), ToEqual, 0)
}

func TestProtocolSlotReader(t *testing.T) {
	r, b := createProtocolReader()
	// last two int8s are just arbitrary binary bits for now
	// ID, Count, Damage, DataSize, Gzipped Data...
	err := writeBytes(b, int16(2), int8(100), int16(99), int16(2), int8(2), int8(3))
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadSlot(r)
	Expect(t, err, ToBeNil)
	slot, ok := v.(Slot)
	Expect(t, ok, ToBeTrue)
	Expect(t, slot.ID, ToEqual, int16(2))
	Expect(t, slot.Count, ToEqual, int8(100))
	Expect(t, slot.Damage, ToEqual, int16(99))
	Expect(t, slot.CompressedNBT, ToBeLengthOf, 2)
	Expect(t, b.Len(), ToEqual, 0)
}

func TestProtocolSlotReaderForEmptySlot(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, int16(-1))
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadSlot(r)
	Expect(t, err, ToBeNil)
	slot, ok := v.(Slot)
	Expect(t, ok, ToBeTrue)
	Expect(t, slot.ID, ToEqual, int16(-1))
	Expect(t, b.Len(), ToEqual, 0)
}

func TestProtocolSlotReaderForEmptyCompressedNBT(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, int16(2), int8(100), int16(99), int16(-1))
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadSlot(r)
	Expect(t, err, ToBeNil)
	slot, ok := v.(Slot)
	Expect(t, ok, ToBeTrue)
	Expect(t, slot.ID, ToEqual, int16(2))
	Expect(t, slot.Count, ToEqual, int8(100))
	Expect(t, slot.Damage, ToEqual, int16(99))
	Expect(t, slot.CompressedNBT, ToBeLengthOf, 0)
	Expect(t, b.Len(), ToEqual, 0)
}

func TestProtocolEntityMetadataTypeReader(t *testing.T) {
}
func TestProtocolDestroyEntityReader(t *testing.T) {
}
func TestProtocolChunkDataReader(t *testing.T) {
}
func TestProtocolMultiBlockChangeReader(t *testing.T) {
}
func TestProtocolMapChunkBulkReader(t *testing.T) {
}
func TestProtocolSetWindowItemsReader(t *testing.T) {
}
func TestProtocolItemDataReader(t *testing.T) {
}
func TestProtocolPluginMessageReader(t *testing.T) {
}
