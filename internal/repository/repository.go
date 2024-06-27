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
	SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomById(roomId int) (models.Room, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)

	UpdateUser(user models.User) error
	GetUserById(id int) (models.User, error)
	Authenticate(email, testPassword string) (int, string, error)
}
