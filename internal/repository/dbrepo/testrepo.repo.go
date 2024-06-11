package dbrepo

import (
	"errors"
	"time"

	"github.com/lucasvictor3/bookingsbackend/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts the reservation in the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}

	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 3 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDates returns true if availability exists for roomId, and false otherwise
func (m *testDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {
	if roomId == 10 {
		return false, errors.New("error")
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, in the given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomById returns a room by id
func (m *testDBRepo) GetRoomById(roomId int) (models.Room, error) {
	var room models.Room

	if roomId == 2 {
		return room, errors.New("Room not found")
	}
	return room, nil
}
