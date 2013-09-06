package simulator

import (
	"mc/protocol"
)

type Vector3Int struct {
	X, Y, Z int32
}

func (v *Vector3Int) Set(x, y, z int32) {
	v.X = x
	v.Y = y
	v.Z = z
}

type Vector3Float struct {
	X, Y, Z float64
}

func (v *Vector3Float) Set(x, y, z float64) {
	v.X = x
	v.Y = y
	v.Z = z
}

type RotationInt struct {
	Yaw, Pitch uint8 // fraction of 360
}

func (r *RotationInt) Set(yaw, pitch uint8) {
	r.Yaw = yaw
	r.Pitch = pitch
}

type RotationFloat struct {
	Yaw, Pitch float32
}

func (r *RotationFloat) Set(yaw, pitch float32) {
	r.Yaw = yaw
	r.Pitch = pitch
}

type Entity struct {
	ID       int32
	OwnerID  int32
	Type     protocol.EntityType
	Position Vector3Float
	Velocity Vector3Float
	Facing   RotationFloat
}

type Player struct {
	Name   string
	Online bool
	Ping   int16
}

type CurrentPlayer struct {
	Player
	Entity                    *Entity
	Stance                    float64
	HeldItemSlot              int16
	GameDifficulty            protocol.GameDifficulty
	Inventory                 []protocol.Slot
	FlyingSpeed, WalkingSpeed float32
	IsGhost                   bool // fly mode
	IsGod                     bool // god mode
	IsFlying                  bool // is in midair
}

type Block struct {
	Type     byte
	Metadata []byte // needs type
	Light    []byte // needs type
	Skylight []byte // needs type
	Add      []byte // needs type
	Biome    []byte // needs type
}

type World struct {
	CurrentPlayer CurrentPlayer // information about the user-controlled player
	Players       map[string]Player
	Entities      map[int32]*Entity
	AgeOfWorld    int64
	TimeOfDay     int64

	LevelType        protocol.LevelType // default, flat, or largeBiomes
	GameMode         protocol.GameMode
	GameState        protocol.GameState
	GameDifficulty   protocol.GameDifficulty
	GameDimension    protocol.GameDimension
	MaxPlayers       byte
	IsRaining        bool
	IsShowingCredits bool
}

func NewWorld() *World {
	return &World{
		Players:   make(map[string]Player, 0),
		Entities:  make(map[int32]*Entity, 0),
		LevelType: protocol.DefaultLevelType,
	}
}

func (w *World) NewEntityWithID(id int32) *Entity {
	e := &Entity{ID: id}
	w.Entities[e.ID] = e
	return e
}

func (w *World) EntityByID(id int32) *Entity {
	return w.Entities[id]
}
