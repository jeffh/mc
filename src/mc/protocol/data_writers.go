package protocol

import (
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

	DefaultDataWriters.Add([]string{}, ProtocolWriteStringSlice)
	DefaultDataWriters.Add([]Slot{}, ProtocolWriteSlotSlice)
	DefaultDataWriters.Add(Slot{}, ProtocolWriteSlot)

	DefaultDataWriters.Add([]EntityMetadata{}, ProtocolWriteEntityMetadataSlice)
}

/////////////////////////////////////////////////////////////////

func entityKey(id EntityMetadataIndex, typ EntityMetadataType) byte {
	return byte(id) | (byte(typ) << 5)
}

/////////////////////////////////////////////////////////////////

func ProtocolWriteEntityMetadataSlice(w *Writer, v interface{}) error {
	metadatas := v.([]EntityMetadata)

	for _, md := range metadatas {
		err := w.WriteValue(entityKey(md.ID, md.Type))
		if err != nil {
			return err
		}
		err = w.WriteDispatch(md.Value)
		if err != nil {
			return err
		}
	}

	return w.WriteValue(byte(127))
}

func ProtocolWriteSlot(w *Writer, v interface{}) error {
	slot := v.(Slot)

	err := w.WriteValue(slot.ID)
	if err != nil {
		return err
	}

	// empty slots don't write anything else
	if slot.IsEmpty() {
		return nil
	}

	err = w.WriteValue(slot.Count)
	if err != nil {
		return err
	}

	err = w.WriteValue(slot.Damage)
	if err != nil {
		return err
	}

	if len(slot.GzippedNBT) == 0 {
		return w.WriteValue(int16(-1))
	}

	err = w.WriteValue(int16(len(slot.GzippedNBT)))
	if err != nil {
		return err
	}

	return nil
}

func ProtocolWriteSlotSlice(w *Writer, v interface{}) error {
	slots := v.([]Slot)
	size := int16(len(slots))
	err := w.WriteValue(size)
	if err != nil {
		return err
	}

	for _, s := range slots {
		err = w.WriteDispatch(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func ProtocolWriteStringSlice(w *Writer, v interface{}) error {
	strings := v.([]string)
	size := int16(len(strings))

	err := w.WriteValue(size)
	if err != nil {
		return err
	}

	for _, s := range strings {
		err = w.WriteDispatch(s)
		if err != nil {
			return err
		}
	}
	return nil
}

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
