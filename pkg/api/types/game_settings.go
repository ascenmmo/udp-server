package types

import (
	"github.com/google/uuid"
	"time"
)

type CreateRoomRequest struct {
	RoomTTl time.Duration
}

type GetDeletedRooms struct {
	GameID uuid.UUID
	RoomID uuid.UUID
}
