package models

import "time"

// // Reservation holds reservation data
// type Reservation struct {
// 	FirstName string
// 	LastName  string
// 	Email     string
// 	Phone     string
// }

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Phone       string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdateAt    time.Time
}

// Room is the rooms model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdateAt  time.Time
}

// Restriction is the restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdateAt        time.Time
}

// Reservation is the reservation model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Phone     string
	Email     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdateAt  time.Time
	Room      Room
}

// RoomRestriction is the reservation model
type RoomRestriction struct {
	ID            int
	RoomID        int
	ReservationID int
	RestrictionID int
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdateAt      time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}
