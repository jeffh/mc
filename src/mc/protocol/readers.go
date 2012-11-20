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

func (r *Reader) UpgradeReader(f ReaderFactory) {
	old := r.stream
	r.stream = f(r.stream)
	fmt.Printf("Upgrading reader: %#v -> %#v\n", old, r.stream)
}

func (r *Reader) ReadValue(v interface{}) error {
	err := binary.Read(r.stream, binary.BigEndian, v)
	value := reflect.ValueOf(v)
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	fmt.Printf("ReadValue: 0x%x\n", value.Interface())
	if err != nil {
		fmt.Printf("Error when reading: %s\n", err)
	}
	return err
}

func (r *Reader) Reader() io.Reader {
	return r.stream
}

func (r *Reader) ReadDispatch(v reflect.Value) error {
	reader, ok := r.readers[v.Type()]
	if ok {
		val, err := reader(r)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(val))
		return err
	}
	return r.ReadValue(v.Addr().Interface())
}

func (r *Reader) ReadSlice(s interface{}) error {
	value := reflect.ValueOf(s)
	size := value.Elem().Len()
	for i := 0; i < size; i++ {
		typ := value.Elem().Index(i).Type()
		val := reflect.New(typ)
		err := r.ReadValue(val.Interface())
		if err != nil {
			return err
		}
		value.Elem().Index(i).Set(val.Elem())
	}
	return nil
}

func (r *Reader) ReadStruct(v interface{}) (err error) {
	value := reflect.ValueOf(v).Elem()
	size := value.NumField()
	for i := 0; i < size; i++ {
		field := value.Field(i)
		err = r.ReadDispatch(field)
		if err != nil {
			return
		}
	}
	return
}

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
