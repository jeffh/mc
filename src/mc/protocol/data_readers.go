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

var DefaultDataReaders = make(DataReaders)

func init() {
	// since encoding/binary supports only fixed-sized data types
	// we need to add custom parsers for the given datatypes
	DefaultDataReaders.Add("", ProtocolReadString) // strings
	DefaultDataReaders.Add(true, ProtocolReadBool) // booleans

	// there are more packets that use (len int16, []byte), so this is default
	// method of parsing unless custom parsers are available for each
	DefaultDataReaders.Add([]byte{}, ProtocolReadByteSlice)
	DefaultDataReaders.Add([]Slot{}, ProtocolReadSlotSlice)
	DefaultDataReaders.Add(Slot{}, ProtocolReadSlot)
	DefaultDataReaders.Add(Int32PrefixedBytes{}, ProtocolReadInt32PrefixedBytes)
}

//////////////////////////////////////////////////////////

func ProtocolReadInt32PrefixedBytes(r *Reader) (v interface{}, err error) {

	var size int32
	err = r.ReadDispatch(&size)
	if err != nil {
		return
	}

	bytes := make(Int32PrefixedBytes, size)
	for i := int32(0); i < size; i++ {
		var byt byte
		err = r.ReadDispatch(&byt)
		if err != nil {
			return
		}
		bytes[i] = byt
	}

	v = bytes
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
	var size uint16
	err = r.ReadValue(&size)
	if err != nil {
		return
	}
	v = make([]Slot, size)
	err = r.ReadSlice(&v)
	return
}

func ProtocolReadByteSlice(r *Reader) (v interface{}, err error) {
	var size uint16
	err = r.ReadValue(&size)
	if err != nil {
		return
	}
	v = make([]byte, size)
	err = r.ReadSlice(&v)
	return
}

func ProtocolReadBool(r *Reader) (v interface{}, err error) {
	var value int8
	err = r.ReadValue(&value)
	v = (value > int8(0))
	return
}

func ProtocolReadString(r *Reader) (v interface{}, err error) {
	var size, ch uint16
	raw := make([]uint16, 0)
	err = r.ReadValue(&size)
	if err != nil {
		return
	}

	fmt.Printf("String Size: %v\n", size)

	for j := uint16(0); j < size; j++ {
		err = r.ReadValue(&ch)
		if err != nil {
			return
		}

		raw = append(raw, ch)
	}

	v = string(utf16.Decode(raw))
	fmt.Printf("String Decoded: %#v\n", v)
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

	s.GzippedNBT = make([]byte, 0)
	for i := int16(0); i < size; i++ {
		var value byte
		err = r.ReadValue(&value)
		if err != nil {
			return
		}
		s.GzippedNBT = append(s.GzippedNBT, value)
	}
	return
}
