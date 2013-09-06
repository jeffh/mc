package simulator

import (
	"ax"
	"bytes"
	"compress/zlib"
	"fmt"
	"mc/protocol"
	"smpm"
)

type MapChunkBulkWrapper struct {
	MapChunkBulk *protocol.MapChunkBulk
	index        int
}

func (w *MapChunkBulkWrapper) HasSkylightData() bool {
	return w.MapChunkBulk.SkylightSent
}

func (w *MapChunkBulkWrapper) ChunkColumnCount() int16 {
	return int16(len(w.MapChunkBulk.Metadatas))
}

func (w *MapChunkBulkWrapper) NextMetadata() smpm.ChunkColumnMetadata {
	metadata := w.MapChunkBulk.Metadatas[w.index]
	w.index += 1
	return smpm.ChunkColumnMetadata{
		X:             metadata.ChunkX,
		Z:             metadata.ChunkY,
		PrimaryBitmap: metadata.PrimaryBitmap,
		AddBitmap:     metadata.AddBitmap,
	}
}

func (w *MapChunkBulkWrapper) IsGroundUpContinuous() bool {
	return true
}

type Simulator struct {
	World  *World
	Logger ax.WrapLogger
}

func NewSimulator(logger ax.Logger) *Simulator {
	return &Simulator{
		World:  NewWorld(),
		Logger: ax.Wrap(ax.Use(logger), ax.NewPrefixLogger("[simulator] ")),
	}
}

func (s *Simulator) ProcessMessage(v interface{}) {
	switch t := v.(type) {
	case *protocol.LoginRequest:
		s.World.CurrentPlayer.Entity = s.World.NewEntityWithID(t.EntityID)
		s.World.LevelType = t.LevelType
		s.World.GameMode = t.GameMode
		s.World.GameDimension = t.Dimension
		s.World.GameDifficulty = t.Difficulty
		s.World.CurrentPlayer.GameDifficulty = t.Difficulty
	case *protocol.SpawnPosition:
		s.World.CurrentPlayer.Entity.Position.Set(float64(t.X), float64(t.Y), float64(t.Z))
	case *protocol.PlayerAbilities:
		s.World.CurrentPlayer.IsGod = t.IsGod()
		s.World.CurrentPlayer.IsGhost = t.IsGhost()
		s.World.CurrentPlayer.FlyingSpeed = t.FlyingSpeed
		s.World.CurrentPlayer.WalkingSpeed = t.WalkingSpeed
	case *protocol.TimeUpdate:
		s.World.AgeOfWorld = t.WorldAge
		s.World.TimeOfDay = t.TimeOfDay
	case *protocol.ChangeGameState:
		switch t.State {
		case protocol.GameStateBeginRain:
			s.World.IsRaining = true
		case protocol.GameStateEndRain:
			s.World.IsRaining = false
		case protocol.GameStateChangeGameMode:
			s.World.GameMode = t.GameMode
		case protocol.GameStateEnterCredits:
			s.World.IsShowingCredits = true
		}
	case *protocol.HeldItemChange:
		s.World.CurrentPlayer.HeldItemSlot = t.SlotID
	case *protocol.PlayerListItem:
		s.World.Players[t.Name] = Player{
			Name:   t.Name,
			Online: t.Online,
			Ping:   t.Ping,
		}
	case *protocol.PlayerPositionLookForClient:
		s.World.CurrentPlayer.Entity.Position.Set(t.X, t.Y, t.Z)
		s.World.CurrentPlayer.Entity.Facing.Set(t.Yaw, t.Pitch)
		s.World.CurrentPlayer.Stance = t.Stance
		s.World.CurrentPlayer.IsFlying = !t.IsOnGround
	case *protocol.SetWindowItems:
		if t.WindowID == protocol.WindowTypeInventory {
			s.World.CurrentPlayer.Inventory = t.Slots
		} else {
			panic(fmt.Errorf("Unknown window id: %d", t.WindowID))
		}
	case *protocol.SetSlot:
		if t.IsHeld() {
			// not sure what to do about this...
		}
	case *protocol.MapChunkBulk:
		buffer := bytes.NewBuffer(t.CompressedData)
		reader, err := zlib.NewReader(buffer)
		if err != nil {
			panic(err)
		}
		file := smpm.NewFile(reader, &MapChunkBulkWrapper{MapChunkBulk: t}, s.Logger.WrappedLogger())
		columns, err := file.Parse()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Columns: %d\n", len(columns))
	case *protocol.SpawnObject:
		entity := s.World.NewEntityWithID(t.EntityID)
		entity.Position.Set(float64(t.X), float64(t.Y), float64(t.Z))
		if t.HasVelocity() {
			entity.Velocity.Set(float64(t.XVelocity), float64(t.YVelocity), float64(t.ZVelocity))
		}
		entity.Facing.Set(float32(t.Yaw), float32(t.Pitch))
		entity.OwnerID = t.OwnerEntityID
		entity.Type = t.Type
	}
}
