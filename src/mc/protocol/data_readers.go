package protocol

import (
	"fmt"
	"reflect"
	"unicode/utf16"
)

type DataReader func(r *Reader) (interface{}, error)

type DataReaders map[reflect.Type]DataReader

func (r *DataReaders) Add(t interface{}, reader DataReader) {
	(*r)[reflect.TypeOf(t)] = reader
}

// The default custom data readers for reading custom types
// from an io.Reader
var DefaultDataReaders = make(DataReaders)

func init() {
	// since encoding/binary supports only fixed-sized data types
	// we need to add custom parsers for the given datatypes
	DefaultDataReaders.Add("", ProtocolReadString)               // strings
	DefaultDataReaders.Add(LevelType(""), ProtocolReadLevelType) // strings
	DefaultDataReaders.Add(true, ProtocolReadBool)               // booleans

	// there are more packets that use (len int16, []byte), so this is default
	// method of parsing unless custom parsers are available for each
	DefaultDataReaders.Add([]string{}, ProtocolReadStringSlice)
	DefaultDataReaders.Add([]byte{}, ProtocolReadByteSlice)
	DefaultDataReaders.Add([]Slot{}, ProtocolReadSlotSlice)
	DefaultDataReaders.Add(Slot{}, ProtocolReadSlot)
	DefaultDataReaders.Add(Int32PrefixedBytes{}, ProtocolReadInt32PrefixedBytes)

	DefaultDataReaders.Add([]EntityMetadata{}, ProtocolReadEntityMetadataSlice)
	DefaultDataReaders.Add(DestroyEntity{}, ProtocolReadDestroyEntity)
	DefaultDataReaders.Add(MapChunkBulk{}, ProtocolReadMapChunkBulk)

	DefaultDataReaders.Add(SpawnObject{}, ProtocolReadSpawnObject)           // needs test
	DefaultDataReaders.Add(EntityProperties{}, ProtocolReadEntityProperties) // needs test
}

//////////////////////////////////////////////////////////

func ProtocolReadLevelType(r *Reader) (v interface{}, err error) {
	var value interface{}
	value, err = ProtocolReadString(r)
	str := value.(string)
	v = LevelType(str)
	return
}

func ProtocolReadDestroyEntity(r *Reader) (v interface{}, err error) {
	var destroyEntity DestroyEntity
	var size byte
	err = r.ReadValue(&size)
	if err != nil {
		return
	}

	destroyEntity.EntityIDs = make([]int32, size)
	err = r.ReadSlice(&destroyEntity.EntityIDs)
	v = destroyEntity
	return
}

// Handles the reading an array of bytes, prefixed by a signed 32-bit
// integer from a given Reader.
func ProtocolReadInt32PrefixedBytes(r *Reader) (v interface{}, err error) {
	var size int32
	err = r.ReadDispatch(&size)
	if err != nil {
		return
	}

	v = make(Int32PrefixedBytes, size)
	err = r.ReadSlice(&v)
	return
}

// Handles reading an array of strings, prefixed by a signed short
// from a given reader.
func ProtocolReadStringSlice(r *Reader) (v interface{}, err error) {
	var size int16
	err = r.ReadValue(&size)
	if err != nil {
		return
	}

	v = make([]string, size)
	err = r.ReadSlice(&v)
	return
}

func ProtocolReadMapChunkBulk(r *Reader) (v interface{}, err error) {
	var chunk MapChunkBulk
	defer func() { v = chunk }()

	var metadataSize int16
	err = r.ReadValue(&metadataSize)
	if err != nil {
		return
	}

	var dataSize int32
	err = r.ReadValue(&dataSize)
	if err != nil {
		return
	}

	err = r.ReadDispatch(&chunk.SkylightSent)
	if err != nil {
		return
	}

	chunk.CompressedData = make([]byte, dataSize)
	err = r.ReadSlice(&chunk.CompressedData)
	if err != nil {
		return
	}

	chunk.Metadatas = make([]ChunkBulkMetadata, metadataSize)
	err = r.ReadSlice(&chunk.Metadatas)
	v = chunk
	return
}

func ProtocolReadEntityMetadataSlice(r *Reader) (v interface{}, err error) {
	slice := make([]EntityMetadata, 0)

	var b byte
	for {
		err = r.ReadValue(&b)
		if err != nil {
			return
		}
		if b == byte(127) {
			break
		}
		// lower 5 bits is ID (keys)
		// upper 3 bits is type
		em := EntityMetadata{
			ID:   EntityMetadataIndex(b & 0x1F),
			Type: EntityMetadataType((b & 0xE0) >> 5),
		}
		switch em.Type {
		case EntityMetadataByte:
			byt := byte(0)
			err = r.ReadDispatch(&byt)
			em.Value = byt
		case EntityMetadataShort:
			i := int16(0)
			err = r.ReadDispatch(&i)
			em.Value = i
		case EntityMetadataInt:
			i := int32(0)
			err = r.ReadDispatch(&i)
			em.Value = i
		case EntityMetadataFloat:
			f := float32(0)
			err = r.ReadDispatch(&f)
			em.Value = f
		case EntityMetadataString:
			s := ""
			err = r.ReadDispatch(&s)
			em.Value = s
		case EntityMetadataSlot:
			s := Slot{}
			err = r.ReadDispatch(&s)
			em.Value = s
		case EntityMetadataPosition:
			p := Position{}
			err = r.ReadDispatch(&p)
			em.Value = p
		default:
			err = fmt.Errorf("Unsupported EntityType: (got 0x%x)", em.Type)
		}

		if err != nil {
			return
		}

		slice = append(slice, em)
	}

	v = slice
	return
}

func ProtocolReadSlotSlice(r *Reader) (v interface{}, err error) {
	var size int16
	err = r.ReadValue(&size)
	if err != nil {
		return
	}
	v = make([]Slot, size)
	err = r.ReadSlice(&v)
	return
}

func ProtocolReadByteSlice(r *Reader) (v interface{}, err error) {
	var size int16
	err = r.ReadValue(&size)
	if err != nil {
		return
	}
	v = make([]byte, size)
	err = r.ReadSlice(&v)
	return
}

func ProtocolReadBool(r *Reader) (v interface{}, err error) {
	var value byte
	err = r.ReadValue(&value)
	v = (value > byte(0))
	return
}

func ProtocolReadString(r *Reader) (v interface{}, err error) {
	var size int16
	var ch uint16
	raw := make([]uint16, 0)
	err = r.ReadValue(&size)
	if err != nil {
		return
	}

	for j := int16(0); j < size; j++ {
		err = r.ReadValue(&ch)
		if err != nil {
			return
		}

		raw = append(raw, ch)
	}

	v = string(utf16.Decode(raw))
	return
}

func ProtocolReadSlot(r *Reader) (v interface{}, err error) {
	var s Slot
	defer func() { v = s }()

	err = r.ReadValue(&s.ID)
	if err != nil {
		return
	}
	if s.ID == -1 {
		return
	}

	err = r.ReadValue(&s.Count)
	if err != nil {
		return
	}

	err = r.ReadValue(&s.Damage)
	if err != nil {
		return
	}

	var size int16
	err = r.ReadValue(&size)
	if err != nil {
		return
	}

	if size < 0 {
		size = 0
	}

	s.GzippedNBT = make([]byte, size)
	err = r.ReadSlice(&s.GzippedNBT)
	return
}

func ProtocolReadSpawnObject(r *Reader) (v interface{}, err error) {
	var spawnObject SpawnObject
	defer func() { v = spawnObject }()

	err = r.ReadValue(&spawnObject.EntityID)
	if err != nil {
		return
	}

	err = r.ReadValue(&spawnObject.Type)
	if err != nil {
		return
	}

	err = r.ReadValue(&spawnObject.X)
	if err != nil {
		return
	}

	err = r.ReadValue(&spawnObject.Y)
	if err != nil {
		return
	}

	err = r.ReadValue(&spawnObject.Z)
	if err != nil {
		return
	}

	err = r.ReadValue(&spawnObject.Pitch)
	if err != nil {
		return
	}

	err = r.ReadValue(&spawnObject.Yaw)
	if err != nil {
		return
	}

	err = r.ReadValue(&spawnObject.Flag)
	if err != nil {
		return
	}

	if spawnObject.Type == EntityItemFrame {
		err = r.ReadValue(&spawnObject.Orientation)
	} else if spawnObject.Type == EntityFallingObject {
		err = r.ReadValue(&spawnObject.BlockType)
	} else if IsProjectileEntity(spawnObject.Type) {
		err = r.ReadValue(&spawnObject.OwnerEntityID)
	}
	if err != nil {
		return
	}

	if spawnObject.Flag > 0 {
		err = r.ReadValue(&spawnObject.XVelocity)
		if err != nil {
			return
		}

		err = r.ReadValue(&spawnObject.YVelocity)
		if err != nil {
			return
		}

		err = r.ReadValue(&spawnObject.ZVelocity)
		if err != nil {
			return
		}
	}
	return
}

func ProtocolReadEntityProperties(r *Reader) (v interface{}, err error) {
	var e EntityProperties
	defer func() { v = e }()

	err = r.ReadValue(&e.EntityID)
	if err != nil {
		return
	}

	var count int32
	err = r.ReadValue(&count)
	if err != nil {
		return
	}

	e.Properties = make([]EntityProperty, count)

	for i := int32(0); i < count; i++ {
		var property EntityProperty
		err = r.ReadDispatch(&property.Key)
		if err != nil {
			return
		}

		err = r.ReadValue(&property.Value)
		if err != nil {
			return
		}

		var size int16
		err = r.ReadValue(&size)
		if err != nil {
			return
		}

		property.Attributes = make([]EntityAttribute, size)
		err = r.ReadSlice(&property.Attributes)
		if err != nil {
			return
		}

		e.Properties[i] = property
	}
	return
}
