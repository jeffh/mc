package protocol

import (
	"encoding/binary"
	"io"
	"reflect"
	"unicode/utf16"
)

func readBytes(b io.Reader, values ...interface{}) error {
	for _, v := range values {
		_, ok := v.(*string)
		if ok {
			var size int16
			err := binary.Read(b, binary.BigEndian, &size)
			if err != nil {
				return err
			}
			raw := make([]uint16, int(size))
			for i := int16(0); i < size; i++ {
				var ch int16
				err = binary.Read(b, binary.BigEndian, &ch)
				if err != nil {
					return err
				}
				raw[int(i)] = uint16(ch)
			}

			s := string(utf16.Decode(raw))
			reflect.ValueOf(v).Elem().Set(reflect.ValueOf(s))
		} else {
			err := binary.Read(b, binary.BigEndian, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writeBytes(b io.Writer, values ...interface{}) error {
	for _, v := range values {
		var err error
		s, ok := v.(string)
		if ok {
			err = binary.Write(b, binary.BigEndian, int16(len(s)))
			if err != nil {
				return err
			}
			raw := utf16.Encode([]rune(s))
			for _, ch := range raw {
				err = binary.Write(b, binary.BigEndian, ch)
				if err != nil {
					return err
				}
			}
		} else {
			err = binary.Write(b, binary.BigEndian, v)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
