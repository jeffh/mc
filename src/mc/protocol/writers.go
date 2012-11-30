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

type Writer struct {
	stream  io.Writer
	writers DataWriters
	mapper  GetPacketTyper
	Logger  Logger
}

func NewWriter(stream io.Writer, m GetPacketTyper, w DataWriters, l Logger) *Writer {
	if w == nil {
		w = DefaultDataWriters
	}
	if m == nil {
		panic(fmt.Errorf("I got a null NewPacketStructer: %#v", m))
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
	fmt.Printf("WriteValue: 0x%x\n", v)
	if err != nil {
		fmt.Printf("Error when writing %#v: %s\n", v, err)
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
		w.WriteStruct(value)
	}

	return w.WriteValue(v.Interface())
}

// Writes a slice's data into the stream. The type of the slice can
// be any value that is supported by ReadDispatch.
//
// The value given should be a slice to write its data
func (w *Writer) WriteStruct(v interface{}) (err error) {
	value := reflect.ValueOf(v).Elem()
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
	pt, err := w.mapper.GetPacketType(v)
	if err != nil {
		return err
	}
	w.WriteValue(pt)
	defer func() {
		fmt.Printf("<- %#v\n", v)
	}()
	return w.WriteStruct(v)
}
