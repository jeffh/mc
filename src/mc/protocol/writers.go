package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type GetPacketTyper interface {
	GetPacketType(v interface{}) (PacketType, error)
}

////////////////////////////////////////////////

// The Writer is the core type to serialize packets into a io.Writer.
// It is pluggable to support serializing arbitrary types.
type Writer struct {
	stream  io.Writer
	writers DataWriters
	mapper  GetPacketTyper
	Logger  Logger
}

// Creates a new writer that can write packets into the given io.Writer.
//
// Writers are the core for parsing bytes in the minecraft protocol.
// It uses a GetPacketTyper to determine the opcode (from the Packet struct),
// then utilizes the DataWriters to parse the appropriate packet.
//
// The Reader can accept a logger to use debugging internals.
//
// The last two arguments are optional, passing nil will use their default values:
// DefaultDataWriters and NullLogger.
func NewWriter(stream io.Writer, m GetPacketTyper, w DataWriters, l Logger) *Writer {
	if w == nil {
		w = DefaultDataWriters
	}
	if m == nil {
		panic(fmt.Errorf("I got a null NewPacketStructer: %#v", m))
	}
	if l == nil {
		l = &NullLogger{}
	}
	return &Writer{
		stream:  stream,
		writers: w,
		mapper:  m,
		Logger:  l,
	}
}

// Accepts a function that takes the current io.Writer and returns a
// new io.Writer for the protocol reader to use.
//
// This is generally used to promote a plain-text connection into
// an encrypted one.
func (w *Writer) UpgradeWriter(f WriterFactory) {
	w.stream = f(w.stream)
}

// Writes a fixed-size go type to the stream.
func (w *Writer) WriteValue(v interface{}) error {
	err := binary.Write(w.stream, binary.BigEndian, v)
	if err != nil {
		w.Logger.Printf("Error when writing %#v: %s\n", v, err)
	}
	return err
}

// Performs a dispatched write. Uses its writer table to invoke
// a custom function that knows how to handle the given struct.
// WriteDispatch invokes WriteStruct if given a struct not known.
// Otherwise, it reverts to WriteValue.
//
// The value given can be anything that is supported by writers
// table, a struct, or fixed-size go type.
func (w *Writer) WriteDispatch(value interface{}) error {
	v := reflect.ValueOf(value)
	writer, ok := w.writers[v.Type()]
	if ok {
		err := writer(w, v.Interface())
		return err
	}

	if v.Kind() == reflect.Struct {
		return w.WriteStruct(value)
	}

	return w.WriteValue(v.Interface())
}

// Writes a slice's data into the stream. The type of the slice can
// be any value that is supported by ReadDispatch.
//
// The value given should be a slice to write its data
func (w *Writer) WriteStruct(v interface{}) (err error) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	size := value.NumField()
	for i := 0; i < size; i++ {
		field := value.Field(i)
		err = w.WriteDispatch(field.Interface())
		if err != nil {
			return err
		}
	}
	return
}

// Writes a minecraft packet into the data stream. It also handles
// writing the proper packet type prefix before writing the struct
// provided.
func (w *Writer) WritePacket(v interface{}) error {
	w.Logger.Printf("C->S %#v\n", v)
	pt, err := w.mapper.GetPacketType(v)
	if err != nil {
		return err
	}
	w.WriteValue(pt)
	return w.WriteStruct(v)
}
