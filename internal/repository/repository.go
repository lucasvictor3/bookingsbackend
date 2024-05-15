package repository

import "github.com/lucasvictor3/bookingsbackend/internal/models"

// methods of the object DatabaseRepo
type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error
}
