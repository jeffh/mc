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

// Reader is the core type to parsing bytes from an io.Reader.
// It is pluggable to support parsing arbitrary types.
type Reader struct {
	stream  io.Reader
	readers DataReaders
	mapper  NewPacketStructer
	Logger  Logger
}

// Creates a new Reader for the given io.Reader.
//
// Readers are the core for parsing bytes in the minecraft protocol.
// It uses a NewPackerStructer to determine the packet type (from the opcode),
// then utilizes the DataReaders to parse the appropriate packet.
//
// The Reader can accept a logger to use debugging internals.
//
// The last two arguments are optional, passing nil will use their default values:
// DefaultDataReaders and NullLogger.
func NewReader(stream io.Reader, m NewPacketStructer, r DataReaders, l Logger) *Reader {
	if r == nil {
		r = DefaultDataReaders
	}
	if m == nil {
		panic(fmt.Errorf("I got a null NewPacketStructer: %#v", m))
	}
	if l == nil {
		l = &NullLogger{}
	}
	return &Reader{
		stream:  stream,
		readers: r,
		mapper:  m,
		Logger:  l,
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
	r.Logger.Printf("Upgrading reader: %#v -> %#v\n", old, r.stream)
}

// Reads a given fixed-size go type and reads it in BigEndian form straight
// from the stream.
//
// The value given should be a pointer
func (r *Reader) ReadValue(v interface{}) error {
	err := binary.Read(r.stream, binary.BigEndian, v)
	// for debugging
	//fmt.Printf("ReadValue: 0x%x\n", reflect.ValueOf(v).Elem().Interface())
	// - end
	if err != nil {
		r.Logger.Printf("Error when reading: %s\n", err)
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
	r.Logger.Printf("S->C 0x%x\n", pt)
	if err != nil {
		return nil, err
	}
	err = r.ReadDispatch(value)
	r.Logger.Printf("      >> %#v\n", value)
	return value, err
}
