package repository

import (
	"time"

	"github.com/lucasvictor3/bookingsbackend/internal/models"
)

// methods of the object DatabaseRepo
type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDates(start, end time.Time, roomId int) (bool, error)
}
