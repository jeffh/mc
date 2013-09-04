package simulator

import (
	"mc/protocol"
)

type Vector3D struct {
	X, Y, Z int32
}

func (v *Vector3D) Set(x, y, z int32) {
	v.X = x
	v.Y = y
	v.Z = z
}

type Rotation struct {
	Yaw, Pitch uint8 // fraction of 360
}

func (r *Rotation) Set(yaw, pitch uint8) {
	r.Yaw = yaw
	r.Pitch = pitch
}

type Entity struct {
	ID       int32
	OwnerID  int32
	Type     byte
	Position Vector3D
	Velocity Vector3D
	Facing   Rotation
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

	LevelType        string // default, flat, or largeBiomes
	GameMode         protocol.GameMode
	GameState        protocol.GameState
	GameDifficulty   protocol.GameDifficulty
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
