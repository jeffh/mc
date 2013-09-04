package simulator

import (
	"fmt"
	"mc/protocol"
	"nbt"
)

type Simulator struct {
	World *World
}

func NewSimulator() *Simulator {
	return &Simulator{
		World: NewWorld(),
	}
}

func (s *Simulator) ProcessMessage(v interface{}) {
	switch t := v.(type) {
	case protocol.LoginRequest:
		s.World.CurrentPlayer.Entity = s.World.NewEntityWithID(t.EntityID)
		s.World.LevelType = t.LevelType
		s.World.GameMode = t.GameMode
		s.World.GameDimension = t.Dimension
		s.World.GameDifficulty = t.Difficulty
		s.World.CurrentPlayer.Difficulty = t.Difficulty
	case protocol.SpawnPosition:
		s.World.CurrentPlayer.Position.Set(t.X, t.Y, t.Z)
	case protocol.PlayerAbilities:
		s.World.CurrentPlayer.IsGod = t.IsGod()
		s.World.CurrentPlayer.IsGhost = t.IsGhost()
		s.World.CurrentPlayer.FlyingSpeed = t.FlyingSpeed
		s.World.CurrentPlayer.WalkingSpeed = t.WalkingSpeed
	case protocol.TimeUpdate:
		s.World.AgeOfWorld = t.WorldAge
		s.World.TimeOfDay = t.TimeOfDay
	case protocol.ChangeGameState:
		switch t.Reason {
		case protocol.GameStateBeginRain:
			s.World.IsRaining = true
		case protocol.GameStateEndRain:
			s.World.IsRaining = false
		case protocol.GameStateChangeGameMode:
			s.World.GameMode = t.GameMode
		case protocol.GameStateEnterCredits:
			s.World.IsShowingCredits = true
		}
	case protocol.HeldItemChange:
		s.World.CurrentPlayer.HeldItemSlot = t.SlotID
	case protocol.PlayerListItem:
		s.World.Players[t.Name] = Player{
			Name:   t.Name,
			Online: t.Online,
			ping:   t.Ping,
		}
	case protocol.PlayerPositionLookForClient:
		s.World.CurrentPlayer.Position.Set(t.X, t.Y, t.Z)
		s.World.CurrentPlayer.Rotation.Set(t.Yaw, t.Pitch)
		s.World.CurrentPlayer.Stance = t.Stance
		s.World.CurrentPlayer.IsFlying = !t.IsOnGround
	case protocol.SetWindowItems:
		if t.WindowID == protocol.WindowTypeInventory {
			s.World.CurrentPlayer.Inventory = t.Slots
		} else {
			panic(fmt.Errorf("Unknown window id: %d", t.WindowID))
		}
	case protocol.SetSlot:
		if t.IsHeld() {
			// not sure what to do about this...
		}
	case protocol.MapChunkBulk:
	case protocol.SpawnObject:
		entity := s.World.NewEntityWithID(t.EntityID)
		entity.Position.Set(t.X, t.Y, t.Z)
		if t.HasVelocity() {
			entity.Velocity.Set(t.XVelocity, t.YVelocity, t.ZVelocity)
		}
		entity.Facing.Set(t.Yaw, t.Pitch)
		entity.OwnerID = t.OwnerID
		entity.Type = t.Type
	}
}
