package models

import "time"

// User, the user model from schema
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room, the rooms model from schema
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction, the restrictions model from schema
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation, the reservation model from schema
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room // we dont always need to have a perfect model mapping from here to the schema. In this case we add another property that will add all the room information in the model
}

// RoomRestriction, the RoomRestrictions model from schema
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	// we may not use all the below, but it is useful to have them here just in case, thus it is easier to pull info from the model
	Room        Room
	Reservation Reservation
	Restriction Restriction
}
