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
}

func NewWriter(stream io.Writer, w DataWriters, m GetPacketTyper) *Writer {
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
	}
}

func (w *Writer) UpgradeWriter(f WriterFactory) {
	w.stream = f(w.stream)
}

func (w *Writer) WriteValue(v interface{}) error {
	err := binary.Write(w.stream, binary.BigEndian, v)
	fmt.Printf("WriteValue: 0x%x\n", v)
	if err != nil {
		fmt.Printf("Error when writing %#v: %s\n", v, err)
	}
	return err
}

func (w *Writer) Writer() io.Writer {
	return w.stream
}

func (w *Writer) WriteDispatch(v reflect.Value) error {
	writer, ok := w.writers[v.Type()]
	if ok {
		err := writer(w, v.Interface())
		return err
	}
	return w.WriteValue(v.Interface())
}

func (w *Writer) WriteStruct(v interface{}) (err error) {
	value := reflect.ValueOf(v).Elem()
	size := value.NumField()
	for i := 0; i < size; i++ {
		field := value.Field(i)
		err = w.WriteDispatch(field)
		if err != nil {
			return err
		}
	}
	return
}

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
