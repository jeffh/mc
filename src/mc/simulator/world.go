package simulator

import (
	"mc/protocol"
)

type Vector3D struct {
	X, Y, Z int32
}

type Rotation struct {
	Yaw, Pitch uint8 // fraction of 360
}

type Entity struct {
	ID        int32
	Type      byte
	Position  Vector3D
	Velocity  Vector3D
	Facing    Rotation
	HeadPitch uint8
}

type Player struct {
	EntityID  int32
	Stance    float64
	Inventory []protocol.Slot
}

type Block struct {
	TYpe     byte
	Metadata []byte // needs type
	Light    []byte // needs type
	Skylight []byte // needs type
	Add      []byte // needs type
	Biome    []byte // needs type
}

type World struct {
	Player     Player // information about the user-controlled player
	Entities   map[int32]Entity
	AgeOfWorld int64
	TimeOfDay  int64

	LevelType      string // default, flat, or largeBiomes
	GameMode       protocol.GameMode
	GameState      protocol.GameState
	GameDifficulty protocol.GameDifficulty
	MaxPlayers     byte
}
