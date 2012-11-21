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
}

//////////////////////////////////////////////////////////

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

	s.Data = make([]byte, 0)
	for i := int16(0); i < size; i++ {
		var value byte
		err = r.ReadValue(&value)
		if err != nil {
			return
		}
		// currently, just toss it all
		// the data is gzipped-NBT format
	}
	return
}
