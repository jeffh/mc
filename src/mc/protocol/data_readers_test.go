package protocol

import (
	. "github.com/jeffh/goexpect"
	"testing"
)

func TestProtocolInt32PrefixedBytes(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b,
		int32(4),
		byte(1), byte(2), byte(3), byte(4),
	)
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadInt32PrefixedBytes(r)
	Expect(t, err, ToBeNil)
	bytes, ok := v.(Int32PrefixedBytes)
	Expect(t, ok, ToBeTrue)
	Expect(t, bytes, ToEqual, Int32PrefixedBytes{1, 2, 3, 4})
}

func TestProtocolEntityMetadataSliceReaderShouldParseBasicTypes(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b,
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
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadEntityMetadataSlice(r)
	Expect(t, err, ToBeNil)
	slice, ok := v.([]EntityMetadata)
	Expect(t, ok, ToBeTrue)
	Expect(t, slice, ToEqual, []EntityMetadata{
		{EntityFlags, EntityMetadataByte, byte(6)},
		{EntityDrowning, EntityMetadataShort, int16(42)},
		{EntityUnderPotionFX, EntityMetadataInt, int32(4432)},
		{EntityAnimalCounter, EntityMetadataFloat, float32(0.5)},
		{EntityState1, EntityMetadataString, "There is no spoon"},
		{EntityState2, EntityMetadataSlot, EmptySlot},
		{EntityState3, EntityMetadataPosition, Position{2, 4, 6}},
	})
	Expect(t, b.Len(), ToBe, 0)
}

func TestProtocolEntityMetadataSliceReaderShouldStopOn127Byte(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, byte(127))
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadEntityMetadataSlice(r)
	Expect(t, err, ToBeNil)
	slice, ok := v.([]EntityMetadata)
	Expect(t, ok, ToBeTrue)
	Expect(t, slice, ToBeLengthOf, 0)
	Expect(t, b.Len(), ToBe, 0)
}

func TestProtocolMapChunkBulkReader(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b,
		int16(2),   // count
		int32(1),   // data size
		byte(1),    // skylight sent
		byte(1),    // data
		int32(10),  // (0) metadata.X
		int32(11),  // (0) metadata.Z
		uint16(15), // (0) metadata.PrimaryBitmap
		uint16(16), // (0) metadata.AddBitmap
		int32(12),  // (1) metadata.X
		int32(13),  // (1) metadata.Z
		uint16(17), // (1) metadata.PrimaryBitmap
		uint16(18), // (1) metadata.AddBitmap
	)
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadMapChunkBulk(r)
	Expect(t, err, ToBeNil)
	chunk := v.(MapChunkBulk)
	Expect(t, chunk.CompressedData, ToEqual, []byte{1})
	Expect(t, chunk.SkylightSent, ToBeTrue)
	Expect(t, chunk.Metadatas, ToEqual, []ChunkBulkMetadata{
		{10, 11, 15, 16},
		{12, 13, 17, 18},
	})
}

func TestProtocolStringReader(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, "hello world")
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadString(r)
	Expect(t, err, ToBeNil)
	Expect(t, v, ToBe, "hello world")
	Expect(t, b.Len(), ToBe, 0)
}

func TestProtocolStringSliceReader(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, int16(2), "John", "Doe")
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadStringSlice(r)
	Expect(t, err, ToBeNil)
	names, ok := v.([]string)
	Expect(t, ok, ToBeTrue)
	Expect(t, names, ToEqual, []string{"John", "Doe"})
}

func TestProtocolSlotSliceReader(t *testing.T) {
	r, b := createProtocolReader()
	// last two int8s are just arbitrary binary bits for now
	// ID, Count, Damage, DataSize, Gzipped Data...
	err := writeBytes(b, int16(1), int16(2), int8(100), int16(99), int16(2), int8(2), int8(3))
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadSlotSlice(r)
	Expect(t, err, ToBeNil)
	slots, ok := v.([]Slot)
	Expect(t, ok, ToBeTrue)
	Expect(t, slots, ToBeLengthOf, 1)
	slot := slots[0]
	Expect(t, slot.ID, ToBe, int16(2))
	Expect(t, slot.Count, ToBe, int8(100))
	Expect(t, slot.Damage, ToBe, int16(99))
	Expect(t, slot.GzippedNBT, ToBeLengthOf, 2)
	Expect(t, b.Len(), ToBe, 0)
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
	Expect(t, slot.ID, ToBe, int16(2))
	Expect(t, slot.Count, ToBe, int8(100))
	Expect(t, slot.Damage, ToBe, int16(99))
	Expect(t, slot.GzippedNBT, ToBeLengthOf, 2)
	Expect(t, b.Len(), ToBe, 0)
}

func TestProtocolSlotReaderForEmptySlot(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, int16(-1))
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadSlot(r)
	Expect(t, err, ToBeNil)
	slot, ok := v.(Slot)
	Expect(t, ok, ToBeTrue)
	Expect(t, slot.ID, ToBe, int16(-1))
	Expect(t, b.Len(), ToBe, 0)
}

func TestProtocolSlotReaderForEmptyGzippedNBT(t *testing.T) {
	r, b := createProtocolReader()
	err := writeBytes(b, int16(2), int8(100), int16(99), int16(-1))
	Expect(t, err, ToBeNil)

	v, err := ProtocolReadSlot(r)
	Expect(t, err, ToBeNil)
	slot, ok := v.(Slot)
	Expect(t, ok, ToBeTrue)
	Expect(t, slot.ID, ToBe, int16(2))
	Expect(t, slot.Count, ToBe, int8(100))
	Expect(t, slot.Damage, ToBe, int16(99))
	Expect(t, slot.GzippedNBT, ToBeLengthOf, 0)
	Expect(t, b.Len(), ToBe, 0)
}
