package protocol

import (
	"fmt"
	"reflect"
	"unicode/utf16"
)

type DataWriter func(r *Writer, v interface{}) error

type DataWriters map[reflect.Type]DataWriter

func (w *DataWriters) Add(t interface{}, writer DataWriter) {
	(*w)[reflect.TypeOf(t)] = writer
}

var DefaultDataWriters = make(DataWriters)

func init() {
	// since encoding/binary supports only fixed-sized data types
	// we need to add custom parsers for the given datatypes
	DefaultDataWriters.Add("", ProtocolWriteString) // strings
	DefaultDataWriters.Add(true, ProtocolWriteBool) // bool
	DefaultDataWriters.Add([]byte{}, ProtocolWriteByteSlice)
}

/////////////////////////////////////////////////////////////////

func ProtocolWriteBool(w *Writer, v interface{}) error {
	var value byte
	if v.(bool) {
		value = 1
	} else {
		value = 0
	}
	return w.WriteValue(value)
}

func ProtocolWriteString(w *Writer, v interface{}) error {
	s := v.(string)
	size := int16(len(s))
	err := w.WriteValue(size)
	if err != nil {
		return err
	}

	raw := utf16.Encode([]rune(s))
	for _, byt := range raw {
		err = w.WriteValue(byt)
		if err != nil {
			return err
		}
	}
	return nil
}

func ProtocolWriteByteSlice(w *Writer, v interface{}) error {
	bytes := v.([]byte)
	size := int16(len(bytes))
	fmt.Printf("Writing ByteSlice of size: %v\n", size)
	err := w.WriteValue(size)
	if err != nil {
		return err
	}

	for _, b := range bytes {
		err = w.WriteValue(b)
		if err != nil {
			return err
		}
	}
	return nil
}
