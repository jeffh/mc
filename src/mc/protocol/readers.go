package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type NewPacketStructer interface {
	NewPacketStruct(typ PacketType) (interface{}, error)
}

////////////////////////////////////////////////

type Reader struct {
	stream  io.Reader
	readers DataReaders
	mapper  NewPacketStructer
}

func NewReader(stream io.Reader, r DataReaders, m NewPacketStructer) *Reader {
	if r == nil {
		r = DefaultDataReaders
	}
	if m == nil {
		panic(fmt.Errorf("I got a null NewPacketStructer: %#v", m))
	}
	return &Reader{
		stream:  stream,
		readers: r,
		mapper:  m,
	}
}

// Accepts a function that takes the current io.Reader and returns a
// new io.Reader for the protocol reader to use.
//
// This is generally used to promote a plain-text connection into
// an encrypted one.
func (r *Reader) UpgradeReader(f ReaderFactory) {
	old := r.stream
	r.stream = f(r.stream)
	fmt.Printf("Upgrading reader: %#v -> %#v\n", old, r.stream)
}

// Reads a given fixed-size go type and reads it in BigEndian form straight
// from the stream.
//
// The value given should be a pointer
func (r *Reader) ReadValue(v interface{}) error {
	err := binary.Read(r.stream, binary.BigEndian, v)
	// for debugging
	value := reflect.ValueOf(v).Elem()
	fmt.Printf("ReadValue: 0x%x\n", value.Interface())
	// - end
	if err != nil {
		fmt.Printf("Error when reading: %s\n", err)
	}
	return err
}

// Performs a dispatched read. Uses its reader table to invoke
// a custom function that knows how to handle the given struct.
// ReadDispatch invokes ReadStruct if given a struct not known.
// Otherwise, it reverts to ReadValue.
//
// The value given should be a pointer that is writable
func (r *Reader) ReadDispatch(value interface{}) error {
	v := reflect.ValueOf(value)
	derefV := v.Elem()
	reader, ok := r.readers[derefV.Type()]
	if ok {
		val, err := reader(r)
		if err != nil {
			return err
		}
		v.Elem().Set(reflect.ValueOf(val))
		return err
	}

	if derefV.Kind() == reflect.Struct {
		return r.ReadStruct(value)
	}

	return r.ReadValue(value)
}

// Reads a data type into a slice. The type of the slice can
// be any value that is supported by ReadDispatch.
//
// The value given should be a pointer to a slice with a given
// size to read.
//
// eg - make([]byte, 5) will read 5 bytes
func (r *Reader) ReadSlice(s interface{}) error {
	value := reflect.ValueOf(s)
	derefValue := reflect.ValueOf(value.Elem().Interface())
	size := derefValue.Len()
	for i := 0; i < size; i++ {
		typ := derefValue.Index(i).Type()
		val := reflect.New(typ)
		err := r.ReadDispatch(val.Interface())
		if err != nil {
			return err
		}
		derefValue.Index(i).Set(val.Elem())
	}
	return nil
}

// Reads data into a struct. The fields a parsed in order of the
// struct's defined fields. All fields are read using ReadDispatch.
//
// The value provided should be a pointer to a struct.
func (r *Reader) ReadStruct(v interface{}) (err error) {
	value := reflect.ValueOf(v).Elem()
	if value.Kind() != reflect.Struct {
		panic(fmt.Errorf("Expected pointer to a struct, got: %#v", v))
	}
	size := value.NumField()
	for i := 0; i < size; i++ {
		field := value.Field(i)
		val := reflect.New(field.Type())
		err = r.ReadDispatch(val.Interface())
		if err != nil {
			return
		}
		field.Set(val.Elem())
	}
	return
}

// Reads an entire minecraft packet. Dispatches based on the message
// type and returns the appropriate struct with all its fields populated.
//
// Returns an error when parsing has failed.
func (r *Reader) ReadPacket() (interface{}, error) {
	var pt PacketType
	err := r.ReadValue(&pt)
	if err != nil {
		return nil, err
	}
	value, err := r.mapper.NewPacketStruct(pt)
	if err != nil {
		return nil, err
	}
	err = r.ReadStruct(value)
	fmt.Printf("-> 0x%x >> %#v\n", pt, value)
	return value, err
}
