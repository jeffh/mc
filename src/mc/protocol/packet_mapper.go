package protocol

import (
	"fmt"
	"reflect"
)

// PacketMapper is an interface that conforms to Minecraft message
// types (also referred to as Packets), to their appropriate structs.
//
// The mapper is bi-directional. So the mapper knows which struct
// corresponds to the given PacketType and how to create a new
// struct of the correct type for the given PacketType
type PacketMapper interface {
	NewPacketStruct(typ PacketType) (interface{}, error)
	GetPacketType(v interface{}) (PacketType, error)
}

///////////////////////////////////////////////////////
// helper functions

// Coerces reflect.Type into a unique string identifier
func typeToStr(t reflect.Type) string {
	return t.PkgPath() + "." + t.Name()
}

///////////////////////////////////////////////////////

// PacketMapper is a generic mapper that can map all of
// PacketType to specific Type and back again.
type StdPacketMapper struct {
	parent         PacketMapper
	packetToStruct map[PacketType]reflect.Type
	structToPacket map[string]PacketType
}

// Creates a standard packet mapper with an optional parent
// to delegate to.
func NewStdPacketMapper(p PacketMapper) *StdPacketMapper {
	return &StdPacketMapper{
		parent:         p,
		packetToStruct: make(map[PacketType]reflect.Type),
		structToPacket: make(map[string]PacketType),
	}
}

// Defines a two-way mapping: PacketType <=> struct
//
// Panics if the given PacketType or struct is already defined.
// Use Set(PacketType, interface{}) to avoid panics
func (m *StdPacketMapper) Define(t PacketType, v interface{}) {
	m.DefineIncoming(t, v)
	m.DefineOutgoing(t, v)
}

// Defines a two-way mapping: PacketType <=> struct
//
// Unlike Define(), Set() will silently overwrite existing definitions.
func (m *StdPacketMapper) Set(t PacketType, v interface{}) {
	m.SetIncoming(t, v)
	m.SetOutgoing(t, v)
}

// Defines a one-way mapping: PacketType => struc
//
// Panics if there is already a definition for the given PacketType
// Use SetIncoming(PacketType, interface{}) if you want to avoid panics
func (m *StdPacketMapper) DefineIncoming(t PacketType, v interface{}) {
	_, ok := m.packetToStruct[t]
	if ok {
		panic(fmt.Errorf("Incoming Packet Type already defined: 0x%x", t))
	}
	m.SetIncoming(t, v)
}

// Defines a one-way mapping: struct => PacketType
//
// Panics if there is already a definition for the given type of struct
// Use SetIncoming(PacketType, interface{}) if you want to avoid panics
func (m *StdPacketMapper) DefineOutgoing(t PacketType, v interface{}) {
	typ := reflect.TypeOf(v)
	_, ok := m.structToPacket[typeToStr(typ)]
	if ok {
		panic(fmt.Errorf("Outgoing Packet Type already defined: 0x%x", t))
	}
	m.SetOutgoing(t, v)
}

// Defines a one-way mapping: PacketType => struct
//
// Unlike DefineIncoming(), SetIncoming() will silently overwrite
// existing definitions.
func (m *StdPacketMapper) SetIncoming(t PacketType, v interface{}) {
	typ := reflect.TypeOf(v)
	m.packetToStruct[t] = typ
}

// Defines a one-way mapping: struct => PacketType
//
// Unlike DefineIncoming(), SetOutgoing() will silently overwrite
// existing definitions.
func (m *StdPacketMapper) SetOutgoing(t PacketType, v interface{}) {
	typ := reflect.TypeOf(v)
	m.structToPacket[typeToStr(typ)] = t
}

// NewPacketStruct creates a new struct that compiles to the given
// packet type.
//
// Returns an error if there is no mapping for the given PacketType
func (m *StdPacketMapper) NewPacketStruct(typ PacketType) (interface{}, error) {
	t, ok := m.packetToStruct[typ]
	if ok {
		return reflect.New(t).Interface(), nil
	}

	if m.parent != nil {
		return m.parent.NewPacketStruct(typ)
	}

	return PacketType(0), fmt.Errorf("Invalid PacketType: %#v", typ)
}

// GetPacketType returns the PacketType that corresponds to the given
// struct.
//
// Returns an error if there is no mapping
func (m *StdPacketMapper) GetPacketType(v interface{}) (PacketType, error) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	tstr := typeToStr(value.Type())
	typ, ok := m.structToPacket[tstr]
	if ok {
		return typ, nil
	}

	if m.parent != nil {
		return m.parent.GetPacketType(v)
	}

	return PacketType(0), fmt.Errorf("Unexpected Struct Type: %#v", v)
}
